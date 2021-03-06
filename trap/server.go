/*
 * Trap
 * An anti-pryer server for better privacy
 *
 * This file is a part of Trap project
 *
 * Copyright 2016 Rain Lee <raincious@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package trap

import (
	"github.com/raincious/trap/trap/core/client"
	"github.com/raincious/trap/trap/core/event"
	"github.com/raincious/trap/trap/core/listen"
	"github.com/raincious/trap/trap/core/logger"
	"github.com/raincious/trap/trap/core/server"
	"github.com/raincious/trap/trap/core/types"

	"sync"
	"time"
)

type Server struct {
	logger                  *logger.Logger
	listen                  *listen.Listen
	event                   *event.Event
	clientMaps              *client.Clients
	clientRWLock            types.Mutex
	clientCronExitCh        chan bool
	clientMaxRecords        types.UInt16
	clientMaxRecordMaxBytes types.UInt32
	serverUpped             bool
	serverUpping            bool
	serverLock              types.Mutex
	serverDownWait          sync.WaitGroup
	timeout                 time.Duration
	tolerate                types.UInt32
	tolerateExpire          time.Duration
	tolerateRestrict        time.Duration
	concurrentLimit         types.UInt16
	onUpCommands            types.Callbacks
	onDownCommands          types.Callbacks
	onUpDownCommands        []types.CallbackPair
	onMarkCommands          []func(server.ClientInfo)
	onUnmarkCommands        []func(types.IP)
	bootTime                time.Time
	totalInbound            types.UInt64
	totalMarked             types.UInt64
	totalHit                types.UInt64
	history                 server.Histories
	distribution            server.Distributions
}

func NewServer() *Server {
	return &Server{
		clientRWLock:            types.Mutex{},
		serverLock:              types.Mutex{},
		serverDownWait:          sync.WaitGroup{},
		timeout:                 1 * time.Second,
		tolerate:                1,
		tolerateExpire:          3600 * time.Second,
		tolerateRestrict:        3600 * time.Second,
		concurrentLimit:         10,
		clientMaxRecords:        16,
		clientMaxRecordMaxBytes: 512,
		history:                 server.Histories{},
		distribution:            server.Distributions{},
	}
}

func (this *Server) SetLogger(l *logger.Logger) {
	this.logger = l.NewContext("Server")

	this.logger.Debugf("`Logger` is set")
}

func (this *Server) SetTolerate(limit types.UInt32, expire time.Duration,
	restrict time.Duration) {
	this.tolerate = limit
	this.tolerateExpire = expire
	this.tolerateRestrict = restrict

	this.logger.Debugf("Tolerate has been set to '%d' attempts within '%s'"+
		", restrict period is '%s'",
		limit, expire, restrict)
}

func (this *Server) SetClientRecordLimit(l types.UInt16) {
	this.clientMaxRecords = l

	this.logger.Debugf("Client Data Record Limit has been set to '%d' items", l)
}

func (this *Server) SetClientRecordDataLimit(limit types.UInt32) {
	this.clientMaxRecordMaxBytes = limit

	this.logger.Debugf("Client retrieve limit now been set to maximum '%d'"+
		" bytes",
		this.clientMaxRecordMaxBytes)
}

func (this *Server) SetTimeout(t time.Duration) {
	this.timeout = t

	this.logger.Debugf("Timeout has been set to '%s'", t)
}

func (this *Server) SetConcurrentLimit(c types.UInt16) {
	this.concurrentLimit = c

	this.logger.Debugf("Concurrent Limit has been set to '%d'", c)
}

func (this *Server) OnUp(f types.Callback) {
	this.onUpCommands = append(this.onUpCommands, f)
}

func (this *Server) OnDown(f types.Callback) {
	this.onDownCommands = append(this.onDownCommands, f)
}

func (this *Server) OnUpDown(up types.Callback, down types.Callback) {
	this.onUpDownCommands = append(this.onUpDownCommands, types.CallbackPair{
		Alpha: up,
		Beta:  down,
	})
}

func (this *Server) OnMark(f func(server.ClientInfo)) {
	this.onMarkCommands = append(this.onMarkCommands, f)
}

func (this *Server) OnUnmark(f func(types.IP)) {
	this.onUnmarkCommands = append(this.onUnmarkCommands, f)
}

func (this *Server) Listen() *listen.Listen {
	if this.listen != nil {
		return this.listen
	}

	this.listen = &listen.Listen{}

	this.listen.Init(&listen.Config{
		Timeout:    this.timeout,
		Logger:     this.logger,
		Concurrent: this.concurrentLimit,
		MaxBytes:   this.clientMaxRecordMaxBytes,
		OnListened: func(lInfo *listen.ListeningInfo) {
			p := event.Parameters{}

			this.Event().Trigger("on.port.registered",
				p.AddString("IP", types.String(lInfo.IP.String())).
					AddInt16("Port", types.Int16(lInfo.Port)).
					AddString("Protocol", types.String(lInfo.Protocol).Lower()))
		},
		OnUnListened: func(lInfo *listen.ListeningInfo) {
			p := event.Parameters{}

			this.Event().Trigger("on.port.unregistered",
				p.AddString("IP", types.String(lInfo.IP.String())).
					AddInt16("Port", types.Int16(lInfo.Port)).
					AddString("Protocol", types.String(lInfo.Protocol).Lower()))
		},
		OnError: func(c listen.ConnectionInfo, e *types.Throw) {
			this.logger.Debugf("An error happened when '%s' "+
				"connected to '%s': %s", c.ClientIP.String(),
				c.ServerAddress.IP, e)
		},
		OnPick: func(c listen.ConnectionInfo, r listen.RespondedResult) {
			switch r.Suggestion {
			case listen.RESPOND_SUGGEST_SKIP:
				return

			case listen.RESPOND_SUGGEST_MARK:

			default:
				this.logger.Errorf("Unknown respond suggestion"+
					" type '%d', abort", r.Suggestion)

				return
			}

			this.serverDownWait.Add(1)

			this.clientRWLock.RoutineExec(func() {
				defer this.serverDownWait.Done()

				this.bumpClient(c, r)
			})
		},
	})

	this.logger.Debugf("`Listen` module now initialized")

	return this.listen
}

func (this *Server) Event() *event.Event {
	if this.event != nil {
		return this.event
	}

	this.event = &event.Event{}

	this.event.Init(&event.Config{
		Logger: this.logger,
	})

	this.logger.Debugf("`Event` module now initialized")

	return this.event
}

func (this *Server) clients() *client.Clients {
	if this.clientMaps != nil {
		return this.clientMaps
	}

	this.clientMaps = client.NewClients(client.Config{
		OnMark: func(c *client.Client, typ client.MarkType) {
			p := event.Parameters{}

			this.Event().Trigger("on.client.marked",
				p.AddString("ClientIP", types.String(
					c.Address().String())).
					AddUInt32("Count", c.Count()))

			switch typ {
			case client.CLIENT_MARK_MANUAL:
				fallthrough
			case client.CLIENT_MARK_PICK:
				lRecord := c.LastRecord()

				clientInfo := server.ClientInfo{
					Client: types.ConvertIP(c.Address()),
					Marked: true,
				}

				if lRecord != nil {
					clientInfo.Server = lRecord.Hitting.IPAddress
					clientInfo.Type = lRecord.Hitting.Type
				}

				for _, markCmd := range this.onMarkCommands {
					markCmd(clientInfo)
				}
			}
		},
		OnUnmark: func(c *client.Client, typ client.UnmarkType) {
			p := event.Parameters{}

			this.Event().Trigger("on.client.marked.out",
				p.AddString("ClientIP", types.String(
					c.Address().String())).
					AddUInt32("Count", c.Count()))

			switch typ {
			case client.CLIENT_UNMARK_MANUAL:
				for _, unmarkCmd := range this.onUnmarkCommands {
					unmarkCmd(types.ConvertIP(c.Address()))
				}
			}
		},
		OnRecord: func(client *client.Client, data client.Record) {
			p := event.Parameters{}

			this.Event().Trigger("on.client.hitting",
				p.AddString("ClientIP", types.String(
					client.Address().String())).
					AddString("ServerIP", types.String(
						data.Hitting.IP.String())).
					AddUInt16("ServerPort", data.Hitting.Port).
					AddString("Type", data.Hitting.Type).
					AddBytes("ReceivedSample", data.Inbound).
					AddBytes("RespondedData", data.Outbound))
		},
	})

	return this.clientMaps
}

func (this *Server) insertClient(c listen.ConnectionInfo, mark bool,
	insertType client.MarkType) (*client.Client, *types.Throw) {
	nowTime := time.Now()

	clientRecord, newClientRec := this.clients().Get(c.ClientIP)

	historyRecord := this.history.GetSlot(
		this.bootTime)

	// If this is a new client, add inbound record
	if newClientRec {
		this.totalInbound += 1
		historyRecord.Inbound += 1

		clientRecord.Tolerate(this.tolerate, this.tolerateExpire,
			this.tolerateRestrict)
	}

	// Don't plus hit as we don't have actual hit here
	this.totalHit += 1
	historyRecord.Hit += 1

	// Update port distribution
	portDis := this.distribution.GetSlot(
		c.ServerAddress.Port, c.Type)

	portDis.Hit += 1

	clientRecord.Record(client.Record{
		Inbound:  []byte{},
		Outbound: []byte{},
		Hitting: client.Hitting{
			IPAddress: c.ServerAddress,
			Type:      c.Type,
		},
		Time: nowTime,
	}, this.clientMaxRecords)

	clientRecord.Bump()

	if mark {
		clientRecord.Mark(insertType)

		this.totalMarked += 1
		historyRecord.Marked += 1
	}

	this.logger.Infof("Client '%s' has been manually added as "+
		"it connects '%s:%d'", c.ClientIP.String(),
		c.ServerAddress.IP.IP(),
		c.ServerAddress.Port)

	return clientRecord, nil
}

func (this *Server) bumpClient(c listen.ConnectionInfo,
	r listen.RespondedResult) (*client.Client, *types.Throw) {
	nowTime := time.Now()

	clientRecord, newClientRec := this.clients().Get(c.ClientIP)

	historyRecord := this.history.GetSlot(
		this.bootTime)

	if newClientRec {
		this.totalInbound += 1
		historyRecord.Inbound += 1

		clientRecord.Tolerate(this.tolerate, this.tolerateExpire,
			this.tolerateRestrict)
	}

	this.totalHit += 1
	historyRecord.Hit += 1

	portDis := this.distribution.GetSlot(
		c.ServerAddress.Port,
		c.Type)

	portDis.Hit += 1

	clientRecord.Record(client.Record{
		Inbound: func(r *listen.RespondedResult) []byte {
			if types.Int32(len(r.ReceivedSample)).UInt32() >
				this.clientMaxRecordMaxBytes {
				return r.ReceivedSample[:this.clientMaxRecordMaxBytes]
			}

			return r.ReceivedSample
		}(&r),
		Outbound: func(r *listen.RespondedResult) []byte {
			if types.Int32(len(r.RespondedData)).UInt32() >
				this.clientMaxRecordMaxBytes {
				return r.RespondedData[:this.clientMaxRecordMaxBytes]
			}

			return r.RespondedData
		}(&r),
		Hitting: client.Hitting{
			IPAddress: c.ServerAddress,
			Type:      c.Type,
		},
		Time: nowTime,
	}, this.clientMaxRecords)

	// Check expiration here, allowing faster expire reset
	if clientRecord.Expired(nowTime) == client.CLIENT_EXPIRED_YES {
		clientRecord.Rebump() // Reset count
	} else {
		clientRecord.Bump() // Update count and last seen
	}

	// Check the connection tolerate limit
	if clientRecord.Count() < this.tolerate {
		this.logger.Infof("Client '%s' connected '%d'"+
			" times, still counting", clientRecord.Address(),
			clientRecord.Count())

		return clientRecord, nil
	}

	if clientRecord.Marked() {
		this.logger.Infof("Marked client '%s' comes again, "+
			"counting", clientRecord.Address())

		return clientRecord, nil
	}

	clientRecord.Mark(client.CLIENT_MARK_PICK)

	this.totalMarked += 1
	historyRecord.Marked += 1

	this.logger.Infof("Client '%s' worth some notice as it "+
		"connected us '%d' times within '%s'", clientRecord.Address(),
		clientRecord.Count(), this.tolerateExpire)

	return clientRecord, nil
}

func (this *Server) clientCron() {
	this.serverDownWait.Add(1)

	defer this.serverDownWait.Done()

	for {
		nowTime := time.Now()

		select {
		case <-this.clientCronExitCh:
			return

		case <-time.After(64 * time.Second):
			this.clientRWLock.Exec(func() {
				this.clients().Scan(func(clientID types.IP,
					clientInfo *client.Client) *types.Throw {
					switch clientInfo.Expired(nowTime) {
					case client.CLIENT_EXPIRED_NO:
						return nil

					case client.CLIENT_EXPIRED_RESTRICTED:
						clientInfo.Unmark(client.CLIENT_UNMARK_EXPIRE)
						return nil

					case client.CLIENT_EXPIRED_YES:
						return this.clients().Delete(clientID,
							client.CLIENT_UNMARK_EXPIRE)
					}

					return nil
				})
			})
		}
	}
}

func (this *Server) Clients() []client.ClientExport {
	var clients []client.ClientExport

	this.clientRWLock.Exec(func() {
		clients = this.clients().Export()
	})

	return clients
}

func (this *Server) Client(addr types.IP) (*client.Client, *types.Throw) {
	var c *client.Client = nil
	var e *types.Throw = nil

	this.clientRWLock.Exec(func() {
		if !this.clients().Has(addr) {
			e = server.ErrClientNotFound.Throw(addr.IP())

			return
		}

		c, _ = this.clients().Get(addr)
	})

	if e != nil {
		return nil, e
	}

	if c == nil {
		return nil, server.ErrClientNotFound.Throw(addr.IP())
	}

	// Make a copy, all change to the client must go through manageable methods
	clientCopy := &client.Client{}

	*clientCopy = *c

	return clientCopy, nil
}

func (this *Server) addClient(clientData server.ClientInfo,
	addType client.MarkType) (*client.Client, *types.Throw) {
	var c *client.Client = nil
	var e *types.Throw = nil

	if clientData.Client.IsEmpty() {
		return nil, server.ErrInvalidClientAddress.Throw(clientData.Client.IP())
	}

	if clientData.Server.IsEmpty() {
		return nil, server.ErrInvalidServerAddress.Throw(
			clientData.Server.IP, clientData.Server.Port)
	}

	if clientData.Type == "" {
		return nil, server.ErrInvalidConnectionType.Throw(
			clientData.Client.IP())
	}

	this.clientRWLock.Exec(func() {
		// Check if it already existed
		if this.clients().Has(clientData.Client) {
			e = server.ErrClientAlreadyExisted.Throw(
				clientData.Client.IP())

			return
		}

		// Add client to data set
		newClient, newClientErr := this.insertClient(listen.ConnectionInfo{
			ClientIP:      clientData.Client,
			ServerAddress: clientData.Server,
			Type:          clientData.Type,
		}, clientData.Marked, addType)

		if newClientErr != nil {
			e = newClientErr

			return
		}

		c = newClient
	})

	return c, e
}

func (this *Server) AddClient(
	clientData server.ClientInfo) (*client.Client, *types.Throw) {
	return this.addClient(clientData, client.CLIENT_MARK_MANUAL)
}

func (this *Server) ImportClient(
	clientData server.ClientInfo) (*client.Client, *types.Throw) {
	return this.addClient(clientData, client.CLIENT_MARK_OTHER)
}

func (this *Server) RemoveClient(addr types.IP) *types.Throw {
	var result *types.Throw = nil

	this.clientRWLock.Exec(func() {
		if !this.clients().Has(addr) {
			result = server.ErrClientNotFound.Throw(addr)

			return
		}

		result = this.clients().Delete(addr, client.CLIENT_UNMARK_MANUAL)
	})

	if result != nil {
		return result
	}

	return result
}

func (this *Server) Status() server.Status {
	sInfo := server.Status{}

	this.clientRWLock.Exec(func() {
		sInfo.Uptime = time.Now().Sub(this.bootTime)

		sInfo.History = this.history.Histories()
		sInfo.Distribution = this.distribution.Distributions()

		sInfo.TotalInbound = this.totalInbound
		sInfo.TotalMarked = this.totalMarked
		sInfo.TotalHit = this.totalHit
		sInfo.TotalClients = types.UInt64(this.clients().Len())
	})

	return sInfo
}

func (this *Server) powerup() *types.Throw {
	if this.serverUpped {
		return server.ErrServerAlreadyUp.Throw()
	}

	if this.serverUpping {
		return server.ErrServerIsBooting.Throw()
	}

	this.logger.Debugf("Powering up")

	this.serverUpping = true
	this.bootTime = time.Now()
	this.clientCronExitCh = make(chan bool)

	go this.clientCron()

	lnErr := this.Listen().Serv()

	// Send the up command after server is up
	for _, upCmd := range this.onUpCommands {
		this.logger.Debugf("Firing `Up` command")

		e := upCmd()

		if e != nil {
			this.logger.Errorf("The last `Up` command throws an error: %s", e)
		}
	}

	// Send the up part of the UpDown commands
	for _, udCmd := range this.onUpDownCommands {
		this.logger.Debugf("Firing `Up` part of an `UpDown` command")

		e := udCmd.Alpha()

		if e != nil {
			this.logger.Errorf("The last `UpDown` command throws an error: %s",
				e)
		}
	}

	// We already up once we get here
	this.serverUpped = true
	this.serverUpping = false

	this.Event().Trigger("on.server.up", event.Parameters{})

	if lnErr != nil {
		this.logger.Errorf("`Server` started, but there is at least "+
			"one problem: %s", lnErr)

		return lnErr
	}

	this.logger.Infof("`Server` started without any serious problem")

	return nil
}

func (this *Server) shutdown() *types.Throw {
	if !this.serverUpped {
		return server.ErrServerNotYetStarted.Throw()
	}

	this.logger.Debugf("Shutting down")

	// Unmark all clients before shutdown
	this.clientRWLock.Exec(func() {
		this.clients().Clear()
	})

	// Send down commands before actually down the server
	for _, downCmd := range this.onDownCommands {
		this.logger.Debugf("Firing `Down` command")

		e := downCmd()

		if e != nil {
			this.logger.Errorf("The last `Down` command throws an error: %s",
				e)
		}
	}

	// Send the down part of the UpDown commands
	for downCmdLoop := len(
		this.onUpDownCommands) - 1; downCmdLoop >= 0; downCmdLoop-- {
		this.logger.Debugf("Firing `Down` part of an `UpDown` command")

		e := this.onUpDownCommands[downCmdLoop].Beta()

		if e != nil {
			this.logger.Errorf("The last `UpDown` command "+
				"throws an error: %s", e)
		}
	}

	// Down the server
	lnErr := this.Listen().Down()

	this.clientCronExitCh <- true

	// Final wait
	this.serverDownWait.Wait()

	// We already successfully shut server down once we get here
	this.serverUpped = false

	this.Event().Trigger("on.server.down", event.Parameters{})

	if lnErr != nil {
		this.logger.Errorf("`Server` is down, but there is at least "+
			"one problem: %s", lnErr)

		return lnErr
	}

	this.logger.Infof("`Server` is successfully down")

	return nil
}

func (this *Server) Serv() *types.Throw {
	var lnErr *types.Throw = nil

	this.serverLock.Exec(func() {
		lnErr = this.powerup()
	})

	return lnErr
}

func (this *Server) Down() *types.Throw {
	var lnErr *types.Throw = nil

	this.serverLock.Exec(func() {
		lnErr = this.shutdown()
	})

	return lnErr
}

func (this *Server) Reload(
	callback func(s *Server) *types.Throw) *types.Throw {
	var lnErr *types.Throw = nil

	this.logger.Debugf("Reloading")

	this.serverLock.Exec(func() {
		lnErr = this.shutdown()

		if lnErr != nil {
			return
		}

		this.listen = nil
		this.event = nil

		this.clientMaps = &client.Clients{}

		this.onUpCommands = types.Callbacks{}
		this.onDownCommands = types.Callbacks{}
		this.onUpDownCommands = []types.CallbackPair{}
		this.onMarkCommands = []func(server.ClientInfo){}
		this.onUnmarkCommands = []func(types.IP){}

		this.totalInbound = 0
		this.totalMarked = 0
		this.totalHit = 0
		this.history = server.Histories{}
		this.distribution = server.Distributions{}

		this.logger.Debugf("Running `Reload` callback")

		lnErr = callback(this)

		if lnErr != nil {
			return
		}

		lnErr = this.powerup()
	})

	if lnErr != nil {
		this.logger.Errorf("There at least one problem when "+
			"reloading `Server`: %s", lnErr)

		return lnErr
	}

	this.logger.Infof("`Server` is successfully reloaded")

	return lnErr
}

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
	"github.com/raincious/trap/trap/core/logger"
	"github.com/raincious/trap/trap/core/server"
	"github.com/raincious/trap/trap/core/sync"
	"github.com/raincious/trap/trap/core/sync/communication"
	"github.com/raincious/trap/trap/core/sync/communication/conn"
	"github.com/raincious/trap/trap/core/sync/communication/controller"
	"github.com/raincious/trap/trap/core/sync/communication/messager"
	"github.com/raincious/trap/trap/core/types"

	"crypto/tls"
	"net"
	"time"
)

type Sync struct {
	listenOn          net.TCPAddr
	trapServer        *Server
	tlsCert           tls.Certificate
	maxReceiveLen     types.UInt16
	logger            *logger.Logger
	passphrase        types.String
	syncNodes         *sync.Nodes
	syncServer        *communication.Server
	syncRetry         uint16
	activeClients     *sync.ActiveClientsTable
	requestTimeout    time.Duration
	connectionTimeout time.Duration
	looseTimeout      time.Duration
	cronDownChan      chan bool
	downing           bool
}

func NewSync() *Sync {
	return &Sync{
		trapServer: nil,
		listenOn: net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 0,
		},
		tlsCert:           tls.Certificate{},
		maxReceiveLen:     65535,
		passphrase:        "",
		syncNodes:         nil,
		syncServer:        nil,
		syncRetry:         3,
		activeClients:     sync.NewActiveClientsTable(),
		requestTimeout:    6 * time.Second,
		connectionTimeout: 120 * time.Second,
		looseTimeout:      120 * time.Second,
		cronDownChan:      make(chan bool),
		downing:           false,
	}
}

func (s *Sync) nodes() *sync.Nodes {
	if s.syncNodes != nil {
		return s.syncNodes
	}

	client := controller.Client{
		Common: controller.Common{
			GetPartners: func() types.SearchableIPAddresses {
				return s.nodes().Partners()
			},
			IsAuthed: func(clientAddr net.Addr) bool {
				return true
			},
			MarkClients: func(
				c *conn.Conn,
				clients []server.ClientInfo,
			) *types.Throw {
				importedClients := []server.ClientInfo{}

				for _, client := range clients {
					if client.Server.IP.IsZero() {
						newIP, ipErr := types.ConvertIPAddress(c.RemoteAddr())

						if ipErr != nil {
							continue
						}

						client.Server.IP = newIP.IP
					}

					s.trapServer.ImportClient(client)

					importedClients = append(importedClients, client)
				}

				go s.nodes().BroadcastMarkClients([]*conn.Conn{c},
					importedClients, s.syncRetry)

				go s.server().BroadcastMarkClients([]*conn.Conn{c},
					importedClients, s.syncRetry)

				return nil
			},
			UnmarkClients: func(
				c *conn.Conn,
				clients []types.IP,
			) *types.Throw {
				for _, client := range clients {
					s.trapServer.RemoveClient(client)
				}

				go s.nodes().BroadcastUnmarkClients([]*conn.Conn{c},
					clients, s.syncRetry)

				go s.server().BroadcastUnmarkClients([]*conn.Conn{c},
					clients, s.syncRetry)

				return nil
			},
		},
		AddPartners: func(
			c *conn.Conn, ips types.SearchableIPAddresses) *types.Throw {
			err := s.server().BroadcastNewPartners([]*conn.Conn{},
				ips.Export(), s.syncRetry)

			if err != nil {
				s.logger.Debugf("Can't broadcast `AddPartners` "+
					"information due to error: %s", err)
			}

			return nil
		},
		RemovePartners: func(
			c *conn.Conn, ips types.SearchableIPAddresses) *types.Throw {
			err := s.server().BroadcastDetachedPartners(
				[]*conn.Conn{}, ips.Export(), s.syncRetry)

			if err != nil {
				s.logger.Debugf("Can't broadcast `RemovePartners` "+
					"information due to error: %s", err)
			}

			return nil
		},
	}

	s.syncNodes = sync.NewNodes(
		client,
		s.logger,
		s.maxReceiveLen,
		s.requestTimeout,
		s.connectionTimeout,
	)

	return s.syncNodes
}

func (s *Sync) server() *communication.Server {
	if s.syncServer != nil {
		return s.syncServer
	}

	contrl := controller.Server{
		Common: controller.Common{
			GetPartners: func() types.SearchableIPAddresses {
				return s.nodes().Partners()
			},
			IsAuthed: func(clientAddr net.Addr) bool {
				return s.activeClients.HasAddr(clientAddr)
			},
			MarkClients: func(
				c *conn.Conn,
				clients []server.ClientInfo,
			) *types.Throw {
				importedClients := []server.ClientInfo{}

				for _, client := range clients {
					if client.Server.IP.IsZero() {
						newIP, ipErr := types.ConvertIPAddress(c.RemoteAddr())

						if ipErr != nil {
							continue
						}

						client.Server.IP = newIP.IP
					}

					s.trapServer.ImportClient(client)

					importedClients = append(importedClients, client)
				}

				go s.nodes().BroadcastMarkClients([]*conn.Conn{c},
					importedClients, s.syncRetry)

				go s.server().BroadcastMarkClients([]*conn.Conn{c},
					importedClients, s.syncRetry)

				return nil
			},
			UnmarkClients: func(
				c *conn.Conn,
				clients []types.IP,
			) *types.Throw {
				for _, client := range clients {
					s.trapServer.RemoveClient(client)
				}

				go s.nodes().BroadcastUnmarkClients([]*conn.Conn{c},
					clients, s.syncRetry)

				go s.server().BroadcastUnmarkClients([]*conn.Conn{c},
					clients, s.syncRetry)

				return nil
			},
		},
		OnAuthed: func(connection *conn.Conn) {
			s.activeClients.Add(connection)
		},
		OnAuthFailed: func(net.Addr) {
			// Call when auth is failed
		},
		GetPassphrase: func() types.String {
			return s.passphrase
		},
		GetLooseTimeout: func() time.Duration {
			if s.looseTimeout < s.connectionTimeout {
				return s.connectionTimeout
			}

			return s.looseTimeout
		},
	}

	handle := messager.Callbacks{}

	handle.Register(messager.SYNC_SIGNAL_HELLO, contrl.Auth)
	handle.Register(messager.SYNC_SIGNAL_HEATBEAT, contrl.Heatbeat)
	handle.Register(messager.SYNC_SIGNAL_CLIENT_MARK, contrl.ClientsMarked)
	handle.Register(messager.SYNC_SIGNAL_CLIENT_UNMARK, contrl.ClientsUnmarked)

	comServer := &communication.Server{
		OnConnected: func(conn *conn.Conn) {},
		OnDisconnected: func(conn *conn.Conn) {
			s.activeClients.Remove(conn)
		},
		Responders:         handle,
		MaxReceiveDataSize: s.maxReceiveLen,
		Logger:             s.logger,
	}

	s.syncServer = comServer

	return s.syncServer
}

func (s *Sync) SetServer(srv *Server) {
	s.trapServer = srv
}

func (s *Sync) SetLogger(l *logger.Logger) {
	s.logger = l.NewContext("Sync")
}

func (s *Sync) SetMaxReceiveLen(maxReceiveLen types.UInt16) {
	s.maxReceiveLen = maxReceiveLen

	s.logger.Debugf("Max receive length has been set to '%d'", s.maxReceiveLen)
}

func (s *Sync) SetRequestTimeout(timeout time.Duration) {
	s.requestTimeout = timeout

	s.logger.Debugf("Request timeout has been set to '%s'", s.requestTimeout)
}

func (s *Sync) SetConnectionTimeout(timeout time.Duration) {
	s.connectionTimeout = timeout

	s.logger.Debugf("Connection timeout has been set to '%s'",
		s.connectionTimeout)
}

func (s *Sync) SetLooseTimeout(timeout time.Duration) {
	s.looseTimeout = timeout

	s.logger.Debugf("Loose timeout has been set to '%s'",
		s.looseTimeout)
}

func (s *Sync) SetPort(port types.UInt16) {
	s.listenOn.Port = int(port.UInt16())

	s.logger.Debugf("Server port has been set to '%d'",
		s.listenOn.Port)
}

func (s *Sync) SetInterface(ifaceIP types.IP) {
	s.listenOn.IP = ifaceIP.IP()

	s.logger.Debugf("Server interface has been set to '%s'",
		s.listenOn.IP)
}

func (s *Sync) SetPassphrase(passphrase types.String) {
	s.passphrase = passphrase
}

func (s *Sync) LoadCert(pem types.String, key types.String) *types.Throw {
	tlsCert, tctErr := tls.LoadX509KeyPair(pem.String(), key.String())

	if tctErr != nil {
		s.logger.Errorf("Can't load certificate due to error: %s", tctErr)

		return types.ConvertError(tctErr)
	}

	s.tlsCert = tlsCert

	s.logger.Debugf("Certificate is loaded")

	return nil
}

func (s *Sync) cron() {
	for {
		select {
		case <-s.cronDownChan:
			return

		case <-time.After(10 * time.Second):
			s.tryConnectToAllNodes()
		}
	}
}

func (s *Sync) connectAllNodes() {
	s.nodes().Scan(func(key types.String, node *sync.Node) *types.Throw {
		if node.IsConnected() {
			return nil
		}

		if !node.IsReconnectable() {
			return nil
		}

		s.logger.Debugf("Connecting to node '%s'", node.Address().String())

		connectErr := node.Connect(s.nodes().Partners(),
			func(c *conn.Conn) {
				s.logger.Debugf("Node '%s' is connected",
					node.Address().String())
			},
			func(c *conn.Conn, ips types.IPAddresses) {
				defer s.server().BroadcastNewPartners([]*conn.Conn{}, ips,
					s.syncRetry)

				s.logger.Infof("Logged in to node '%s'",
					node.Address().String())
			},
			func(rmPartners types.IPAddresses, c *conn.Conn, err *types.Throw) {
				defer s.server().BroadcastDetachedPartners([]*conn.Conn{},
					rmPartners, s.syncRetry)

				if err != nil && !err.Is(messager.ErrMessageEOFReached) {
					s.logger.Warningf("Node '%s' is dropped due to error: %s",
						node.Address().String(), err)

					return
				}

				s.logger.Infof("Node '%s' is disconnected",
					node.Address().String())
			})

		if connectErr != nil {
			s.logger.Errorf("Can't connect to node '%s' due to error: %s",
				node.Address().String(), connectErr)
		}

		return nil
	})
}

func (s *Sync) disconnectAllNodes() {
	s.nodes().Clear()
}

func (s *Sync) tryConnectToAllNodes() {
	if s.downing {
		return
	}

	s.connectAllNodes()
}

func (s *Sync) AddNode(nodeAddr types.IPAddress,
	passphrase types.String) *types.Throw {
	err := s.nodes().Register(nodeAddr, passphrase)

	if err != nil {
		s.logger.Errorf("Can't add node '%s' due to error: %s",
			nodeAddr.String(), err)
	} else {
		s.logger.Debugf("Node '%s' has been added", nodeAddr.String())
	}

	return err
}

func (s *Sync) Status() sync.Status {
	status := sync.Status{
		Nodes: []sync.NodeInfo{},
		Server: sync.ServerInfo{
			Listen: types.IPAddress{
				IP:   types.ConvertIP(s.listenOn.IP),
				Port: types.UInt16(s.listenOn.Port),
			},
			Clients: []sync.ClientInfo{},
		},
	}

	s.nodes().Scan(func(key types.String, n *sync.Node) *types.Throw {
		nPartner := n.Partners()

		status.Nodes = append(status.Nodes, sync.NodeInfo{
			Address:   n.Address(),
			Delay:     n.Delay(),
			Stats:     n.Stats(),
			Connected: n.IsConnected(),
			Partner:   nPartner.Export(),
		})

		return nil
	})

	s.server().Scan(
		[]*conn.Conn{},
		func(key string, sess *communication.Session) *types.Throw {
			ipAddr, ipErr := types.ConvertIPAddress(sess.Conn().RemoteAddr())

			if ipErr != nil {
				ipAddr = types.IPAddress{}
			}

			status.Server.Clients = append(
				status.Server.Clients,
				sync.ClientInfo{
					Remote: ipAddr,
					Stats:  sess.Request().Stats(),
				})

			return nil
		})

	return status
}

func (s *Sync) Serv() *types.Throw {
	s.logger.Debugf("Booting up")

	s.trapServer.OnMark(func(client server.ClientInfo) {
		clients := []server.ClientInfo{client}

		go s.nodes().BroadcastMarkClients([]*conn.Conn{},
			clients, s.syncRetry)

		go s.server().BroadcastMarkClients([]*conn.Conn{},
			clients, s.syncRetry)
	})

	s.trapServer.OnUnmark(func(client types.IP) {
		clients := []types.IP{client}

		go s.nodes().BroadcastUnmarkClients([]*conn.Conn{},
			clients, s.syncRetry)

		go s.server().BroadcastUnmarkClients([]*conn.Conn{},
			clients, s.syncRetry)
	})

	sErr := s.server().Listen(
		s.listenOn,
		s.tlsCert,
		s.connectionTimeout,
	)

	if sErr != nil {
		s.logger.Debugf("Can't serve due to error: %s", sErr)

		return sErr
	}

	s.logger.Infof("`Sync` Server is serving at '%s'", s.listenOn.String())

	s.downing = false

	s.connectAllNodes()

	go s.cron()

	s.logger.Debugf("`Sync` is up")

	return nil
}

func (s *Sync) Down() *types.Throw {
	s.downing = true

	s.logger.Debugf("Disconnect from nodes")

	s.disconnectAllNodes()

	s.logger.Debugf("Shutting down server")

	sErr := s.server().Down()

	if sErr != nil {
		return sErr
	}

	s.logger.Debugf("`Sync` is down")

	return nil
}

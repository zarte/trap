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

package communication

import (
	"github.com/raincious/trap/trap/core/sync/communication/conn"
	"github.com/raincious/trap/trap/core/sync/communication/messager"
	"github.com/raincious/trap/trap/core/types"

	"crypto/tls"
	"sync"
	"time"
)

var (
	ErrClientAlreadyConnected *types.Error = types.NewError(
		"`Sync` Client already connected host '%s'")

	ErrClientNotConnected *types.Error = types.NewError(
		"`Sync` Client is not connected")

	ErrClientSessionUnavailable *types.Error = types.NewError(
		"Can't get session as it's unavailable")
)

type Client struct {
	Common

	upped             bool
	connected         bool
	connectedLock     types.Mutex
	defaultResponders messager.Callbacks
	connTimeout       time.Duration
	reqTimeout        time.Duration
	connDelay         time.Duration
	heatbeatPeriod    time.Duration
	delayLock         types.Mutex
	session           *Session
	sessions          *Sessions
	sessionWait       sync.WaitGroup
	sessionLock       types.Mutex
	shutdownSignal    chan bool
	shutdownLock      types.Mutex
	heatbeatTicker    <-chan time.Time
	clientAuthed      bool
}

func NewClient(defaultResponders messager.Callbacks, connTimeout time.Duration,
	reqTimeout time.Duration) *Client {
	return &Client{
		upped:             false,
		connected:         false,
		connectedLock:     types.Mutex{},
		defaultResponders: defaultResponders,
		connTimeout:       connTimeout,
		reqTimeout:        reqTimeout,
		connDelay:         time.Duration(0),
		heatbeatPeriod:    time.Duration(0),
		delayLock:         types.Mutex{},
		session:           nil,
		sessions:          nil,
		sessionWait:       sync.WaitGroup{},
		sessionLock:       types.Mutex{},
		shutdownSignal:    make(chan bool),
		shutdownLock:      types.Mutex{},
		clientAuthed:      false,
	}
}

func (c *Client) dialup(ip types.IPAddress, onConnected func(*conn.Conn),
	onDisconnected func(*conn.Conn, *types.Throw)) *types.Throw {
	if c.upped {
		return ErrClientAlreadyConnected.Throw(ip.String())
	}

	tlsConn, tlsConnErr := tls.Dial("tcp", ip.String().String(),
		&tls.Config{
			//Certificates:           []tls.Certificate{
			//                            cert
			//                        },
			InsecureSkipVerify: true,
			ServerName:         ip.String().String(),
		})

	if tlsConnErr != nil {
		return types.ConvertError(tlsConnErr)
	}

	defer func() {
		if c.connected {
			return
		}

		c.hungup()
	}()

	syncConn := &conn.Conn{
		Conn: tlsConn,
	}

	syncConn.SetTimeout(c.connTimeout)

	ready := make(chan bool)

	c.clientAuthed = false
	c.heatbeatTicker = time.Tick(c.connTimeout)
	c.connected = false
	c.sessions = NewSessions(c.defaultResponders, c.reqTimeout,
		func() {
			c.sessionWait.Add(1)
		}, func() {
			c.sessionWait.Done()
		})

	session, sessErr := c.sessions.Register(syncConn)

	if sessErr != nil {
		return sessErr
	}

	c.sessionLock.Exec(func() {
		c.session = session
	})

	c.sessionWait.Add(2)

	go func() {
		var serveErr *types.Throw = nil

		defer func() {
			c.connectedLock.Exec(func() {
				c.connected = false
			})

			c.sessions.Unregister(syncConn)

			c.sessionLock.Exec(func() {
				c.session = nil
			})

			onDisconnected(syncConn, serveErr)

			c.sessionWait.Done() // Must before c.Disconnect() Call or deadlock

			c.Disconnect()
		}()

		c.connectedLock.Exec(func() {
			c.connected = true
		})

		onConnected(syncConn)

		serveErr = session.Serve(ready)
	}()

	go func() {
		defer func() {
			syncConn.Close()
			c.sessions.Clear()
			c.sessionWait.Done()
		}()

		for {
			select {
			case <-c.shutdownSignal:
				return

			case <-c.heatbeatTicker:
				if !c.clientAuthed {
					continue
				}

				c.heatbeat()
			}
		}
	}()

	<-ready

	c.upped = true

	return nil
}

func (c *Client) hungup() *types.Throw {
	if !c.upped {
		return ErrClientNotConnected.Throw()
	}

	c.upped = false

	c.shutdownSignal <- true

	c.sessionWait.Wait()

	return nil
}

func (c *Client) Connected() bool {
	var connected bool = false

	c.connectedLock.Exec(func() {
		connected = c.connected
	})

	return connected
}

func (c *Client) Connect(ip types.IPAddress, onConnected func(*conn.Conn),
	onDisconnect func(*conn.Conn, *types.Throw)) *types.Throw {
	var err *types.Throw = nil

	c.shutdownLock.Exec(func() {
		err = c.dialup(ip, onConnected, onDisconnect)
	})

	return err
}

func (c *Client) Disconnect() *types.Throw {
	var err *types.Throw = nil

	c.shutdownLock.Exec(func() {
		err = c.hungup()
	})

	return err
}

func (c *Client) Delay() time.Duration {
	delay := time.Duration(0)

	c.delayLock.Exec(func() {
		delay = c.connDelay
	})

	return delay
}

func (c *Client) getSession() (*Session, *types.Throw) {
	if c.session == nil {
		return nil, ErrClientSessionUnavailable.Throw()
	}

	return c.session, nil
}

func (c *Client) heatbeat() *types.Throw {
	var err *types.Throw = nil

	c.sessionLock.Exec(func() {
		sess, sessErr := c.getSession()

		if sessErr != nil {
			err = sessErr

			return
		}

		delay, hbErr := sess.Heatbeat()

		if hbErr != nil {
			err = hbErr

			return
		}

		c.delayLock.Exec(func() {
			c.connDelay = delay
		})
	})

	return err
}

func (c *Client) Auth(password types.String,
	connects types.IPAddresses) (types.IPAddresses, *types.Throw) {
	var err *types.Throw = nil
	var serverPartners types.IPAddresses = nil

	c.sessionLock.Exec(func() {
		sess, sessErr := c.getSession()

		if sessErr != nil {
			err = sessErr

			return
		}

		heatb, ips, authErr := sess.Auth(password, connects)

		if authErr != nil {
			err = authErr

			return
		}

		c.delayLock.Exec(func() {
			c.heatbeatTicker = time.Tick(heatb)

			c.clientAuthed = true
		})

		serverPartners = ips
	})

	return serverPartners, err
}

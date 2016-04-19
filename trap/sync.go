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
	tlsCert           tls.Certificate
	logger            *logger.Logger
	passphrase        types.String
	syncNodes         *sync.Nodes
	syncServer        *communication.Server
	activeClients     map[string]bool
	requestTimeout    time.Duration
	connectionTimeout time.Duration
	looseTimeout      time.Duration
	cronDownChan      chan bool
	downing           bool
}

func NewSync() *Sync {
	return &Sync{
		listenOn: net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 0,
		},
		tlsCert:           tls.Certificate{},
		passphrase:        "",
		syncNodes:         nil,
		syncServer:        nil,
		activeClients:     map[string]bool{},
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
			GetPartners: func() types.IPAddresses {
				return types.IPAddresses{}
			},
			IsAuthed: func(clientAddr net.Addr) bool {
				return true
			},
			MarkClients: func(c []server.ClientInfo) *types.Throw {
				return nil
			},
			UnmarkClients: func(c []server.ClientInfo) *types.Throw {
				return nil
			},
		},
		AddPartners: func(ips types.IPAddresses) *types.Throw {
			return nil
		},
		RemovePartners: func(ips types.IPAddresses) *types.Throw {
			return nil
		},
	}
	handle := messager.Callbacks{}

	handle.Register(messager.SYNC_SIGNAL_PARTNER_ADD, client.PartnersAdded)
	handle.Register(
		messager.SYNC_SIGNAL_PARTNER_REMOVE, client.PartnersRemoved)
	handle.Register(messager.SYNC_SIGNAL_CLIENT_MARK, client.ClientsMarked)
	handle.Register(messager.SYNC_SIGNAL_CLIENT_UNMARK, client.ClientsUnmarked)

	s.syncNodes = sync.NewNodes(
		handle,
		s.requestTimeout,
		s.connectionTimeout,
	)

	return s.syncNodes
}

func (s *Sync) startServer() (*communication.Server, *types.Throw) {
	if s.syncServer != nil {
		return s.syncServer, nil
	}

	contrl := controller.Server{
		Common: controller.Common{
			GetPartners: func() types.IPAddresses {
				return types.IPAddresses{}
			},
			IsAuthed: func(clientAddr net.Addr) bool {
				addrStr := clientAddr.String()

				if _, ok := s.activeClients[addrStr]; !ok {
					return false
				}

				if !s.activeClients[addrStr] {
					return false
				}

				return true
			},
			MarkClients: func(
				ips []server.ClientInfo,
			) *types.Throw {
				return nil
			},
			UnmarkClients: func(
				ips []server.ClientInfo,
			) *types.Throw {
				return nil
			},
		},
		OnAuthed: func(conn *conn.Conn) {
			clientAddr := conn.RemoteAddr().String()

			s.activeClients[clientAddr] = true
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

	comServer := &communication.Server{
		OnConnected: func(conn *conn.Conn) {
			clientAddr := conn.RemoteAddr().String()

			s.activeClients[clientAddr] = false
		},
		OnDisconnected: func(conn *conn.Conn) {
			delete(s.activeClients, conn.RemoteAddr().String())
		},
	}

	listenErr := comServer.Listen(
		s.listenOn,
		handle,
		s.tlsCert,
		s.connectionTimeout,
	)

	if listenErr != nil {
		return nil, listenErr
	}

	s.syncServer = comServer

	return s.syncServer, nil
}

func (s *Sync) SetLogger(l *logger.Logger) {
	s.logger = l.NewContext("Sync")
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
}

func (s *Sync) SetInterface(ifaceIP types.IP) {
	s.listenOn.IP = ifaceIP.IP()
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

		case <-time.After(5 * time.Second):
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

		connectErr := node.Connect(s.nodes().Partners(),
			func(conn *conn.Conn) {
				s.logger.Errorf("Node '%s' is connected",
					node.Address().String())
			},
			func(conn *conn.Conn, err *types.Throw) {
				if err != nil {
					s.logger.Debugf("Node '%s' is dropped due to error: %s",
						node.Address().String(), err)

					return
				}

				s.logger.Debugf("Node '%s' is disconnected",
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
	s.nodes().Scan(func(key types.String, node *sync.Node) *types.Throw {
		if !node.IsConnected() {
			return nil
		}

		dconnectErr := node.Disconnect()

		if dconnectErr != nil {
			s.logger.Errorf("Can't disconnect from node '%s' due to error: %s",
				node.Address().String(), dconnectErr)

			return nil
		}

		return nil
	})
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

func (s *Sync) Serv() *types.Throw {
	_, sErr := s.startServer()

	if sErr != nil {
		s.logger.Debugf("Can't serve due to error: %s", sErr)

		return sErr
	}

	s.logger.Debugf("`Sync` Server is serving at '%s'", s.listenOn.String())

	s.downing = false

	s.connectAllNodes()

	go s.cron()

	return nil
}

func (s *Sync) Down() *types.Throw {
	s.downing = true

	server, sErr := s.startServer()

	if sErr != nil {
		return sErr
	}

	server.Down()

	s.disconnectAllNodes()

	return nil
}

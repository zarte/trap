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
	"net"
	"sync"
	"time"
)

var (
	ErrServerAlreadyUp *types.Error = types.NewError(
		"`Sync` Server is already running at '%s'")

	ErrServerNotUp *types.Error = types.NewError(
		"`Sync` Server is not started")
)

type Server struct {
	Common

	timeout        time.Duration
	server         net.Listener
	wait           sync.WaitGroup
	callbacks      messager.Callbacks
	OnConnected    func(*conn.Conn)
	OnDisconnected func(*conn.Conn)
	sessions       *Sessions
}

func (s *Server) Listen(listenOn net.TCPAddr, responders messager.Callbacks,
	cert tls.Certificate, timeout time.Duration) *types.Throw {
	if s.server != nil {
		return ErrServerAlreadyUp.Throw(listenOn.String())
	}

	// Init variables
	s.wait = sync.WaitGroup{}
	s.timeout = timeout
	s.sessions = NewSessions(responders, timeout, func() {
		s.wait.Add(1)
	}, func() {
		s.wait.Done()
	})

	listener, lsErr := tls.Listen("tcp", listenOn.String(), &tls.Config{
		InsecureSkipVerify: true,
		Certificates: []tls.Certificate{
			cert,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	})

	if lsErr != nil {
		return types.ConvertError(lsErr)
	}

	go func() {
		if s.server != nil {
			return
		}

		listener.Close()
	}()

	s.server = listener

	s.wait.Add(1)

	go s.serve()

	return nil
}

func (s *Server) serve() {
	defer s.sessions.Clear()
	defer s.wait.Done()

	onConnect := s.OnConnected
	onDisconnect := s.OnDisconnected

	for {
		connection, conErr := s.server.Accept()

		if conErr != nil {
			return
		}

		syncConn := &conn.Conn{
			Conn: connection,
		}

		syncConn.SetTimeout(s.timeout)

		session, sessErr := s.sessions.Register(syncConn)

		if sessErr != nil {
			continue
		}

		s.wait.Add(1)

		go func(c *conn.Conn) {
			ready := make(chan bool)

			defer func() {
				s.sessions.Unregister(c)

				s.wait.Done()

				if onDisconnect != nil {
					onDisconnect(c)
				}

				c.Close()
			}()

			if onConnect != nil {
				onConnect(c)
			}

			session.Serve(ready)
		}(syncConn)
	}
}

func (s *Server) Scan(excludedConns []*conn.Conn) *types.Throw {
	return s.sessions.Scan(excludedConns,
		func(key string, sess *Session) *types.Throw {
			return nil
		})
}

func (s *Server) Down() *types.Throw {
	if s.server == nil {
		return ErrServerNotUp.Throw()
	}

	downErr := s.server.Close()

	if downErr != nil {
		return types.ConvertError(downErr)
	}

	return nil
}

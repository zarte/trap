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

package controller

import (
	"github.com/raincious/trap/trap/core/sync/communication/conn"
	"github.com/raincious/trap/trap/core/sync/communication/data"
	"github.com/raincious/trap/trap/core/sync/communication/messager"
	"github.com/raincious/trap/trap/core/types"

	"net"
	"time"
)

var (
	ErrControllerServerClientAlreadyAuthed *types.Error = types.NewError(
		"Client is already authed")

	ErrControllerServerClientAuthDenied *types.Error = types.NewError(
		"Client '%s' has been denied from auth")
)

type Server struct {
	Common

	OnAuthed        func(*conn.Conn)
	OnAuthFailed    func(net.Addr)
	GetPassphrase   func() types.String
	GetLooseTimeout func() time.Duration
}

func (s *Server) Auth(req messager.Request) *types.Throw {
	var rqErr *types.Throw = nil

	helloData := data.Hello{}

	if s.IsAuthed(req.RemoteAddr()) {
		rqErr = req.Reply(messager.SYNC_SIGNAL_HELLO_DENIED,
			&data.Undefined{})

		req.Close()

		return ErrControllerServerClientAlreadyAuthed.Throw(req.RemoteAddr())
	}

	parseErr := helloData.Parse(req.Data())

	if parseErr != nil {
		req.Reply(messager.SYNC_SIGNAL_HELLO_DENIED, &data.Undefined{})

		req.Close()

		return ErrControllerServerClientAlreadyAuthed.Throw(req.RemoteAddr())
	}

	if helloData.Password != s.GetPassphrase() {
		rqErr = req.Reply(messager.SYNC_SIGNAL_HELLO_DENIED,
			&data.Undefined{})

		s.OnAuthFailed(req.RemoteAddr())

		req.Close()

		return ErrControllerServerClientAuthDenied.Throw(req.RemoteAddr())
	}

	serverPartners := s.GetPartners()
	sConnected := helloData.Connected.Searchable()

	intersection := serverPartners.Intersection(&sConnected)

	if intersection.Len() > 0 {
		rqErr = req.Reply(messager.SYNC_SIGNAL_HELLO_CONFLICT,
			&data.HelloConflict{
				Confilct: intersection.Export(),
			})

		s.OnAuthFailed(req.RemoteAddr())

		req.Close()

		return rqErr
	}

	req.Conn().SetTimeout(s.GetLooseTimeout())
	req.SetMaxSendLength(helloData.MaxLength)

	rqErr = req.Reply(messager.SYNC_SIGNAL_HELLO_ACCEPT, &data.HelloAccept{
		MaxLength:      req.GetMaxReceiveLength(),
		HeatBeatPeriod: s.GetLooseTimeout() / 2,
		Timeout:        s.GetLooseTimeout(),
		Connected:      serverPartners.Export(),
	})

	s.OnAuthed(req.Conn())

	return rqErr
}

func (s *Server) Heatbeat(req messager.Request) *types.Throw {
	heatbeat := &data.HeatBeat{}

	if !s.IsAuthed(req.RemoteAddr()) {
		req.Reply(messager.SYNC_SIGNAL_HEATBEAT_DENIED, &data.Undefined{})

		req.Close()

		return ErrControllerServerClientNotLoggedIn.Throw(req.RemoteAddr())
	}

	parseErr := heatbeat.Parse(req.Data())

	if parseErr != nil {
		req.Reply(messager.SYNC_SIGNAL_HEATBEAT_DENIED, &data.Undefined{})

		req.Close()

		return parseErr
	}

	return req.Reply(messager.SYNC_SIGNAL_HEATBEAT, heatbeat)
}

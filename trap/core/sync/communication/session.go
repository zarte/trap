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
	"github.com/raincious/trap/trap/core/server"
	"github.com/raincious/trap/trap/core/sync/communication/conn"
	"github.com/raincious/trap/trap/core/sync/communication/data"
	"github.com/raincious/trap/trap/core/sync/communication/messager"
	"github.com/raincious/trap/trap/core/types"

	"sync"
	"time"
)

var (
	ErrSessionAlreadyRegistered *types.Error = types.NewError(
		"`Sync` Session for '%s' is already registered")

	ErrSessionNotRegistered *types.Error = types.NewError(
		"`Sync` Session for '%s' is not yet registered")

	ErrSessionAuthFailedDenied *types.Error = types.NewError(
		"Auth request has been denied by server '%s'")

	ErrSessionAuthFailedConflicted *types.Error = types.NewError(
		"Current client is conflicted with target server '%s'")

	ErrSessionHeatBeatDenied *types.Error = types.NewError(
		"'%s' refused to reply 'Heatbeat' request")

	ErrSessionPartnerAddDenied *types.Error = types.NewError(
		"'%s' refused to handle 'Partner Add' request")

	ErrSessionPartnerRemoveDenied *types.Error = types.NewError(
		"'%s' refused to handle 'Partner Remove' request")

	ErrSessionClientMarkDenied *types.Error = types.NewError(
		"'%s' refused to handle 'Client Mark' request")

	ErrSessionClientUnmarkDenied *types.Error = types.NewError(
		"'%s' refused to handle 'Client Unmark' request")
)

type Session struct {
	conn           *conn.Conn
	messager       *messager.Messager
	wait           sync.WaitGroup
	requestTimeout time.Duration
	enabled        bool
	enabledLock    types.Mutex
}

func (s *Session) Serve(serving chan bool) *types.Throw {
	var err *types.Throw = nil

	listenReady := make(chan bool)

	s.wait.Add(1)

	go func() {
		defer func() {
			s.enabledLock.Exec(func() {
				s.enabled = false
			})

			s.wait.Done()
		}()

		s.enabledLock.Exec(func() {
			s.enabled = true
		})

		err = s.messager.Listen(s.conn, listenReady)
	}()

	<-listenReady

	serving <- true

	s.wait.Wait()

	return err
}

func (s *Session) Conn() *conn.Conn {
	return s.conn
}

func (s *Session) Request() *messager.Messager {
	return s.messager
}

func (s *Session) Close() *types.Throw {
	err := s.conn.Close()

	if err != nil {
		return types.ConvertError(err)
	}

	return nil
}

func (s *Session) Enabled() bool {
	var enabled bool = false

	s.enabledLock.Exec(func() {
		if !s.enabled {
			return
		}

		enabled = true
	})

	return enabled
}

func (s *Session) registering() *types.Throw {
	s.wait = sync.WaitGroup{}

	return nil
}

func (s *Session) unregistering() *types.Throw {
	s.wait.Wait()

	return nil
}

func (s *Session) Auth(
	maxLength types.UInt16,
	password types.String,
	connects types.IPAddresses,
	onAuthed func(*conn.Conn, types.IPAddresses),
) (
	types.UInt16,
	time.Duration,
	time.Duration,
	types.IPAddresses,
	*types.Throw,
) {
	pnMaxLen := types.UInt16(0)
	heatbeatPeriod := time.Duration(0)
	newPartners := types.IPAddresses{}
	confilcteds := types.IPAddresses{}
	startTime := time.Time{}
	endTime := time.Time{}

	hello := &data.Hello{
		MaxLength: maxLength,
		Password:  password,
		Connected: connects,
	}

	handle := messager.Callbacks{}

	handle.Register(messager.SYNC_SIGNAL_HELLO_ACCEPT,
		func(req messager.Request) *types.Throw {
			accept := data.HelloAccept{}

			acceptErr := accept.Parse(req.Data())

			if acceptErr != nil {
				return acceptErr
			}

			pnMaxLen = accept.MaxLength
			heatbeatPeriod = accept.HeatBeatPeriod
			newPartners = accept.Connected

			req.Conn().SetTimeout(accept.Timeout)
			req.SetMaxSendLength(accept.MaxLength)

			endTime = time.Now()

			return nil
		})

	handle.Register(messager.SYNC_SIGNAL_HELLO_DENIED,
		func(req messager.Request) *types.Throw {
			return ErrSessionAuthFailedDenied.Throw(req.RemoteAddr())
		})

	handle.Register(messager.SYNC_SIGNAL_HELLO_CONFLICT,
		func(req messager.Request) *types.Throw {
			confilct := data.HelloConflict{}

			confilctErr := confilct.Parse(req.Data())

			if confilctErr != nil {
				return confilctErr
			}

			confilcteds = confilct.Confilct

			return ErrSessionAuthFailedConflicted.Throw(req.RemoteAddr())
		})

	startTime = time.Now()

	reqErr := s.Request().Query(
		messager.SYNC_SIGNAL_HELLO,
		hello,
		handle,
		s.requestTimeout,
	)

	if reqErr != nil {
		if reqErr.Is(ErrSessionAuthFailedConflicted) {
			return 0, time.Duration(0), time.Duration(0), confilcteds, reqErr
		}

		return 0, time.Duration(0), time.Duration(0),
			types.IPAddresses{}, reqErr
	}

	onAuthed(s.conn, newPartners)

	return pnMaxLen, endTime.Sub(startTime), heatbeatPeriod, newPartners, nil
}

func (s *Session) Heatbeat() (time.Duration, *types.Throw) {
	startTime := time.Time{}
	endTime := time.Time{}
	handle := messager.Callbacks{}
	heatbeat := &data.HeatBeat{}

	handle.Register(messager.SYNC_SIGNAL_HEATBEAT,
		func(req messager.Request) *types.Throw {
			endTime = time.Now()

			return nil
		})

	handle.Register(messager.SYNC_SIGNAL_HEATBEAT_DENIED,
		func(req messager.Request) *types.Throw {
			return ErrSessionHeatBeatDenied.Throw(req.RemoteAddr())
		})

	startTime = time.Now()

	reqErr := s.Request().Query(
		messager.SYNC_SIGNAL_HEATBEAT,
		heatbeat,
		handle,
		s.requestTimeout,
	)

	if reqErr != nil {
		return time.Duration(0), reqErr
	}

	return endTime.Sub(startTime), nil
}

func (s *Session) AddPartners(partners types.IPAddresses) *types.Throw {
	handle := messager.Callbacks{}
	partner := &data.Partner{}

	partner.Added = partners

	handle.Register(messager.SYNC_SIGNAL_PARTNER_ADD_ACCEPT,
		func(req messager.Request) *types.Throw {
			return nil
		})

	handle.Register(messager.SYNC_SIGNAL_PARTNER_ADD_DENIED,
		func(req messager.Request) *types.Throw {
			return ErrSessionPartnerAddDenied.Throw(req.RemoteAddr())
		})

	reqErr := s.Request().Query(
		messager.SYNC_SIGNAL_PARTNER_ADD,
		partner,
		handle,
		s.requestTimeout,
	)

	if reqErr != nil {
		return reqErr
	}

	return nil
}

func (s *Session) RemovePartners(partners types.IPAddresses) *types.Throw {
	handle := messager.Callbacks{}
	partner := &data.Partner{}

	partner.Removed = partners

	handle.Register(messager.SYNC_SIGNAL_PARTNER_REMOVE_ACCEPT,
		func(req messager.Request) *types.Throw {
			return nil
		})

	handle.Register(messager.SYNC_SIGNAL_PARTNER_REMOVE_DENIED,
		func(req messager.Request) *types.Throw {
			return ErrSessionPartnerRemoveDenied.Throw(req.RemoteAddr())
		})

	reqErr := s.Request().Query(
		messager.SYNC_SIGNAL_PARTNER_REMOVE,
		partner,
		handle,
		s.requestTimeout,
	)

	if reqErr != nil {
		return reqErr
	}

	return nil
}

func (s *Session) MarkClients(clients []server.ClientInfo) *types.Throw {
	handle := messager.Callbacks{}
	mark := &data.ClientMark{}

	handle.Register(messager.SYNC_SIGNAL_CLIENT_MARK_ACCEPT,
		func(req messager.Request) *types.Throw {
			return nil
		})

	handle.Register(messager.SYNC_SIGNAL_CLIENT_MARK_DENIED,
		func(req messager.Request) *types.Throw {
			return ErrSessionClientMarkDenied.Throw(req.RemoteAddr())
		})

	mark.Addresses = clients

	reqErr := s.Request().Query(
		messager.SYNC_SIGNAL_CLIENT_MARK,
		mark,
		handle,
		s.requestTimeout,
	)

	if reqErr != nil {
		return reqErr
	}

	return nil
}

func (s *Session) UnmarkClients(clients []types.IP) *types.Throw {
	handle := messager.Callbacks{}
	um := &data.ClientUnmark{}

	handle.Register(messager.SYNC_SIGNAL_CLIENT_UNMARK_ACCEPT,
		func(req messager.Request) *types.Throw {
			return nil
		})

	handle.Register(messager.SYNC_SIGNAL_CLIENT_UNMARK_DENIED,
		func(req messager.Request) *types.Throw {
			return ErrSessionClientUnmarkDenied.Throw(req.RemoteAddr())
		})

	um.Addresses = clients

	reqErr := s.Request().Query(
		messager.SYNC_SIGNAL_CLIENT_UNMARK,
		um,
		handle,
		s.requestTimeout,
	)

	if reqErr != nil {
		return reqErr
	}

	return nil
}

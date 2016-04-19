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
	"net"

	"github.com/raincious/trap/trap/core/server"
	"github.com/raincious/trap/trap/core/sync/communication/data"
	"github.com/raincious/trap/trap/core/sync/communication/messager"

	"github.com/raincious/trap/trap/core/types"
)

var (
	ErrControllerInvalidData *types.Error = types.NewError("Failed to parse data sent by '%s': %s")

	ErrControllerServerClientNotLoggedIn *types.Error = types.NewError("Client '%s' is not logged in")
)

type Common struct {
	GetPartners func() types.IPAddresses

	IsAuthed func(net.Addr) bool

	MarkClients   func([]server.ClientInfo) *types.Throw
	UnmarkClients func([]server.ClientInfo) *types.Throw
}

func (c *Common) ClientsMarked(req messager.Request) *types.Throw {
	marked := &data.ClientMark{}

	if !c.IsAuthed(req.RemoteAddr()) {
		req.Reply(messager.SYNC_SIGNAL_CLIENT_MARK_DENIED, &data.Undefined{})

		req.Close()

		return ErrControllerServerClientNotLoggedIn.Throw(req.RemoteAddr())
	}

	parseErr := marked.Parse(req.Data())

	if parseErr != nil {
		req.Reply(messager.SYNC_SIGNAL_CLIENT_MARK_DENIED, &data.Undefined{})

		req.Close()

		return ErrControllerInvalidData.Throw(req.RemoteAddr(), parseErr)
	}

	markErr := c.MarkClients(marked.Addresses)

	if markErr != nil {
		req.Reply(messager.SYNC_SIGNAL_CLIENT_MARK_DENIED, &data.Undefined{})

		req.Close()

		return markErr
	}

	return req.Reply(messager.SYNC_SIGNAL_CLIENT_MARK_ACCEPT, &data.Undefined{})
}

func (c *Common) ClientsUnmarked(req messager.Request) *types.Throw {
	unmarked := &data.ClientUnmark{}

	if !c.IsAuthed(req.RemoteAddr()) {
		req.Reply(messager.SYNC_SIGNAL_CLIENT_UNMARK_DENIED, &data.Undefined{})

		req.Close()

		return ErrControllerServerClientNotLoggedIn.Throw(req.RemoteAddr())
	}

	parseErr := unmarked.Parse(req.Data())

	if parseErr != nil {
		req.Reply(messager.SYNC_SIGNAL_CLIENT_UNMARK_DENIED, &data.Undefined{})

		req.Close()

		return ErrControllerInvalidData.Throw(req.RemoteAddr(), parseErr)
	}

	unmarkErr := c.UnmarkClients(unmarked.Addresses)

	if unmarkErr != nil {
		req.Reply(messager.SYNC_SIGNAL_CLIENT_UNMARK_DENIED, &data.Undefined{})

		req.Close()

		return unmarkErr
	}

	return req.Reply(messager.SYNC_SIGNAL_CLIENT_UNMARK_ACCEPT,
		&data.Undefined{})
}

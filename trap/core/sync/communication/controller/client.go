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
)

var (
	ErrControllerPartnerConflicted *types.Error = types.NewError(
		"We already connected at least one of the partners of '%s'")
)

type Client struct {
	Common

	AddPartners func(
		*conn.Conn, types.SearchableIPAddresses) *types.Throw
	ConfilctedPartners func(types.SearchableIPAddresses)
	RemovePartners     func(
		*conn.Conn, types.SearchableIPAddresses) *types.Throw
}

func (c *Client) PartnersAdded(req messager.Request) *types.Throw {
	partner := &data.Partner{}

	if !c.IsAuthed(req.RemoteAddr()) {
		req.Reply(messager.SYNC_SIGNAL_PARTNER_ADD_DENIED, &data.Undefined{})

		req.Close()

		return ErrControllerServerClientNotLoggedIn.Throw(req.RemoteAddr())
	}

	parseErr := partner.Parse(req.Data())

	if parseErr != nil {
		req.Reply(messager.SYNC_SIGNAL_PARTNER_ADD_DENIED, &data.Undefined{})

		req.Close()

		return ErrControllerInvalidData.Throw(req.RemoteAddr(), parseErr)
	}

	serverPartners := c.GetPartners()

	searchableNewPartners := partner.Added.Searchable()
	intersection := serverPartners.Intersection(&searchableNewPartners)

	if intersection.Len() > 0 {
		req.Reply(messager.SYNC_SIGNAL_PARTNER_ADD_DENIED, &data.Undefined{})

		req.Close()

		c.ConfilctedPartners(intersection)

		return ErrControllerPartnerConflicted.Throw(req.RemoteAddr())
	}

	addErr := c.AddPartners(req.Conn(), searchableNewPartners)

	if addErr != nil {
		req.Reply(messager.SYNC_SIGNAL_PARTNER_ADD_DENIED, &data.Undefined{})

		req.Close()

		return addErr
	}

	return req.Reply(messager.SYNC_SIGNAL_PARTNER_ADD_ACCEPT, &data.Undefined{})
}

func (c *Client) PartnersRemoved(req messager.Request) *types.Throw {
	partner := &data.Partner{}

	if !c.IsAuthed(req.RemoteAddr()) {
		req.Reply(messager.SYNC_SIGNAL_PARTNER_REMOVE_DENIED, &data.Undefined{})

		req.Close()

		return ErrControllerServerClientNotLoggedIn.Throw(req.RemoteAddr())
	}

	parseErr := partner.Parse(req.Data())

	if parseErr != nil {
		req.Reply(messager.SYNC_SIGNAL_PARTNER_REMOVE_DENIED, &data.Undefined{})

		req.Close()

		return ErrControllerInvalidData.Throw(req.RemoteAddr(), parseErr)
	}

	delErr := c.RemovePartners(req.Conn(), partner.Removed.Searchable())

	if delErr != nil {
		req.Reply(messager.SYNC_SIGNAL_PARTNER_ADD_DENIED, &data.Undefined{})

		req.Close()

		return delErr
	}

	return req.Reply(messager.SYNC_SIGNAL_PARTNER_REMOVE_ACCEPT,
		&data.Undefined{})
}

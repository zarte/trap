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

package client

import (
	"github.com/raincious/trap/trap/core/types"

	"time"
)

type MarkType byte
type UnmarkType byte

const (
	CLIENT_MARK_MANUAL MarkType = iota
	CLIENT_MARK_PICK
	CLIENT_MARK_OTHER
)

const (
	CLIENT_UNMARK_MANUAL UnmarkType = iota
	CLIENT_UNMARK_EXPIRE
)

type Clients struct {
	clients  map[types.IP]*Client
	onMark   func(*Client, MarkType)
	onUnmark func(*Client, UnmarkType)
	onRecord func(*Client, Record)
}

func NewClients(config Config) *Clients {
	return &Clients{
		clients:  map[types.IP]*Client{},
		onMark:   config.OnMark,
		onUnmark: config.OnUnmark,
		onRecord: config.OnRecord,
	}
}

func (c *Clients) Get(ip types.IP) (*Client, bool) {
	isNew := false

	if _, ok := c.clients[ip]; !ok {
		c.clients[ip] = &Client{
			address:        ip.IP(),
			firstSeen:      time.Now(),
			lastSeen:       time.Now(),
			count:          0,
			records:        []Record{},
			lastRecord:     nil,
			marked:         false,
			onMark:         c.onMark,
			onUnmark:       c.onUnmark,
			onRecord:       c.onRecord,
			tolerateCount:  0,
			tolerateExpire: time.Duration(0),
			restrictExpire: time.Duration(0),
		}

		isNew = true
	}

	return c.clients[ip], isNew
}

func (c *Clients) Has(ip types.IP) bool {
	if _, ok := c.clients[ip]; !ok {
		return false
	}

	return true
}

func (c *Clients) Len() int {
	return len(c.clients)
}

func (c *Clients) Delete(ip types.IP, ty UnmarkType) *types.Throw {
	if !c.Has(ip) {
		return ErrClientNotFound.Throw(ip)
	}

	if c.clients[ip].marked {
		c.clients[ip].Unmark(ty)
	}

	delete(c.clients, ip)

	return nil
}

func (c *Clients) Scan(
	callback func(types.IP, *Client) *types.Throw) *types.Throw {
	var err *types.Throw = nil

	for key, val := range c.clients {
		callbackErr := callback(key, val)

		if callbackErr == nil {
			continue
		}

		err = callbackErr

		break
	}

	return err
}

func (c *Clients) Clear() *types.Throw {
	var err *types.Throw = nil

	for key, _ := range c.clients {
		deleteErr := c.Delete(key, CLIENT_UNMARK_EXPIRE)

		if deleteErr != nil {
			err = deleteErr
		}
	}

	return err
}

func (c *Clients) Export() []ClientExport {
	clients := []ClientExport{}

	for _, clientInfo := range c.clients {
		clients = append(clients, ClientExport{
			Address:   clientInfo.Address(),
			FirstSeen: clientInfo.FirstSeen(),
			LastSeen:  clientInfo.LastSeen(),
			Count:     clientInfo.Count(),
			Records:   clientInfo.Records(),
			Marked:    clientInfo.Marked(),
		})
	}

	return clients
}

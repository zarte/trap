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

package sync

import (
	"net"

	"github.com/raincious/trap/trap/core/sync/communication/conn"
	"github.com/raincious/trap/trap/core/types"
)

type ActiveClientsTable struct {
	clients map[string]bool
	lock    types.Mutex
}

func NewActiveClientsTable() *ActiveClientsTable {
	return &ActiveClientsTable{
		lock:    types.Mutex{},
		clients: map[string]bool{},
	}
}

func (a *ActiveClientsTable) Has(conn *conn.Conn) bool {
	addr := conn.RemoteAddr().String()

	return a.has(addr)
}

func (a *ActiveClientsTable) HasAddr(addr net.Addr) bool {
	return a.has(addr.String())
}

func (a *ActiveClientsTable) has(key string) bool {
	var hasIt bool = false

	a.lock.Exec(func() {
		if _, ok := a.clients[key]; !ok {
			return
		}

		hasIt = true
	})

	return hasIt
}

func (a *ActiveClientsTable) Add(conn *conn.Conn) {
	addr := conn.RemoteAddr().String()

	a.lock.Exec(func() {
		a.clients[addr] = true
	})
}

func (a *ActiveClientsTable) Remove(conn *conn.Conn) {
	addr := conn.RemoteAddr().String()

	a.lock.Exec(func() {
		delete(a.clients, addr)
	})
}

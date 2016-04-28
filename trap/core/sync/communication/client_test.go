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
	"github.com/raincious/trap/trap/core/logger"
	"github.com/raincious/trap/trap/core/sync/communication/conn"
	"github.com/raincious/trap/trap/core/sync/communication/messager"
	"github.com/raincious/trap/trap/core/types"

	"net"
	"testing"
	"time"

	"log"
)

func TestClientConnect(t *testing.T) {
	responder := messager.Callbacks{}

	sessions := NewSessions(
		logger.NewLogger(),
		1024,
		10*time.Second,
		func() {}, func() {})

	client := NewClient(sessions, responder, 60*time.Second)

	addr := types.IPAddress{
		IP:   types.ConvertIP(net.ParseIP("127.0.0.1")),
		Port: 4430,
	}

	conErr := client.Connect(addr,
		func(conn *conn.Conn) {

		},
		func(conn *conn.Conn, err *types.Throw) {

		})

	if conErr != nil {
		panic(conErr)
	}

	defer func() {
		client.Disconnect()
	}()

	parentPath := types.IPAddresses{
		types.IPAddress{
			IP:   types.ConvertIP(net.ParseIP("127.0.0.1")),
			Port: types.UInt16(8080),
		},
		types.IPAddress{
			IP:   types.ConvertIP(net.ParseIP("127.0.0.2")),
			Port: types.UInt16(8081),
		},
		types.IPAddress{
			IP:   types.ConvertIP(net.ParseIP("127.0.0.2")),
			Port: types.UInt16(8081),
		},
	}

	serverPartners, authErr := client.Auth("PASSWORD", parentPath,
		func(conn *conn.Conn, ips types.IPAddresses) {

		})

	if authErr != nil {
		log.Printf("Auth failed: %s", authErr)

		return
	}

	log.Printf("Partners: %v", serverPartners)
}

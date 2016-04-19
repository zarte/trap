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
	"github.com/raincious/trap/trap/core/sync/communication"
	"github.com/raincious/trap/trap/core/sync/communication/conn"
	"github.com/raincious/trap/trap/core/sync/communication/messager"
	"github.com/raincious/trap/trap/core/types"

	"time"
)

var (
	ErrNodeNotConnected *types.Error = types.NewError(
		"The node '%s' is not been connected")

	ErrNodeRetryAfterTime *types.Error = types.NewError(
		"The node '%s' can only be connected after '%s'")
)

type Node struct {
	addr                  types.IPAddress
	password              types.String
	client                *communication.Client
	callbacks             messager.Callbacks
	requestTimeout        time.Duration
	connectionTimeout     time.Duration
	connectRetryPeriod    time.Duration
	maxConnectRetryPeriod time.Duration
	nextConnectAfter      time.Time
	connectionFailedCount uint64
	partners              NodeMap
	partnersLock          types.Mutex
}

func (n *Node) Client() *communication.Client {
	if n.client != nil {
		return n.client
	}

	n.client = communication.NewClient(
		n.callbacks,
		n.requestTimeout,
		n.connectionTimeout,
	)

	return n.client
}

func (n *Node) addConnectAfterWait() {
	retryPeriod := n.connectRetryPeriod * time.Duration(n.connectionFailedCount)

	n.connectionFailedCount++

	if retryPeriod > n.maxConnectRetryPeriod {
		retryPeriod = n.maxConnectRetryPeriod
	}

	n.nextConnectAfter = time.Now().Add(retryPeriod)
}

func (n *Node) Connect(
	connectedPartners types.IPAddresses,
	onConnected func(*conn.Conn),
	onDisconnected func(*conn.Conn, *types.Throw)) *types.Throw {
	currentTime := time.Now()

	if !currentTime.After(n.nextConnectAfter) {
		return ErrNodeRetryAfterTime.Throw(n.addr.String(), n.nextConnectAfter)
	}

	e := n.Client().Connect(n.addr, onConnected,
		func(conn *conn.Conn, err *types.Throw) {
			n.partnersLock.Exec(func() {
				n.partners = NodeMap{}
			})

			onDisconnected(conn, err)
		})

	if e != nil {
		n.addConnectAfterWait()

		return e
	}

	partners, authErr := n.Client().Auth(n.password, connectedPartners)

	if authErr != nil {
		n.addConnectAfterWait()

		return authErr
	}

	n.partnersLock.Exec(func() {
		for _, partner := range partners {
			n.partners[partner.String()] = partner
		}
	})

	n.connectionFailedCount = 0

	return nil
}

func (n *Node) Disconnect() *types.Throw {
	if n.client == nil {
		return ErrNodeNotConnected.Throw(n.addr.String())
	}

	return n.Client().Disconnect()
}

func (n *Node) IsConnected() bool {
	return n.Client().Connected()
}

func (n *Node) IsReconnectable() bool {
	if n.IsConnected() {
		return false
	}

	if !time.Now().After(n.nextConnectAfter) {
		return false
	}

	return true
}

func (n *Node) Address() types.IPAddress {
	return n.addr
}

func (n *Node) IsPartner(partner types.IPAddress) bool {
	isPartner := false

	n.partnersLock.Exec(func() {
		if _, ok := n.partners[partner.String()]; ok {
			isPartner = true
		}
	})

	return isPartner
}

func (n *Node) Partners() types.IPAddresses {
	partners := types.IPAddresses{}

	n.partnersLock.Exec(func() {
		for _, partner := range n.partners {
			partners = append(partners, partner)
		}
	})

	return partners
}

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
	"github.com/raincious/trap/trap/core/sync/communication/controller"
	"github.com/raincious/trap/trap/core/sync/communication/messager"
	"github.com/raincious/trap/trap/core/types"

	"time"
)

var (
	ErrNodeNotConnected *types.Error = types.NewError(
		"The node '%s' is not been connected")

	ErrNodeAlreadyConnected *types.Error = types.NewError(
		"The node '%s' is already been connected")

	ErrNodeIsMutexed *types.Error = types.NewError(
		"The node '%s' has been mutexed by another active node")

	ErrNodeRetryAfterTime *types.Error = types.NewError(
		"The node '%s' can only be connected after '%s'")
)

type Node struct {
	nodes                 *Nodes
	addr                  types.IPAddress
	addrStr               types.String
	password              types.String
	client                *communication.Client
	controller            *controller.Client
	requestTimeout        time.Duration
	connectionTimeout     time.Duration
	connectRetryPeriod    time.Duration
	maxConnectRetryPeriod time.Duration
	nextConnectAfter      time.Time
	connectionFailedCount uint64
	partners              NodeMap
	partnersLock          types.Mutex
	mutexWith             map[types.String]*Node
}

func (n *Node) Client() *communication.Client {
	if n.client != nil {
		return n.client
	}

	commonController := n.controller.Common

	commonController.Logger = n.controller.Common.Logger.NewContext(
		n.addr.String())

	clientController := controller.Client{
		Common: commonController,
		AddPartners: func(c *conn.Conn, ips types.IPAddresses) *types.Throw {
			n.partnersLock.Exec(func() {
				for _, partner := range ips {
					n.partners[partner.String()] = partner
				}
			})

			n.nodes.Scan(func(key types.String, node *Node) *types.Throw {
				for _, partner := range ips {
					if !node.addr.IsEqual(&partner) {
						continue
					}

					n.addMutexWith(node)
				}

				return nil
			})

			return n.controller.AddPartners(c, ips)
		},
		ConfilctedPartners: func(confilctedPartners types.IPAddresses) {
			n.nodes.Scan(func(key types.String, node *Node) *types.Throw {
				for _, partner := range confilctedPartners {
					if !node.IsPartner(partner) {
						continue
					}

					node.addMutexWith(n)
				}

				return nil
			})
		},
		RemovePartners: func(c *conn.Conn, ips types.IPAddresses) *types.Throw {
			n.partnersLock.Exec(func() {
				for _, partner := range ips {
					delete(n.partners, partner.String())
				}
			})

			n.nodes.Scan(func(key types.String, node *Node) *types.Throw {
				for _, partner := range ips {
					if !node.addr.IsEqual(&partner) {
						continue
					}

					n.removeMutexWith(node)
				}

				return nil
			})

			return n.controller.RemovePartners(c, ips)
		},
	}

	handle := messager.Callbacks{}

	handle.Register(messager.SYNC_SIGNAL_PARTNER_ADD,
		clientController.PartnersAdded)
	handle.Register(messager.SYNC_SIGNAL_PARTNER_REMOVE,
		clientController.PartnersRemoved)
	handle.Register(messager.SYNC_SIGNAL_CLIENT_MARK,
		clientController.ClientsMarked)
	handle.Register(messager.SYNC_SIGNAL_CLIENT_UNMARK,
		clientController.ClientsUnmarked)

	n.client = communication.NewClient(
		handle,
		n.requestTimeout,
		n.connectionTimeout,
	)

	return n.client
}

func (n *Node) addMutexWith(node *Node) {
	n.partnersLock.Exec(func() {
		n.mutexWith[node.addrStr] = node
	})
}

func (n *Node) removeMutexWith(node *Node) {
	n.partnersLock.Exec(func() {
		delete(n.mutexWith, node.addrStr)
	})
}

func (n *Node) isMutexWith(target *Node) bool {
	var hasIt bool = false

	if !n.IsConnected() {
		return false
	}

	n.partnersLock.Exec(func() {
		if _, ok := n.mutexWith[target.addrStr]; !ok {
			return
		}

		if n.mutexWith[target.addrStr] != target {
			return
		}

		hasIt = true
	})

	return hasIt
}

func (n *Node) addConnectAfterWait(nowTime time.Time) {
	retryPeriod := n.connectRetryPeriod * time.Duration(n.connectionFailedCount)

	n.connectionFailedCount++

	if retryPeriod > n.maxConnectRetryPeriod {
		retryPeriod = n.maxConnectRetryPeriod
	}

	n.nextConnectAfter = nowTime.Add(retryPeriod)
}

func (n *Node) Connect(
	connectedPartners types.IPAddresses,
	onConnected func(*conn.Conn),
	onAuthed func(*conn.Conn, types.IPAddresses),
	onDisconnected func(types.IPAddresses, *conn.Conn, *types.Throw),
) *types.Throw {
	currentTime := time.Now()

	if n.IsConnected() {
		return ErrNodeAlreadyConnected.Throw(n.addr.String())
	}

	if n.IsMutexed() {
		return ErrNodeIsMutexed.Throw(n.addr.String())
	}

	if !currentTime.After(n.nextConnectAfter) {
		return ErrNodeRetryAfterTime.Throw(n.addr.String(), n.nextConnectAfter)
	}

	e := n.Client().Connect(n.addr,
		func(conn *conn.Conn) {
			onConnected(conn)
		},
		func(conn *conn.Conn, err *types.Throw) {
			oldPartners := types.IPAddresses{}

			n.partnersLock.Exec(func() {
				for _, pVal := range n.partners {
					oldPartners = append(oldPartners, pVal)
				}

				n.partners = NodeMap{}
				n.mutexWith = map[types.String]*Node{}
			})

			onDisconnected(oldPartners, conn, err)
		})

	if e != nil {
		n.addConnectAfterWait(currentTime)

		return e
	}

	partners, authErr := n.Client().Auth(
		n.password,
		connectedPartners,
		func(conn *conn.Conn, ips types.IPAddresses) {
			newPartnerIPs := ips

			n.partnersLock.Exec(func() {
				for _, partner := range ips {
					n.partners[partner.String()] = partner
				}

				n.partners[n.addrStr] = n.addr
			})

			n.nodes.Scan(func(key types.String, node *Node) *types.Throw {
				for _, partner := range ips {
					if !node.addr.IsEqual(&partner) {
						continue
					}

					n.addMutexWith(node)
				}

				return nil
			})

			newPartnerIPs = append(newPartnerIPs, n.addr)

			onAuthed(conn, newPartnerIPs)
		})

	if authErr != nil {
		if authErr.Is(communication.ErrSessionAuthFailedConflicted) {
			n.nodes.Scan(func(key types.String, node *Node) *types.Throw {
				for _, partner := range partners {
					if !node.IsPartner(partner) {
						continue
					}

					node.addMutexWith(n)
				}

				return nil
			})
		}

		n.addConnectAfterWait(currentTime)

		return authErr
	}

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

	if n.IsMutexed() {
		return false
	}

	return true
}

func (n *Node) IsMutexed() bool {
	var isMutexed bool = false

	n.nodes.Scan(func(key types.String, node *Node) *types.Throw {
		if !node.isMutexWith(n) {
			return nil
		}

		isMutexed = true

		return ErrNodeScanBreakScan.Throw()
	})

	return isMutexed
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

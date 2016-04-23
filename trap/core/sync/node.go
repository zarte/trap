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
	nclient               *communication.Client
	controller            *controller.Client
	requestTimeout        time.Duration
	connectionTimeout     time.Duration
	connectRetryPeriod    time.Duration
	maxConnectRetryPeriod time.Duration
	nextConnectAfter      time.Time
	connectionFailedCount uint64
	partners              types.SearchableIPAddresses
	partnersLock          types.Mutex
	mutexWith             nodeMutexes
}

func (n *Node) client() *communication.Client {
	if n.nclient != nil {
		return n.nclient
	}

	commonController := n.controller.Common

	commonController.Logger = n.controller.Common.Logger.NewContext(
		n.addr.String())

	clientController := controller.Client{
		Common: commonController,
		AddPartners: func(
			c *conn.Conn, ips types.SearchableIPAddresses) *types.Throw {
			n.nodes.Scan(func(key types.String, node *Node) *types.Throw {
				if node == n {
					return nil
				}

				ips.Through(func(
					key types.IPAddressString,
					partner types.IPAddressWrapped,
				) *types.Throw {
					partnerAddr := partner.IPAddress()

					if node.addr.IsEqual(&partnerAddr) {
						n.addMutexWith(node)
					}

					if node.HasPartner(&partner) {
						n.addMutexWithPastners(node, &ips)
					}

					return nil
				})

				return nil
			})

			n.partnersLock.Exec(func() {
				ips.Through(func(
					key types.IPAddressString,
					partner types.IPAddressWrapped,
				) *types.Throw {
					n.partners.Insert(partner)

					return nil
				})
			})

			return n.controller.AddPartners(c, ips)
		},
		ConfilctedPartners: func(
			confilctedPartners types.SearchableIPAddresses) {
			n.nodes.Scan(func(key types.String, node *Node) *types.Throw {
				if node == n {
					return nil
				}

				confilctedPartners.Through(func(
					key types.IPAddressString,
					partner types.IPAddressWrapped,
				) *types.Throw {
					partnerAddr := partner.IPAddress()

					if node.addr.IsEqual(&partnerAddr) {
						node.addMutexWith(n)
					}

					if node.HasPartner(&partner) {
						node.addMutexWithPastners(n, &confilctedPartners)
					}

					return nil
				})

				return nil
			})
		},
		RemovePartners: func(
			c *conn.Conn, ips types.SearchableIPAddresses) *types.Throw {
			n.nodes.Scan(func(key types.String, node *Node) *types.Throw {
				if node == n {
					return nil
				}

				ips.Through(func(
					key types.IPAddressString,
					partner types.IPAddressWrapped,
				) *types.Throw {
					partnerAddr := partner.IPAddress()

					if !node.addr.IsEqual(&partnerAddr) &&
						!node.HasPartner(&partner) {
						return nil
					}

					n.removeMutexWith(node)

					return nil
				})

				return nil
			})

			n.partnersLock.Exec(func() {
				ips.Through(func(
					key types.IPAddressString,
					partner types.IPAddressWrapped,
				) *types.Throw {
					n.partners.Delete(&partner)

					return nil
				})
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

	n.nclient = communication.NewClient(
		handle,
		n.requestTimeout,
		n.connectionTimeout,
	)

	return n.nclient
}

func (n *Node) addMutexWith(node *Node) {
	n.partnersLock.Exec(func() {
		addrs := types.NewSearchableIPAddresses()

		addrs.Insert(node.addr.Wrapped())

		n.mutexWith.Append(node, addrs)
	})
}

func (n *Node) addMutexWithPastners(
	node *Node, confilcted *types.SearchableIPAddresses) {
	n.partnersLock.Exec(func() {
		n.mutexWith.Append(node, n.partners.Intersection(confilcted))
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

		if n.mutexWith[target.addrStr].With != target {
			return
		}

		due := n.mutexWith[target.addrStr].Due

		if !due.Contains(&n.partners) {
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
	connectedPartners types.SearchableIPAddresses,
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

	e := n.client().Connect(n.addr,
		func(conn *conn.Conn) {
			onConnected(conn)
		},
		func(conn *conn.Conn, err *types.Throw) {
			oldPartners := types.IPAddresses{}

			n.partnersLock.Exec(func() {
				oldPartners = n.partners.Export()

				n.partners = types.NewSearchableIPAddresses()
				n.mutexWith = nodeMutexes{}
			})

			onDisconnected(oldPartners, conn, err)
		})

	if e != nil {
		n.addConnectAfterWait(currentTime)

		return e
	}

	partners, authErr := n.client().Auth(
		n.password,
		connectedPartners.Export(),
		func(conn *conn.Conn, ips types.IPAddresses) {
			searchableNewPartners := ips.Searchable()

			n.nodes.Scan(func(key types.String, node *Node) *types.Throw {
				if node == n {
					return nil
				}

				searchableNewPartners.Through(func(
					key types.IPAddressString,
					partner types.IPAddressWrapped,
				) *types.Throw {
					partnerAddr := partner.IPAddress()

					if node.addr.IsEqual(&partnerAddr) {
						n.addMutexWith(node)
					}

					if node.HasPartner(&partner) {
						n.addMutexWithPastners(node, &searchableNewPartners)
					}

					if node.IsConnected() {
						node.Disconnect()
					}

					return nil
				})

				return nil
			})

			n.partnersLock.Exec(func() {
				n.partners.Insert(n.addr.Wrapped())

				searchableNewPartners.Through(func(
					key types.IPAddressString,
					partner types.IPAddressWrapped,
				) *types.Throw {
					n.partners.Insert(partner)

					return nil
				})
			})

			newPartnerIPs := ips

			newPartnerIPs = append(newPartnerIPs, n.addr)

			onAuthed(conn, newPartnerIPs)
		})

	if authErr != nil {
		if authErr.Is(communication.ErrSessionAuthFailedConflicted) {
			searchableNewPartners := partners.Searchable()

			n.nodes.Scan(func(key types.String, node *Node) *types.Throw {
				if node == n {
					return nil
				}

				searchableNewPartners.Through(func(
					key types.IPAddressString,
					p types.IPAddressWrapped,
				) *types.Throw {
					partnerAddr := p.IPAddress()

					if node.addr.IsEqual(&partnerAddr) {
						node.addMutexWith(n)
					}

					if node.HasPartner(&p) {
						node.addMutexWithPastners(n, &searchableNewPartners)
					}

					return nil
				})

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
	if n.nclient == nil {
		return ErrNodeNotConnected.Throw(n.addr.String())
	}

	return n.client().Disconnect()
}

func (n *Node) Delay() time.Duration {
	return n.client().Delay()
}

func (n *Node) IsConnected() bool {
	return n.client().Connected()
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

func (n *Node) HasPartner(partner *types.IPAddressWrapped) bool {
	isPartner := false

	n.partnersLock.Exec(func() {
		if !n.partners.Has(partner) {
			return
		}

		isPartner = true
	})

	return isPartner
}

func (n *Node) Partners() types.SearchableIPAddresses {
	partners := types.NewSearchableIPAddresses()

	n.partnersLock.Exec(func() {
		n.partners.Through(func(
			key types.IPAddressString,
			val types.IPAddressWrapped,
		) *types.Throw {
			partners.Insert(val)

			return nil
		})
	})

	return partners
}

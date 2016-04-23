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
	"github.com/raincious/trap/trap/core/sync/communication/controller"
	"github.com/raincious/trap/trap/core/types"

	"time"
)

var (
	ErrNodeAlreadyExisted *types.Error = types.NewError(
		"The node '%s' is already registered")

	ErrNodeNotExisted *types.Error = types.NewError(
		"The node '%s' is not existed")

	ErrNodeScanBreakScan *types.Error = types.NewError(
		"Break scaning")
)

const (
	MAX_NODE_CONNECT_RETRY = 32
)

type Nodes struct {
	nodeList          map[types.String]*Node
	nodeLock          types.Mutex
	controller        controller.Client
	requestTimeout    time.Duration
	connectionTimeout time.Duration
}

func NewNodes(defaultResponders controller.Client,
	requestTimeout time.Duration, connectionTimeout time.Duration) *Nodes {
	return &Nodes{
		controller:        defaultResponders,
		nodeList:          map[types.String]*Node{},
		nodeLock:          types.Mutex{},
		requestTimeout:    requestTimeout,
		connectionTimeout: connectionTimeout,
	}
}

func (n *Nodes) has(nodeKey types.String) bool {
	if _, ok := n.nodeList[nodeKey]; !ok {
		return false
	}

	return true
}

func (n *Nodes) Has(node types.IPAddress) bool {
	var hasIt bool = false

	n.nodeLock.Exec(func() {
		if !n.has(node.String()) {
			return
		}

		hasIt = true
	})

	return hasIt
}

func (n *Nodes) Remove(node types.IPAddress) *types.Throw {
	var err *types.Throw = nil

	nodeKey := node.String()

	n.nodeLock.Exec(func() {
		if n.has(nodeKey) {
			err = ErrNodeNotExisted.Throw(nodeKey)

			return
		}

		selectedNode := n.nodeList[nodeKey]

		disErr := selectedNode.Disconnect()

		if disErr != nil {
			if disErr.Is(ErrNodeNotConnected) {
				return
			}

			err = disErr

			return
		}

		delete(n.nodeList, nodeKey)
	})

	return err
}

func (n *Nodes) Register(ipAddr types.IPAddress,
	pass types.String) *types.Throw {
	var err *types.Throw = nil

	nodeKey := ipAddr.String()

	n.nodeLock.Exec(func() {
		if n.has(nodeKey) {
			err = ErrNodeAlreadyExisted.Throw(nodeKey)

			return
		}

		n.nodeList[nodeKey] = &Node{
			nodes:                 n,
			addr:                  ipAddr,
			addrStr:               ipAddr.String(),
			password:              pass,
			nclient:               nil,
			controller:            &n.controller,
			requestTimeout:        n.requestTimeout,
			connectionTimeout:     n.connectionTimeout,
			connectRetryPeriod:    n.connectionTimeout,
			maxConnectRetryPeriod: n.connectionTimeout * MAX_NODE_CONNECT_RETRY,
			nextConnectAfter:      time.Time{},
			connectionFailedCount: 0,
			partners:              types.NewSearchableIPAddresses(),
			partnersLock:          types.Mutex{},
			mutexWith:             nodeMutexes{},
		}
	})

	return err
}

func (n *Nodes) scan(
	scanner func(types.String, *Node) *types.Throw) *types.Throw {
	var scanResult *types.Throw = nil
	var err *types.Throw = nil
	var nodes map[types.String]*Node = map[types.String]*Node{}

	n.nodeLock.Exec(func() {
		for key, node := range n.nodeList {
			nodes[key] = node
		}
	})

	for key, node := range nodes {
		scanResult = scanner(key, node)

		if scanResult == nil {
			continue
		}

		if scanResult.Is(ErrNodeScanBreakScan) {
			break
		}

		err = scanResult

		break
	}

	return err
}

func (n *Nodes) Scan(
	scanner func(types.String, *Node) *types.Throw) *types.Throw {
	return n.scan(scanner)
}

func (n *Nodes) HasPartner(addr types.IPAddress) bool {
	var isPartner bool = false

	sAddr := addr.Wrapped()

	n.scan(func(key types.String, node *Node) *types.Throw {
		if !node.HasPartner(&sAddr) {
			return nil
		}

		isPartner = true

		return ErrNodeScanBreakScan.Throw()
	})

	return isPartner
}

func (n *Nodes) Partners() types.SearchableIPAddresses {
	partners := types.NewSearchableIPAddresses()

	n.scan(func(key types.String, node *Node) *types.Throw {
		p := node.Partners()

		p.Through(func(
			key types.IPAddressString,
			val types.IPAddressWrapped,
		) *types.Throw {
			partners.Insert(val)

			return nil
		})

		return nil
	})

	return partners
}

func (n *Nodes) Clear() *types.Throw {
	var err *types.Throw = nil

	n.scan(func(key types.String, node *Node) *types.Throw {
		err = n.Remove(node.Address())

		return nil
	})

	return err
}

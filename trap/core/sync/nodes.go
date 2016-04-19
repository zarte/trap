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
	"github.com/raincious/trap/trap/core/sync/communication/messager"
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
	MAX_NODE_CONNECT_RETRY_FCT = 32
)

type Nodes struct {
	nodeList map[types.String]*Node

	responders messager.Callbacks

	requestTimeout    time.Duration
	connectionTimeout time.Duration
}

func NewNodes(defaultResponders messager.Callbacks,
	requestTimeout time.Duration, connectionTimeout time.Duration) *Nodes {
	return &Nodes{
		responders:        defaultResponders,
		nodeList:          map[types.String]*Node{},
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
	if !n.has(node.String()) {
		return false
	}

	return true
}

func (n *Nodes) Remove(node types.IPAddress) *types.Throw {
	nodeKey := node.String()

	if n.has(nodeKey) {
		return ErrNodeNotExisted.Throw(nodeKey)
	}

	selectedNode := n.nodeList[nodeKey]

	disErr := selectedNode.Disconnect()

	if disErr != nil {
		if disErr.Is(ErrNodeNotConnected) {
			return nil
		}

		return disErr
	}

	delete(n.nodeList, nodeKey)

	return nil
}

func (n *Nodes) Register(ipAddr types.IPAddress,
	pass types.String) *types.Throw {
	nodeKey := ipAddr.String()

	if n.has(nodeKey) {
		return ErrNodeAlreadyExisted.Throw(nodeKey)
	}

	n.nodeList[nodeKey] = &Node{
		addr:                  ipAddr,
		password:              pass,
		client:                nil,
		callbacks:             n.responders,
		requestTimeout:        n.requestTimeout,
		connectionTimeout:     n.connectionTimeout,
		connectRetryPeriod:    n.connectionTimeout,
		maxConnectRetryPeriod: n.connectionTimeout * MAX_NODE_CONNECT_RETRY_FCT,
		nextConnectAfter:      time.Time{},
		connectionFailedCount: 0,
		partners:              NodeMap{},
		partnersLock:          types.Mutex{},
	}

	return nil
}

func (n *Nodes) scan(
	scanner func(types.String, *Node) *types.Throw) *types.Throw {
	var scanResult *types.Throw = nil
	var err *types.Throw = nil

	for key, node := range n.nodeList {
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

func (n *Nodes) IsPartner(addr types.IPAddress) bool {
	var isPartner bool = false

	n.scan(func(key types.String, node *Node) *types.Throw {
		if !node.IsPartner(addr) {
			return nil
		}

		isPartner = true

		return ErrNodeScanBreakScan.Throw()
	})

	return isPartner
}

func (n *Nodes) Partners() types.IPAddresses {
	partners := types.IPAddresses{}

	n.scan(func(key types.String, node *Node) *types.Throw {
		partners = append(partners, node.Partners()...)

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

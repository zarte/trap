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
	"github.com/raincious/trap/trap/core/types"
)

type nodeMutex struct {
	With *Node
	Due  types.SearchableIPAddresses
}

type nodeMutexes map[types.String]nodeMutex

func (m nodeMutexes) has(n *Node) bool {
	if _, ok := m[n.addrStr]; !ok {
		return false
	}

	return true
}

func (m nodeMutexes) Append(n *Node, ip types.SearchableIPAddresses) {
	if !m.has(n) {
		m[n.addrStr] = nodeMutex{
			With: n,
			Due:  types.NewSearchableIPAddresses(),
		}
	}

	mut := m[n.addrStr]

	ip.Through(func(
		key types.IPAddressString,
		val types.IPAddressWrapped,
	) *types.Throw {
		mut.Due.Insert(val)

		return nil
	})

	m[n.addrStr] = mut
}

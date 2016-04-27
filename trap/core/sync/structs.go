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

type nodeMap types.SearchableIPAddresses

type NodeInfo struct {
	Address   types.IPAddress
	Delay     time.Duration
	Stats     messager.Stats
	Connected bool
	Partner   types.IPAddresses
}

type ClientInfo struct {
	Remote types.IPAddress
	Stats  messager.Stats
}

type ServerInfo struct {
	Listen  types.IPAddress
	Clients []ClientInfo
}

type Status struct {
	Nodes  []NodeInfo
	Server ServerInfo
}

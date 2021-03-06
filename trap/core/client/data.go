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

	"net"
	"time"
)

type Hitting struct {
	types.IPAddress

	Type types.String
}

type Record struct {
	Inbound  []byte
	Outbound []byte
	Hitting  Hitting
	Time     time.Time
}

type ClientExport struct {
	Address   net.IP
	FirstSeen time.Time
	LastSeen  time.Time
	Count     types.UInt32
	Records   []Record
	Marked    bool
}

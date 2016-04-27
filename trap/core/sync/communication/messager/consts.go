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

package messager

const (
	SYNC_BUFFER_LENGTH = 256
)

const (
	SYNC_SIGNAL_UNDEFINED             = byte(0)
	SYNC_SIGNAL_BYE                   = byte(1)
	SYNC_SIGNAL_HEATBEAT              = byte(2)
	SYNC_SIGNAL_HEATBEAT_DENIED       = byte(3)
	SYNC_SIGNAL_HELLO                 = byte(4)
	SYNC_SIGNAL_HELLO_ACCEPT          = byte(5)
	SYNC_SIGNAL_HELLO_DENIED          = byte(6)
	SYNC_SIGNAL_HELLO_CONFLICT        = byte(7)
	SYNC_SIGNAL_PARTNER_ADD           = byte(8)
	SYNC_SIGNAL_PARTNER_ADD_ACCEPT    = byte(9)
	SYNC_SIGNAL_PARTNER_ADD_DENIED    = byte(10)
	SYNC_SIGNAL_PARTNER_REMOVE        = byte(11)
	SYNC_SIGNAL_PARTNER_REMOVE_ACCEPT = byte(12)
	SYNC_SIGNAL_PARTNER_REMOVE_DENIED = byte(13)
	SYNC_SIGNAL_CLIENT_MARK           = byte(14)
	SYNC_SIGNAL_CLIENT_MARK_ACCEPT    = byte(15)
	SYNC_SIGNAL_CLIENT_MARK_DENIED    = byte(16)
	SYNC_SIGNAL_CLIENT_UNMARK         = byte(17)
	SYNC_SIGNAL_CLIENT_UNMARK_ACCEPT  = byte(18)
	SYNC_SIGNAL_CLIENT_UNMARK_DENIED  = byte(19)
)

const (
	SYNC_CONTROLCHAR_ESCAPE      = byte(0)
	SYNC_CONTROLCHAR_TRANSMITTED = byte(1)
	SYNC_CONTROLCHAR_SEPARATOR   = byte(2)
	SYNC_CONTROLCHAR_RESERVED    = byte(255)
)

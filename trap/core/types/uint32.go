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

package types

import (
	"strconv"
)

type UInt32 uint32

func (i UInt32) String() String {
	return String(strconv.FormatUint(uint64(i), 10))
}

func (i UInt32) Int16() Int16 {
	if i > UINT32_MAX_INT16 {
		return Int16(MAX_INT16)
	}

	return Int16(i)
}

func (i UInt32) Int32() Int32 {
	if i > UINT32_MAX_INT32 {
		return Int32(MAX_INT32)
	}

	return Int32(i)
}

func (i UInt32) Int64() Int64 {
	return Int64(i)
}

func (i UInt32) UInt16() UInt16 {
	if i > UINT32_MAX_UINT16 {
		return UInt16(MAX_UINT16)
	}

	return UInt16(i)
}

func (i UInt32) UInt32() uint32 {
	return uint32(i)
}

func (i UInt32) UInt64() UInt64 {
	return UInt64(i)
}

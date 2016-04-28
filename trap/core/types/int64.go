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
	"encoding/binary"
	"strconv"
)

type Int64 int64

func (i Int64) String() String {
	return String(strconv.FormatInt(int64(i), 10))
}

func (i Int64) Int16() Int16 {
	if i > INT64_MAX_INT16 {
		return Int16(MAX_INT16)
	} else if i < INT64_MIN_INT16 {
		return Int16(MIN_INT16)
	}

	return Int16(i)
}

func (i Int64) Int32() Int32 {
	if i > INT64_MAX_INT32 {
		return Int32(MAX_INT32)
	} else if i < INT64_MIN_INT32 {
		return Int32(MIN_INT32)
	}

	return Int32(i)
}

func (i Int64) Int64() int64 {
	return int64(i)
}

func (i Int64) UInt16() UInt16 {
	if i < INT64_MIN_UINT16 {
		return UInt16(MIN_UINT16)
	} else if i > INT64_MAX_UINT16 {
		return UInt16(MAX_UINT16)
	}

	return UInt16(i)
}

func (i Int64) UInt32() UInt32 {
	if i < INT64_MIN_UINT32 {
		return UInt32(MIN_UINT32)
	} else if i > INT64_MAX_UINT32 {
		return UInt32(MAX_UINT32)
	}

	return UInt32(i)
}

func (i Int64) UInt64() UInt64 {
	if i < INT64_MIN_UINT64 {
		return UInt64(0)
	}

	return UInt64(i)
}

func (i Int64) Serialize() ([]byte, *Throw) {
	buf := make([]byte, 10)

	binary.PutVarint(buf, i.Int64())

	return buf, nil
}

func (i *Int64) Unserialize(text []byte) *Throw {
	if len(text) != 10 {
		return ErrTypesUnserializeInvalidDataLength.Throw(10)
	}

	num, n := binary.Varint(text)

	if n <= 0 {
		return ErrTypesUnserializeInvalidResult.Throw()
	}

	*i = Int64(num)

	return nil
}

type Int64Slice []Int64

func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }

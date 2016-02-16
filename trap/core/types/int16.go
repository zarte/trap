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

type Int16 int16

func (i Int16) String() (String) {
    return String(strconv.FormatInt(int64(i), 10))
}

func (i Int16) Int16() (int16) {
    return int16(i)
}

func (i Int16) Int32() (Int32) {
    return Int32(i)
}

func (i Int16) Int64() (Int64) {
    return Int64(i)
}

func (i Int16) UInt16() (UInt16) {
    if i < INT16_MIN_UINT16 {
        return UInt16(MIN_UINT16)
    }

    return UInt16(i)
}

func (i Int16) UInt32() (UInt32) {
    if i < INT16_MIN_UINT32 {
        return UInt32(MIN_UINT32)
    }

    return UInt32(i)
}

func (i Int16) UInt64() (UInt64) {
    if i < INT16_MIN_UINT64 {
        return UInt64(MIN_UINT64)
    }

    return UInt64(i)
}
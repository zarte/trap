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
    "encoding/binary"
)

type UInt16 uint16

func (i UInt16) String() (String) {
    return String(strconv.FormatUint(uint64(i), 10))
}

func (i UInt16) Int16() (Int16) {
    if i > UINT16_MAX_INT16 {
        return Int16(MAX_INT16)
    }

    return Int16(i)
}

func (i UInt16) Int32() (Int32) {
    return Int32(i)
}

func (i UInt16) Int64() (Int64) {
    return Int64(i)
}

func (i UInt16) UInt16() (uint16) {
    return uint16(i)
}

func (i UInt16) UInt32() (UInt32) {
    return UInt32(i)
}

func (i UInt16) UInt64() (UInt64) {
    return UInt64(i)
}

func (i UInt16) Serialize() ([]byte, *Throw) {
    uint16Byte :=  make([]byte, 2)

    binary.LittleEndian.PutUint16(uint16Byte, uint16(i))

    return uint16Byte, nil
}

func (i *UInt16) Unserialize(data []byte) (*Throw) {
    if len(data) != 2 {
        return ErrTypesUnserializeInvalidDataLength.Throw(2)
    }

    *i = UInt16(binary.LittleEndian.Uint16(data[:]))

    return nil
}
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

const (
    MAX_UINT            = ^uint(0)
    MAX_UINT8           = ^uint8(0)
    MAX_UINT16          = ^uint16(0)
    MAX_UINT32          = ^uint32(0)
    MAX_UINT64          = ^uint64(0)

    MIN_UINT            = uint(0)
    MIN_UINT8           = uint8(0)
    MIN_UINT16          = uint16(0)
    MIN_UINT32          = uint32(0)
    MIN_UINT64          = uint64(0)

    MAX_INT             = int(MAX_UINT >> 1)
    MAX_INT8            = int8(MAX_UINT8 >> 1)
    MAX_INT16           = int16(MAX_UINT16 >> 1)
    MAX_INT32           = int32(MAX_UINT32 >> 1)
    MAX_INT64           = int64(MAX_UINT64 >> 1)

    MIN_INT             = -MAX_INT - 1
    MIN_INT8            = -MAX_INT8 - 1
    MIN_INT16           = -MAX_INT16 - 1
    MIN_INT32           = -MAX_INT32 - 1
    MIN_INT64           = -MAX_INT64 - 1
)

const (
    INT16_MAX_INT16     = Int16(MAX_INT16)
    INT16_MIN_INT16     = Int16(MIN_INT16)
    INT32_MAX_INT16     = Int32(MAX_INT16)
    INT32_MIN_INT16     = Int32(MIN_INT16)
    INT64_MAX_INT16     = Int64(MAX_INT16)
    INT64_MIN_INT16     = Int64(MIN_INT16)

    INT32_MAX_INT32     = Int32(MAX_INT32)
    INT32_MIN_INT32     = Int32(MIN_INT32)
    INT64_MAX_INT32     = Int64(MAX_INT32)
    INT64_MIN_INT32     = Int64(MIN_INT32)

    INT64_MAX_INT64     = Int64(MAX_INT64)
    INT64_MIN_INT64     = Int64(MIN_INT64)

    INT32_MAX_UINT16    = Int32(MAX_UINT16)
    INT32_MIN_UINT16    = Int32(MIN_UINT16)
    INT64_MAX_UINT16    = Int64(MAX_UINT16)
    INT64_MIN_UINT16    = Int64(MIN_UINT16)

    INT64_MAX_UINT32    = Int64(MAX_UINT16)
    INT64_MIN_UINT32    = Int64(MIN_UINT16)


    UINT16_MAX_INT16    = UInt16(MAX_INT16)
    UINT32_MAX_INT16    = UInt32(MAX_INT16)
    UINT64_MAX_INT16    = UInt64(MAX_INT16)

    UINT32_MAX_INT32    = UInt32(MAX_INT32)
    UINT64_MAX_INT32    = UInt64(MAX_INT32)

    UINT64_MAX_INT64    = UInt64(MAX_INT64)


    INT16_MIN_UINT16    = Int16(MIN_UINT16)

    INT16_MIN_UINT32    = Int16(MIN_UINT32)
    INT32_MIN_UINT32    = Int32(MIN_UINT32)

    INT16_MIN_UINT64    = Int16(MIN_UINT64)
    INT32_MIN_UINT64    = Int32(MIN_UINT64)
    INT64_MIN_UINT64    = Int64(MIN_UINT64)


    UINT16_MAX_UINT16   = UInt16(MAX_UINT16)
    UINT16_MIN_UINT16   = UInt16(MIN_UINT16)

    UINT32_MAX_UINT32   = UInt32(MAX_UINT32)
    UINT32_MIN_UINT32   = UInt32(MIN_UINT32)

    UINT64_MAX_UINT64   = UInt64(MAX_UINT64)
    UINT64_MIN_UINT64   = UInt64(MIN_UINT64)
)
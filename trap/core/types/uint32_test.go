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
    "testing"
    "strconv"
)

func TestUInt32ToString(t *testing.T) {
    max             :=  UInt32(MAX_UINT32)
    min             :=  UInt32(MIN_UINT32)

    if max.String().String() != strconv.FormatUint(uint64(MAX_UINT32), 10) {
        t.Error("UInt32.String() failed convert number '%d' to string", max)
    }

    if min.String().String() != strconv.FormatUint(uint64(MIN_UINT32), 10) {
        t.Error("UInt32.String() failed convert number '%d' to string", min)
    }
}

func TestUInt32ToInt16(t *testing.T) {
    max             :=  UInt32(MAX_UINT32)
    min             :=  UInt32(MIN_UINT32)

    if max.Int16() != INT16_MAX_INT16 {
        t.Errorf("UInt32.Int16() failed convert number '%d'. " +
            "Excepting '%d', got '%d'", max, INT16_MAX_INT16, max.Int16())
    }

    if min.Int16() != Int16(0) {
        t.Errorf("UInt32.Int16() failed convert number '%d'. " +
            "Excepting '%d', got '%d'", min, Int16(0), min.Int16())
    }
}

func TestUInt32ToInt32(t *testing.T) {
    max             :=  UInt32(MAX_UINT32)
    min             :=  UInt32(MIN_UINT32)

    if max.Int32() != INT32_MAX_INT32 {
        t.Errorf("UInt32.Int32() failed convert number '%d'. " +
            "Excepting '%d', got '%d'", max, INT32_MAX_INT32, max.Int32())
    }

    if min.Int32() != Int32(0) {
        t.Errorf("UInt32.Int16() failed convert number '%d'. " +
            "Excepting '%d', got '%d'", min, Int32(0), min.Int32())
    }
}

func TestUInt32ToInt64(t *testing.T) {
    max             :=  UInt32(MAX_UINT32)
    min             :=  UInt32(MIN_UINT32)

    if max.Int64() != INT64_MAX_UINT32 {
        t.Errorf("UInt32.Int64() failed convert number '%d'. " +
            "Excepting '%d', got '%d'", max, INT64_MAX_UINT32, max.Int64())
    }

    if min.Int64() != Int64(0) {
        t.Errorf("UInt32.Int16() failed convert number '%d'. " +
            "Excepting '%d', got '%d'", min, Int64(0), min.Int64())
    }
}

func TestUInt32ToUInt16(t *testing.T) {
    max             :=  UInt32(MAX_UINT32)
    min             :=  UInt32(MIN_UINT32)

    if max.UInt16() != UINT16_MAX_UINT16 {
        t.Errorf("UInt32.UInt16() failed convert number '%d'. " +
            "Excepting '%d', got '%d'", max, UINT16_MAX_UINT16, max.UInt16())
    }

    if min.UInt16() != UInt16(0) {
        t.Errorf("UInt32.Int16() failed convert number '%d'. " +
            "Excepting '%d', got '%d'", min, UInt16(0), min.UInt16())
    }
}

func TestUInt32ToRealUInt32(t *testing.T) {
    max             :=  UInt32(MAX_UINT32)
    min             :=  UInt32(MIN_UINT32)

    if max.UInt32() != MAX_UINT32 {
        t.Errorf("UInt32.UInt32() failed convert number '%d'. " +
            "Excepting '%d', got '%d'", max, MAX_UINT32, max.UInt32())
    }

    if min.UInt32() != uint32(0) {
        t.Errorf("UInt32.Int16() failed convert number '%d'. " +
            "Excepting '%d', got '%d'", min, uint32(0), min.UInt32())
    }
}

func TestUInt32ToUInt64(t *testing.T) {
    max             :=  UInt32(MAX_UINT32)
    min             :=  UInt32(MIN_UINT32)

    if max.UInt64() != UINT64_MAX_UINT32 {
        t.Errorf("UInt32.UInt64() failed convert number '%d'. " +
            "Excepting '%d', got '%d'", max, UINT64_MAX_UINT32, max.UInt64())
    }

    if min.UInt64() != UInt64(0) {
        t.Errorf("UInt32.Int16() failed convert number '%d'. " +
            "Excepting '%d', got '%d'", min, UInt64(0), min.UInt64())
    }
}
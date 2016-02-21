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

func TestInt32ToString(t *testing.T) {
    // Try max
    maxInt32        :=  Int32(MAX_INT32)

    if maxInt32.String().String() != strconv.FormatInt(int64(MAX_INT32), 10) {
        t.Error("Int32.String() failed convert number '%d' to string", maxInt32)
    }

    // Try min
    minInt32        :=  Int32(MIN_INT32)

    if minInt32.String().String() != strconv.FormatInt(int64(MIN_INT32), 10) {
        t.Error("Int32.String() failed convert number '%d' to string", minInt32)
    }
}

func TestInt32ToInt16(t *testing.T) {
    // Try max
    maxInt32        :=  Int32(MAX_INT32)

    if maxInt32.Int16() != INT16_MAX_INT16 {
        t.Errorf("Int32.Int16() failed convert number '%d'. " +
            "Excepting '%d', got '%d'",
            maxInt32, INT16_MAX_INT16, maxInt32.Int16())
    }

    // Try min
    minInt32        :=  Int32(MIN_INT32)

    if minInt32.Int16() != INT16_MIN_INT16 {
        t.Errorf("Int32.Int16() failed convert number '%d'. " +
            "Excepting '%d', got '%d'",
            minInt32, INT16_MIN_INT16, minInt32.Int16())
    }
}

func TestInt32ToRealInt32(t *testing.T) {
    // Try max
    maxInt32        :=  Int32(MAX_INT32)

    if maxInt32.Int32() != int32(MAX_INT32) {
        t.Errorf("Int32.Int32() failed convert number '%d'. " +
            "Excepting '%d', got '%d'",
            maxInt32, int32(MAX_INT32), maxInt32.Int32())
    }

    // Try min
    minInt32        :=  Int32(MIN_INT32)

    if minInt32.Int32() != int32(MIN_INT32) {
        t.Errorf("Int32.Int32() failed convert number '%d'. " +
            "Excepting '%d', got '%d'",
            maxInt32, int32(MIN_INT32), maxInt32.Int32())
    }
}

func TestInt32ToInt64(t *testing.T) {
    // Try max
    maxInt32        :=  Int32(MAX_INT32)

    if maxInt32.Int64() != INT64_MAX_INT32 {
        t.Errorf("Int32.Int64() failed convert number '%d'. " +
            "Excepting '%d', got '%d'",
            maxInt32, INT64_MAX_INT32, maxInt32.Int64())
    }

    // Try min
    minInt32        :=  Int32(MIN_INT32)

    if minInt32.Int64() != INT64_MIN_INT32 {
        t.Errorf("Int32.Int64() failed convert number '%d'. " +
            "Excepting '%d', got '%d'",
            minInt32, INT64_MIN_INT32, minInt32.Int64())
    }
}

func TestInt32ToUInt16(t *testing.T) {
    // Try max
    maxInt32        :=  Int32(MAX_INT32)

    if maxInt32.UInt16() != UINT16_MAX_UINT16 {
        t.Errorf("Int32.UInt16() failed convert number '%d'. " +
            "Excepting '%d', got '%d'",
            maxInt32, UINT16_MAX_UINT16, maxInt32.UInt16())
    }

    // Try min
    minInt32        :=  Int32(MIN_INT32)

    if minInt32.UInt16() != UInt16(0) {
        t.Errorf("Int32.UInt16() failed convert number '%d'. " +
            "Excepting '%d', got '%d'",
            maxInt32, UInt16(0), maxInt32.UInt16())
    }
}

func TestInt32ToUInt32(t *testing.T) {
    // Try max
    maxInt32        :=  Int32(MAX_INT32)

    if maxInt32.UInt32() != UINT32_MAX_INT32 {
        t.Errorf("Int32.UInt32() failed convert number '%d'. " +
            "Excepting '%d', got '%d'",
            maxInt32, UINT32_MAX_INT32, maxInt32.UInt32())
    }

    // Try min
    minInt32        :=  Int32(MIN_INT32)

    if minInt32.UInt32() != UInt32(0) {
        t.Errorf("Int32.UInt32() failed convert number '%d'. " +
            "Excepting '%d', got '%d'",
            maxInt32, UInt32(0), maxInt32.UInt32())
    }
}

func TestInt32ToUInt64(t *testing.T) {
    // Try max
    maxInt32        :=  Int32(MAX_INT32)

    if maxInt32.UInt64() != UINT64_MAX_INT32 {
        t.Errorf("Int32.UInt64() failed convert number '%d'. " +
            "Excepting '%d', got '%d'",
            maxInt32, UINT64_MAX_INT32, maxInt32.UInt64())
    }

    // Try min
    minInt32        :=  Int32(MIN_INT32)

    if minInt32.UInt64() != UInt64(0) {
        t.Errorf("Int32.UInt64() failed convert number '%d'. " +
            "Excepting '%d', got '%d'",
            maxInt32, UInt64(0), maxInt32.UInt64())
    }
}
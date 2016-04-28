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
	"testing"
)

func TestInt64ToString(t *testing.T) {
	// Try max
	maxInt64 := Int64(MAX_INT64)

	if maxInt64.String().String() != strconv.FormatInt(int64(MAX_INT64), 10) {
		t.Errorf("Int64.String() failed convert number '%d' to string", maxInt64)
	}

	// Try min
	minInt64 := Int64(MAX_INT64)

	if minInt64.String().String() != strconv.FormatInt(int64(MAX_INT64), 10) {
		t.Errorf("Int64.String() failed convert number '%d' to string", minInt64)
	}
}

func TestInt64ToInt16(t *testing.T) {
	// Try max
	maxInt64 := Int64(MAX_INT64)

	if maxInt64.Int16() != INT16_MAX_INT16 {
		t.Errorf("Int64.Int16() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			maxInt64, INT16_MAX_INT16, maxInt64.Int16())
	}

	// Try min
	minInt64 := Int64(MIN_INT64)

	if minInt64.Int16() != INT16_MIN_INT16 {
		t.Errorf("Int64.Int16() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			minInt64, INT16_MIN_INT16, minInt64.Int16())
	}
}

func TestInt64ToInt32(t *testing.T) {
	// Try max
	maxInt64 := Int64(MAX_INT64)

	if maxInt64.Int32() != INT32_MAX_INT32 {
		t.Errorf("Int64.Int32() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			maxInt64, INT32_MAX_INT32, maxInt64.Int32())
	}

	// Try min
	minInt64 := Int64(MIN_INT64)

	if minInt64.Int32() != INT32_MIN_INT32 {
		t.Errorf("Int64.Int32() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			minInt64, INT32_MIN_INT32, minInt64.Int32())
	}
}

func TestInt64ToRealInt64(t *testing.T) {
	// Try max
	maxInt64 := Int64(MAX_INT64)

	if maxInt64.Int64() != MAX_INT64 {
		t.Errorf("Int64.Int64() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			maxInt64, MAX_INT64, maxInt64.Int64())
	}

	// Try min
	minInt64 := Int64(MIN_INT64)

	if minInt64.Int64() != MIN_INT64 {
		t.Errorf("Int64.Int64() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			minInt64, MIN_INT64, minInt64.Int64())
	}
}

func TestInt64ToUInt16(t *testing.T) {
	// Try max
	maxInt64 := Int64(MAX_INT64)

	if maxInt64.UInt16() != UINT16_MAX_UINT16 {
		t.Errorf("Int64.UInt16() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			maxInt64, UINT16_MAX_UINT16, maxInt64.UInt16())
	}

	// Try min
	minInt64 := Int64(MIN_INT64)

	if minInt64.UInt16() != UINT16_MIN_UINT16 {
		t.Errorf("Int64.UInt16() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			minInt64, UINT16_MIN_UINT16, minInt64.UInt16())
	}
}

func TestInt64ToUInt32(t *testing.T) {
	// Try max
	maxInt64 := Int64(MAX_INT64)

	if maxInt64.UInt32() != UINT32_MAX_UINT32 {
		t.Errorf("Int64.UInt32() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			maxInt64, UINT32_MAX_UINT32, maxInt64.UInt32())
	}

	// Try min
	minInt64 := Int64(MIN_INT64)

	if minInt64.UInt32() != UINT32_MIN_UINT32 {
		t.Errorf("Int64.UInt32() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			minInt64, UINT32_MIN_UINT32, minInt64.UInt32())
	}
}

func TestInt64ToUInt64(t *testing.T) {
	// Try max
	maxInt64 := Int64(MAX_INT64)

	if maxInt64.UInt64() != UINT64_MAX_INT64 {
		t.Errorf("Int64.UInt64() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			maxInt64, UINT32_MAX_UINT32, maxInt64.UInt64())
	}

	// Try min
	minInt64 := Int64(MIN_INT64)

	if minInt64.UInt64() != UInt64(0) {
		t.Errorf("Int64.UInt64() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			minInt64, UInt64(0), minInt64.UInt64())
	}
}

func TestInt64SerializeUnserialize(t *testing.T) {
	maxInt64 := Int64(MAX_INT64)

	n, _ := maxInt64.Serialize()

	num := Int64(0)

	num.Unserialize(n)

	if maxInt64 != num {
		t.Errorf("Int64.Serialize() or  Int64.Serialize() is failed. "+
			"Excepting result to be '%d', got '%d'", maxInt64, num)
	}
}

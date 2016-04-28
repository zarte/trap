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

func TestUInt64ToString(t *testing.T) {
	// Try max
	maxUInt64 := UInt64(MAX_UINT64)

	if maxUInt64.String().String() != strconv.FormatUint(uint64(MAX_UINT64), 10) {
		t.Errorf("UInt64.String() failed convert number '%d' to string", maxUInt64.UInt64())
	}

	// Try min
	minUInt64 := UInt64(MIN_UINT64)

	if minUInt64.String().String() != strconv.FormatUint(uint64(MIN_UINT64), 10) {
		t.Errorf("UInt64.String() failed convert number '%d' to string", minUInt64.UInt64())
	}
}

func TestUInt64ToInt16(t *testing.T) {
	// Try max
	maxUInt64 := UInt64(MAX_UINT64)

	if maxUInt64.Int16() != INT16_MAX_INT16 {
		t.Errorf("UInt64.Int16() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			maxUInt64, INT16_MAX_INT16, maxUInt64.Int16())
	}

	// Try min
	minUInt64 := UInt64(MIN_UINT64)

	if minUInt64.Int16() != Int16(0) {
		t.Errorf("UInt64.Int16() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			minUInt64, Int16(0), minUInt64.Int16())
	}
}

func TestUInt64ToInt32(t *testing.T) {
	// Try max
	maxUInt64 := UInt64(MAX_UINT64)

	if maxUInt64.Int32() != INT32_MAX_INT32 {
		t.Errorf("UInt64.Int32() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			maxUInt64, INT32_MAX_INT32, maxUInt64.Int32())
	}

	// Try min
	minUInt64 := UInt64(MIN_UINT64)

	if minUInt64.Int32() != Int32(0) {
		t.Errorf("UInt64.Int32() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			minUInt64, Int32(0), minUInt64.Int32())
	}
}

func TestUInt64ToInt64(t *testing.T) {
	// Try max
	maxUInt64 := UInt64(MAX_UINT64)

	if maxUInt64.Int64() != INT64_MAX_INT64 {
		t.Errorf("UInt64.Int64() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			maxUInt64, INT64_MAX_INT64, maxUInt64.Int64())
	}

	// Try min
	minUInt64 := UInt64(MIN_UINT64)

	if minUInt64.Int64() != Int64(0) {
		t.Errorf("UInt64.Int64() failed convert number '%d'. "+
			"Excepting '%d', got '%d'",
			minUInt64, Int64(0), minUInt64.Int64())
	}
}

func TestUInt64ToUInt16(t *testing.T) {
	max := UInt64(MAX_UINT64)
	min := UInt64(MIN_UINT64)

	if max.UInt16() != UINT16_MAX_UINT16 {
		t.Errorf("UInt64.UInt16() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", max, UINT16_MAX_UINT16, max.UInt16())
	}

	if min.UInt16() != UInt16(0) {
		t.Errorf("UInt64.UInt16() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", min, UInt16(0), min.UInt16())
	}
}

func TestUInt64ToUInt32(t *testing.T) {
	max := UInt64(MAX_UINT64)
	min := UInt64(MIN_UINT64)

	if max.UInt32() != UINT32_MAX_UINT32 {
		t.Errorf("UInt64.UInt32() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", max, UINT32_MAX_UINT32, max.UInt32())
	}

	if min.UInt32() != UInt32(0) {
		t.Errorf("UInt64.UInt32() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", min, UInt32(0), min.UInt32())
	}
}

func TestUInt64ToRealUInt64(t *testing.T) {
	max := UInt64(MAX_UINT64)
	min := UInt64(MIN_UINT64)

	if max.UInt64() != MAX_UINT64 {
		t.Errorf("UInt64.UInt64() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", max, MAX_UINT64, max.UInt64())
	}

	if min.UInt64() != MIN_UINT64 {
		t.Errorf("UInt64.UInt64() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", min, MIN_UINT64, min.UInt64())
	}
}

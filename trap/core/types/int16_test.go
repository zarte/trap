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

func TestInt16ToString(t *testing.T) {
	// Try a randomly picked number
	testNumber := Int16(1024)

	if testNumber.String() != "1024" {
		t.Error("Int16.String() failed convert number to string")
	}

	// Try max
	maxInt16 := Int16(MAX_INT16)

	if maxInt16.String().String() != strconv.FormatInt(int64(MAX_INT16), 10) {
		t.Errorf("Int16.String() failed convert number '%d' to string", maxInt16.Int16())
	}

	// Try min
	minInt16 := Int16(MIN_INT16)

	if minInt16.String().String() != strconv.FormatInt(int64(MIN_INT16), 10) {
		t.Errorf("Int16.String() failed convert number '%d' to string", minInt16.Int16())
	}
}

func TestInt16ToRealInt16(t *testing.T) {
	// Try a randomly picked number
	testNumber := Int16(2414)

	if testNumber.Int16() != int16(2414) {
		t.Error("Int16.Int16() failed convert Int16 to real int16")
	}

	// Try max
	maxInt16 := Int16(MAX_INT16)

	if maxInt16.Int16() != int16(MAX_INT16) {
		t.Errorf("Int16.Int16() failed convert number '%d' to real int16",
			maxInt16)
	}

	// Try min
	minInt16 := Int16(MIN_INT16)

	if minInt16.Int16() != int16(MIN_INT16) {
		t.Errorf("Int16.Int16() failed convert number '%d' to real int16",
			minInt16)
	}
}

func TestInt16ToInt32(t *testing.T) {
	// Try a randomly picked number
	testNumber := Int16(23333)

	if testNumber.Int32() != Int32(23333) {
		t.Error("Int16.Int32() failed convert Int16 to Int32")
	}

	// Try max
	maxInt16 := Int16(MAX_INT16)

	if maxInt16.Int32() != INT32_MAX_INT16 {
		t.Errorf("Int16.Int32() failed convert number '%d' to Int32. "+
			"Excepting '%d', got '%d'",
			maxInt16, INT32_MAX_INT16, maxInt16.Int32())
	}

	// Try min
	minInt16 := Int16(MIN_INT16)

	if minInt16.Int32() != INT32_MIN_INT16 {
		t.Errorf("Int16.Int32() failed convert number '%d' to Int32. "+
			"Excepting '%d', got '%d'",
			minInt16, INT32_MIN_INT16, minInt16.Int32())
	}
}

func TestInt16ToInt64(t *testing.T) {
	// Try a randomly picked number
	testNumber := Int16(32451)

	if testNumber.Int64() != Int64(32451) {
		t.Error("Int16.Int64() failed convert Int16 to Int64")
	}

	// Try max
	maxInt16 := Int16(MAX_INT16)

	if maxInt16.Int64() != INT64_MAX_INT16 {
		t.Errorf("Int16.Int64() failed convert number '%d' to Int64. "+
			"Excepting '%d', got '%d'",
			maxInt16, INT64_MAX_INT16, maxInt16.Int64())
	}

	// Try min
	minInt16 := Int16(MIN_INT16)

	if minInt16.Int64() != INT64_MIN_INT16 {
		t.Errorf("Int16.Int64() failed convert number '%d' to Int64. "+
			"Excepting '%d', got '%d'",
			minInt16, INT64_MIN_INT16, minInt16.Int64())
	}
}

func TestInt16ToUInt16(t *testing.T) {
	// Try max
	maxInt16 := Int16(MAX_INT16)

	if maxInt16.UInt16() != UINT16_MAX_INT16 {
		t.Errorf("Int16.UInt16() failed convert number '%d' to UInt16. "+
			"Excepting '%d', got '%d'",
			maxInt16, UINT16_MAX_INT16, maxInt16.UInt16())
	}

	// Try min
	minInt16 := Int16(MIN_INT16)

	if minInt16.UInt16() != UInt16(0) {
		t.Errorf("Int16.UInt16() failed convert number '%d' to UInt16. "+
			"Excepting '%d', got '%d'",
			minInt16, UInt16(0), minInt16.UInt16())
	}
}

func TestInt16ToUInt32(t *testing.T) {
	// Try max
	maxInt16 := Int16(MAX_INT16)

	if maxInt16.UInt32() != UINT32_MAX_INT16 {
		t.Errorf("Int16.UInt32() failed convert number '%d' to UInt32. "+
			"Excepting '%d', got '%d'",
			maxInt16, UINT32_MAX_INT16, maxInt16.UInt32())
	}

	// Try min
	minInt16 := Int16(MIN_INT16)

	if minInt16.UInt32() != UInt32(0) {
		t.Errorf("Int16.UInt32() failed convert number '%d' to UInt32. "+
			"Excepting '%d', got '%d'",
			minInt16, UInt32(0), minInt16.UInt32())
	}
}

func TestInt16ToUInt64(t *testing.T) {
	// Try max
	maxInt16 := Int16(MAX_INT16)

	if maxInt16.UInt64() != UINT64_MAX_INT16 {
		t.Errorf("Int16.UInt64() failed convert number '%d' to UInt64. "+
			"Excepting '%d', got '%d'",
			maxInt16, UINT64_MAX_INT16, maxInt16.UInt64())
	}

	// Try min
	minInt16 := Int16(MIN_INT16)

	if minInt16.UInt64() != UInt64(0) {
		t.Errorf("Int16.UInt64() failed convert number '%d' to UInt64. "+
			"Excepting '%d', got '%d'",
			minInt16, UInt64(0), minInt16.UInt64())
	}
}

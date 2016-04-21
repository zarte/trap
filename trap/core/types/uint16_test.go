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

func TestUInt16ToString(t *testing.T) {
	max := UInt16(MAX_UINT16)
	min := UInt16(MIN_UINT16)

	if max.String().String() != strconv.FormatUint(uint64(MAX_UINT16), 10) {
		t.Errorf("UInt16.String() failed convert number '%d' to string", max.UInt16())
	}

	if min.String().String() != strconv.FormatUint(uint64(MIN_UINT16), 10) {
		t.Errorf("UInt16.String() failed convert number '%d' to string", min.UInt16())
	}
}

func TestUInt16ToInt16(t *testing.T) {
	max := UInt16(MAX_UINT16)
	min := UInt16(MIN_UINT16)

	if max.Int16() != INT16_MAX_INT16 {
		t.Errorf("UInt16.Int16() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", max, INT16_MAX_INT16, max.Int16())
	}

	if min.Int16() != Int16(0) {
		t.Errorf("UInt16.Int16() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", min, Int16(0), min.Int16())
	}
}

func TestUInt16ToInt32(t *testing.T) {
	max := UInt16(MAX_UINT16)
	min := UInt16(MIN_UINT16)

	if max.Int32() != INT32_MAX_UINT16 {
		t.Errorf("UInt16.Int32() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", max, INT32_MAX_UINT16, max.Int32())
	}

	if min.Int32() != Int32(0) {
		t.Errorf("UInt16.Int32() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", min, Int32(0), min.Int32())
	}
}

func TestUInt16ToInt64(t *testing.T) {
	max := UInt16(MAX_UINT16)
	min := UInt16(MIN_UINT16)

	if max.Int64() != INT64_MAX_UINT16 {
		t.Errorf("UInt16.Int64() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", max, INT64_MAX_UINT16, max.Int64())
	}

	if min.Int64() != Int64(0) {
		t.Errorf("UInt16.Int64() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", min, Int64(0), min.Int64())
	}
}

func TestUInt16ToRealUInt16(t *testing.T) {
	max := UInt16(MAX_UINT16)
	min := UInt16(MIN_UINT16)

	if max.UInt16() != MAX_UINT16 {
		t.Errorf("UInt16.UInt16() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", max, MAX_UINT16, max.UInt16())
	}

	if min.UInt16() != uint16(0) {
		t.Errorf("UInt16.UInt16() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", min, uint16(0), min.UInt16())
	}
}

func TestUInt16ToUInt32(t *testing.T) {
	max := UInt16(MAX_UINT16)
	min := UInt16(MIN_UINT16)

	if max.UInt32() != UINT32_MAX_UINT16 {
		t.Errorf("UInt16.UInt32() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", max, UINT32_MAX_UINT16, max.UInt32())
	}

	if min.UInt32() != UInt32(0) {
		t.Errorf("UInt16.UInt32() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", min, UInt32(0), min.UInt32())
	}
}

func TestUInt16ToUInt64(t *testing.T) {
	max := UInt16(MAX_UINT16)
	min := UInt16(MIN_UINT16)

	if max.UInt64() != UINT64_MAX_UINT16 {
		t.Errorf("UInt16.UInt64() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", max, UINT64_MAX_UINT16, max.UInt64())
	}

	if min.UInt64() != UInt64(0) {
		t.Errorf("UInt16.UInt64() failed convert number '%d'. "+
			"Excepting '%d', got '%d'", min, UInt64(0), min.UInt64())
	}
}

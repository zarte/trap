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

package event

import (
	"github.com/raincious/trap/trap/core/types"

	"testing"
)

func TestParameterSetGetInt16(t *testing.T) {
	param1 := Parameter{}
	param2 := Parameter{}

	param1.setInt16(types.INT16_MAX_INT16)
	param2.setInt16(types.INT16_MIN_INT16)

	if types.INT16_MAX_INT16 != param1.GetInt16() {
		t.Error("Parameter.GetInt16() didn't outputs the original value")

		return
	}

	if types.INT16_MIN_INT16 != param2.GetInt16() {
		t.Error("Parameter.GetInt16() didn't outputs the original value")

		return
	}

	if types.INT16_MAX_INT16.String() != param1.String() {
		t.Error("Parameter.GetInt16() didn't outputs the original value")

		return
	}

	if types.INT16_MIN_INT16.String() != param2.String() {
		t.Error("Parameter.GetInt16() didn't outputs the original value")

		return
	}

	if param1.GetInt32() != 0 ||
		param1.GetInt64() != 0 ||
		param1.GetUInt16() != 0 ||
		param1.GetUInt32() != 0 ||
		param1.GetUInt64() != 0 ||
		param1.GetStr() != "" ||
		string(param1.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}

	if param2.GetInt32() != 0 ||
		param2.GetInt64() != 0 ||
		param2.GetUInt16() != 0 ||
		param2.GetUInt32() != 0 ||
		param2.GetUInt64() != 0 ||
		param2.GetStr() != "" ||
		string(param2.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}
}

func TestParameterSetGetInt32(t *testing.T) {
	param1 := Parameter{}
	param2 := Parameter{}

	param1.setInt32(types.INT32_MAX_INT32)
	param2.setInt32(types.INT32_MIN_INT32)

	if types.INT32_MAX_INT32 != param1.GetInt32() {
		t.Error("Parameter.GetInt32() didn't outputs the original value")

		return
	}

	if types.INT32_MIN_INT32 != param2.GetInt32() {
		t.Error("Parameter.GetInt32() didn't outputs the original value")

		return
	}

	if types.INT32_MAX_INT32.String() != param1.String() {
		t.Error("Parameter.GetInt32() didn't outputs the original value")

		return
	}

	if types.INT32_MIN_INT32.String() != param2.String() {
		t.Error("Parameter.GetInt32() didn't outputs the original value")

		return
	}

	if param1.GetInt16() != 0 ||
		param1.GetInt64() != 0 ||
		param1.GetUInt16() != 0 ||
		param1.GetUInt32() != 0 ||
		param1.GetUInt64() != 0 ||
		param1.GetStr() != "" ||
		string(param1.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}

	if param2.GetInt16() != 0 ||
		param2.GetInt64() != 0 ||
		param2.GetUInt16() != 0 ||
		param2.GetUInt32() != 0 ||
		param2.GetUInt64() != 0 ||
		param2.GetStr() != "" ||
		string(param2.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}
}

func TestParameterSetGetInt64(t *testing.T) {
	param1 := Parameter{}
	param2 := Parameter{}

	param1.setInt64(types.INT64_MAX_INT64)
	param2.setInt64(types.INT64_MIN_INT64)

	if types.INT64_MAX_INT64 != param1.GetInt64() {
		t.Error("Parameter.GetInt64() didn't outputs the original value")

		return
	}

	if types.INT64_MIN_INT64 != param2.GetInt64() {
		t.Error("Parameter.GetInt64() didn't outputs the original value")

		return
	}

	if types.INT64_MAX_INT64.String() != param1.String() {
		t.Error("Parameter.GetInt64() didn't outputs the original value")

		return
	}

	if types.INT64_MIN_INT64.String() != param2.String() {
		t.Error("Parameter.GetInt64() didn't outputs the original value")

		return
	}

	if param1.GetInt16() != 0 ||
		param1.GetInt32() != 0 ||
		param1.GetUInt16() != 0 ||
		param1.GetUInt32() != 0 ||
		param1.GetUInt64() != 0 ||
		param1.GetStr() != "" ||
		string(param1.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}

	if param2.GetInt16() != 0 ||
		param2.GetInt32() != 0 ||
		param2.GetUInt16() != 0 ||
		param2.GetUInt32() != 0 ||
		param2.GetUInt64() != 0 ||
		param2.GetStr() != "" ||
		string(param2.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}
}

func TestParameterSetGetUInt16(t *testing.T) {
	param1 := Parameter{}
	param2 := Parameter{}

	param1.setUInt16(types.UINT16_MAX_UINT16)
	param2.setUInt16(types.UINT16_MIN_UINT16)

	if types.UINT16_MAX_UINT16 != param1.GetUInt16() {
		t.Error("Parameter.GetUInt16() didn't outputs the original value")

		return
	}

	if types.UINT16_MIN_UINT16 != param2.GetUInt16() {
		t.Error("Parameter.GetUInt16() didn't outputs the original value")

		return
	}

	if types.UINT16_MAX_UINT16.String() != param1.String() {
		t.Error("Parameter.GetUInt16() didn't outputs the original value")

		return
	}

	if types.UINT16_MIN_UINT16.String() != param2.String() {
		t.Error("Parameter.GetUInt16() didn't outputs the original value")

		return
	}

	if param1.GetInt16() != 0 ||
		param1.GetInt32() != 0 ||
		param1.GetInt64() != 0 ||
		param1.GetUInt32() != 0 ||
		param1.GetUInt64() != 0 ||
		param1.GetStr() != "" ||
		string(param1.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}

	if param2.GetInt16() != 0 ||
		param2.GetInt32() != 0 ||
		param2.GetInt64() != 0 ||
		param2.GetUInt32() != 0 ||
		param2.GetUInt64() != 0 ||
		param2.GetStr() != "" ||
		string(param2.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}
}

func TestParameterSetGetUInt32(t *testing.T) {
	param1 := Parameter{}
	param2 := Parameter{}

	param1.setUInt32(types.UINT32_MAX_UINT32)
	param2.setUInt32(types.UINT32_MIN_UINT32)

	if types.UINT32_MAX_UINT32 != param1.GetUInt32() {
		t.Error("Parameter.GetUInt32() didn't outputs the original value")

		return
	}

	if types.UINT32_MIN_UINT32 != param2.GetUInt32() {
		t.Error("Parameter.GetUInt32() didn't outputs the original value")

		return
	}

	if types.UINT32_MAX_UINT32.String() != param1.String() {
		t.Error("Parameter.GetUInt32() didn't outputs the original value")

		return
	}

	if types.UINT32_MIN_UINT32.String() != param2.String() {
		t.Error("Parameter.GetUInt32() didn't outputs the original value")

		return
	}

	if param1.GetInt16() != 0 ||
		param1.GetInt32() != 0 ||
		param1.GetInt64() != 0 ||
		param1.GetUInt16() != 0 ||
		param1.GetUInt64() != 0 ||
		param1.GetStr() != "" ||
		string(param1.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}

	if param2.GetInt16() != 0 ||
		param2.GetInt32() != 0 ||
		param2.GetInt64() != 0 ||
		param2.GetUInt16() != 0 ||
		param2.GetUInt64() != 0 ||
		param2.GetStr() != "" ||
		string(param2.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}
}

func TestParameterSetGetUInt64(t *testing.T) {
	param1 := Parameter{}
	param2 := Parameter{}

	param1.setUInt64(types.UINT64_MAX_UINT64)
	param2.setUInt64(types.UINT64_MIN_UINT64)

	if types.UINT64_MAX_UINT64 != param1.GetUInt64() {
		t.Error("Parameter.GetUInt64() didn't outputs the original value")

		return
	}

	if types.UINT64_MIN_UINT64 != param2.GetUInt64() {
		t.Error("Parameter.GetUInt64() didn't outputs the original value")

		return
	}

	if types.UINT64_MAX_UINT64.String() != param1.String() {
		t.Error("Parameter.GetUInt64() didn't outputs the original value")

		return
	}

	if types.UINT64_MIN_UINT64.String() != param2.String() {
		t.Error("Parameter.GetUInt64() didn't outputs the original value")

		return
	}

	if param1.GetInt16() != 0 ||
		param1.GetInt32() != 0 ||
		param1.GetInt64() != 0 ||
		param1.GetUInt16() != 0 ||
		param1.GetUInt32() != 0 ||
		param1.GetStr() != "" ||
		string(param1.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}

	if param2.GetInt16() != 0 ||
		param2.GetInt32() != 0 ||
		param2.GetInt64() != 0 ||
		param2.GetUInt16() != 0 ||
		param2.GetUInt32() != 0 ||
		param2.GetStr() != "" ||
		string(param2.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}
}

func TestParameterSetGetStr(t *testing.T) {
	param1 := Parameter{}
	param1.setStr("This is a simple string")

	if "This is a simple string" != param1.GetStr() {
		t.Error("Parameter.GetStr() didn't outputs the original value")

		return
	}

	if param1.GetInt16() != 0 ||
		param1.GetInt32() != 0 ||
		param1.GetInt64() != 0 ||
		param1.GetUInt16() != 0 ||
		param1.GetUInt32() != 0 ||
		param1.GetUInt64() != 0 ||
		string(param1.GetBytes()) != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}
}

func TestParameterSetGetBytes(t *testing.T) {
	param1 := Parameter{}
	param1.setBytes([]byte("This is a simple string"))

	if "This is a simple string" != string(param1.GetBytes()) {
		t.Error("Parameter.GetBytes() didn't outputs the original value")

		return
	}

	if param1.GetInt16() != 0 ||
		param1.GetInt32() != 0 ||
		param1.GetInt64() != 0 ||
		param1.GetUInt16() != 0 ||
		param1.GetUInt32() != 0 ||
		param1.GetUInt64() != 0 ||
		param1.GetStr() != "" {
		t.Error("The Parameter outputs unexpected value")

		return
	}
}

func TestParameters(t *testing.T) {
	format := types.String("$((INT16)) $((INT32)) $((INT64)) $((UINT16)) " +
		"$((UINT32)) $((UINT64)) $((STRING)) $((BYTES))")

	labels := []types.String{
		"$((INT16))",
		"$((INT32))",
		"$((INT64))",
		"$((UINT16))",
		"$((UINT32))",
		"$((UINT64))",
		"$((STRING))",
		"$((BYTES))",
	}

	excepted := types.INT16_MAX_INT16.String().Join(" ",
		types.INT32_MAX_INT32.String(),
		" ",
		types.INT64_MAX_INT64.String(),
		" ",
		types.UINT16_MAX_UINT16.String(),
		" ",
		types.UINT32_MAX_UINT32.String(),
		" ",
		types.UINT64_MAX_UINT64.String(),
		" ",
		"TEST String ",
		"Test BYTES",
	)

	params := Parameters{}.
		AddInt16("INT16", types.INT16_MAX_INT16).
		AddInt32("INT32", types.INT32_MAX_INT32).
		AddInt64("INT64", types.INT64_MAX_INT64).
		AddUInt16("UINT16", types.UINT16_MAX_UINT16).
		AddUInt32("UINT32", types.UINT32_MAX_UINT32).
		AddUInt64("UINT64", types.UINT64_MAX_UINT64).
		AddString("STRING", "TEST String").
		AddBytes("BYTES", []byte("Test BYTES"))

	if params.Parse(format, labels) != excepted {
		t.Error("Parameters parsed a wrong result")

		return
	}
}

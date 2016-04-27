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

package server

import (
	"github.com/raincious/trap/trap/core/types"

	"testing"
)

func TestClientInfoSerialize(t *testing.T) {
	expected := []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0, 144, 31, 0, 84, 104, 105, 115,
		32, 105, 115, 32, 97, 32, 116, 101, 115, 116,
	}
	data := ClientInfo{}

	testIP, _ := types.ConvertIPFromString("0.0.0.0")

	data.Client = testIP
	data.Server.IP = testIP
	data.Server.Port = 8080
	data.Type = "This is a test"
	data.Marked = false

	result, _ := data.Serialize()

	if string(expected) != string(result) {
		t.Errorf("ClientInfo.Serialize() failed to serialize data. "+
			"Expecting %d, got %d", expected, result)

		return
	}
}

func TestClientInfoUnserialize(t *testing.T) {
	expected := []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0, 144, 31, 0, 84, 104, 105, 115,
		32, 105, 115, 32, 97, 32, 116, 101, 115, 116,
	}
	data := ClientInfo{}

	testIP, _ := types.ConvertIPFromString("0.0.0.0")

	data.Unserialize(expected)

	if !data.Client.IsEqual(&testIP) {
		t.Errorf("ClientInfo.Unserialize() failed to unserialize ClientInfo:"+
			"Expecting 'Client' to be '%d', got '%d'", testIP, data.Client)

		return
	}

	if !data.Server.IP.IsEqual(&testIP) || data.Server.Port != 8080 {
		t.Errorf("ClientInfo.Unserialize() failed to unserialize ClientInfo:"+
			"Expecting 'Server' to be '%s:%d', got '%s:%d'",
			testIP, 8080,
			data.Server.IP, data.Server.Port)

		return
	}

	if data.Type != "This is a test" {
		t.Errorf("ClientInfo.Unserialize() failed to unserialize ClientInfo:"+
			"Expecting 'Type' to be '%s', got '%s'",
			"This is a test", data.Type)

		return
	}

	if data.Marked != false {
		t.Errorf("ClientInfo.Unserialize() failed to unserialize ClientInfo:"+
			"Expecting 'Marked' to be '%v', got '%v'",
			false, data.Marked)

		return
	}
}

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

package permission

import (
	"github.com/raincious/trap/trap/core/types"

	"testing"
)

func TestRegisterNewRegister(t *testing.T) {
	reg := NewRegister()

	if reg.currentIdx != 0 || reg.registered == nil {
		t.Error("Inital value is incorrect")

		return
	}
}

func noItemIsTheSame(arr []types.UInt64) bool {
	arrmaps := map[types.UInt64]bool{}

	for _, val := range arr {
		if _, ok := arrmaps[val]; ok {
			return false
		}

		arrmaps[val] = true
	}

	return true
}

func TestRegisterGet(t *testing.T) {
	reg := NewRegister()
	vecs := []types.UInt64{}

	for i := types.UInt16(0); i < 64; i++ {
		vecs = append(vecs, reg.Get(types.UInt16(i).String()))
	}

	if !noItemIsTheSame(vecs) {
		t.Error("There is at least one conflict in the result vectors")

		return
	}
}

func TestRegisterAll(t *testing.T) {
	reg := NewRegister()
	vec := reg.Get("Name of the testing permission")

	if vec != reg.Get("Name of the testing permission") {
		t.Error("Failed to get permission vector")

		return
	}

	allPermissions := reg.All()

	if len(allPermissions) != 1 {
		t.Error("Permission register should only contain one item, but it's not")

		return
	}

	if allPermissions["Name of the testing permission"] != vec {
		t.Error("Regisrer.All() exports an invalid value")

		return
	}

	vec2 := reg.Get("Another permission key name")

	if len(allPermissions) != 2 {
		t.Error("Permission register should contain two items, but it's not")

		return
	}

	if allPermissions["Name of the testing permission"] != vec {
		t.Error("Regisrer.All() exports an invalid value")

		return
	}

	if allPermissions["Another permission key name"] != vec2 {
		t.Error("Regisrer.All() exports an invalid value")

		return
	}
}

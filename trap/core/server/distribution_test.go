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
	"testing"
)

func TestDistributionsGetSlot(t *testing.T) {
	dists := Distributions{}

	// Create a new slot
	slot := dists.GetSlot(16, "Test Port")

	if slot.Port != 16 || slot.Type != "Test Port" || slot.Hit != 0 {
		t.Error("Distributions.GetSlot() created an invalid new slot")

		return
	}

	// Update the slot, only update this one as it should be
	slot.Hit += 1

	// Re-get
	slot = dists.GetSlot(16, "Test Port")

	if slot.Hit != 1 {
		t.Error("Distributions.GetSlot() get an invalid slot")

		return
	}

	slot = dists.GetSlot(18, "Test Port 2")

	if slot.Port != 18 || slot.Type != "Test Port 2" || slot.Hit != 0 {
		t.Error("Distributions.GetSlot() created an invalid new slot")

		return
	}
}

func TestDistributionsDistributions(t *testing.T) {
	dists := Distributions{}

	if len(dists.Distributions()) != 0 {
		t.Error("Distributions.Distributions() not returning an empty slice")

		return
	}

	dists.GetSlot(2222, "Test Port 1")

	if len(dists.Distributions()) != 1 {
		t.Error("Distributions.Distributions() returning invalid amount of items")

		return
	}

	dists.GetSlot(2233, "Test Port 2")

	if len(dists.Distributions()) != 2 {
		t.Error("Distributions.Distributions() returning invalid amount of items")

		return
	}

	// Re-get, not create any new
	dists.GetSlot(2222, "Test Port 1")

	// Should be the same amount
	if len(dists.Distributions()) != 2 {
		t.Error("Distributions.Distributions() returning invalid amount of items")

		return
	}
}

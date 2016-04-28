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
	"time"
)

func TestHistoriesGetSlot(t *testing.T) {
	currentTime := time.Now()
	histories := Histories{}

	// Test to add invalid amount of items, it should loops around and
	// not cause any problem of that
	for i := types.Int64(1); i <= HISTORY_LENGTH*2; i++ {
		record := histories.GetSlot(
			currentTime.Add(-(time.Duration(i) * time.Hour)))

		if record.Hours != i+1 {
			t.Errorf("Histories.GetSlot() can't create slot correctly. "+
				"Expecting 'Hours' to be '%d', got '%d'",
				i, record.Hours)

			return
		}

		if record.Marked != 0 {
			t.Errorf("Histories.GetSlot() can't create slot correctly. "+
				"Expecting 'Marked' to be '%d', got '%d'",
				0, record.Marked)

			return
		}

		if record.Inbound != 0 {
			t.Errorf("Histories.GetSlot() can't create slot correctly. "+
				"Expecting 'Inbound' to be '%d', got '%d'",
				0, record.Inbound)

			return
		}

		if record.Hit != 0 {
			t.Errorf("Histories.GetSlot() can't create slot correctly. "+
				"Expecting 'Hit' to be '%d', got '%d'",
				0, record.Hit)

			return
		}
	}

	if types.Int64(len(histories.Histories())) != HISTORY_LENGTH {
		t.Error("Histories.GetSlot() creates unexpected amount of items")

		return
	}
}

func TestHistoriesHistories(t *testing.T) {
	currentTime := time.Now()
	histories := Histories{}

	if types.Int64(len(histories.Histories())) != 0 {
		t.Error("Histories.Histories() exports unexpected amount of items")

		return
	}

	histories.GetSlot(currentTime.Add(-(1 * time.Hour)))

	if types.Int64(len(histories.Histories())) != 1 {
		t.Error("Histories.Histories() exports unexpected amount of items")

		return
	}

	// Loop it back, we should get the same data as we created earlier on
	histories.GetSlot(currentTime.Add(
		-(time.Duration(HISTORY_LENGTH+1) * time.Hour)))

	if types.Int64(len(histories.Histories())) != 1 {
		t.Error("Histories.Histories() exports unexpected amount of items")

		return
	}

	// More one more step, we should now have to items exported
	histories.GetSlot(currentTime.Add(
		-(time.Duration(HISTORY_LENGTH+2) * time.Hour)))

	if types.Int64(len(histories.Histories())) != 2 {
		t.Error("Histories.Histories() exports unexpected amount of items")

		return
	}

	// And we should actually be able to get the item at the 2 hour slot
	// in every each circles
	histories.GetSlot(currentTime.Add(-(2 * time.Hour)))

	// No new item created (still export 2 items)
	if types.Int64(len(histories.Histories())) != 2 {
		t.Error("Histories.Histories() exports unexpected amount of items")

		return
	}
}

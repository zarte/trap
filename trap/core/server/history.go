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

	"math"
	"time"
)

const (
	HISTORY_LENGTH = types.Int64(12)
)

type History struct {
	Marked  types.UInt32
	Inbound types.UInt32
	Hit     types.UInt32

	Hours types.Int64
}

type Histories [HISTORY_LENGTH]*History

func (h *Histories) GetSlot(referenceTime time.Time) *History {
	hour := types.Int64(
		math.Ceil(
			time.Now().Sub(referenceTime).Hours()))
	slot := hour % HISTORY_LENGTH

	if h[slot] == nil {
		h[slot] = &History{
			Marked:  0,
			Inbound: 0,
			Hit:     0,
			Hours:   0,
		}
	}

	his := h[slot]

	if his.Hours != hour {
		his.Marked = 0
		his.Inbound = 0
		his.Hit = 0
		his.Hours = hour
	}

	return his
}

func (h *Histories) Histories() []History {
	historyLogs := []History{}

	for _, val := range h {
		if val == nil {
			continue
		}

		historyLogs = append(historyLogs, *val)
	}

	return historyLogs
}

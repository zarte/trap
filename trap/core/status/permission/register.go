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
)

const (
	registerFactor = types.UInt64(2)
)

type Register struct {
	currentIdx types.UInt64
	registered map[types.String]types.UInt64
}

func NewRegister() Register {
	return Register{
		currentIdx: 0,
		registered: map[types.String]types.UInt64{},
	}
}

func (r *Register) Get(name types.String) types.UInt64 {
	if _, ok := r.registered[name]; ok {
		return r.registered[name]
	}

	if r.currentIdx < 1 {
		r.currentIdx = 1
	}

	r.currentIdx = r.currentIdx * registerFactor
	r.registered[name] = r.currentIdx

	return r.registered[name]
}

func (r *Register) All() map[types.String]types.UInt64 {
	return r.registered
}

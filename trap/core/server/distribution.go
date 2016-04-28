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
)

type Distribution struct {
	Port types.UInt16
	Type types.String
	Hit  types.UInt32
}

type Distributions map[types.String]*Distribution

func (d Distributions) GetSlot(
	port types.UInt16, typeName types.String) *Distribution {
	pType := port.String() + ":" + typeName

	if _, ok := d[pType]; !ok {
		d[pType] = &Distribution{
			Port: port,
			Type: typeName,
			Hit:  0,
		}
	}

	return d[pType]
}

func (d Distributions) Distributions() []Distribution {
	result := []Distribution{}

	for _, val := range d {
		result = append(result, *val)
	}

	return result
}

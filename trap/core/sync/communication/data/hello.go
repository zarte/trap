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

package data

import (
	"github.com/raincious/trap/trap/core/types"
)

type Hello struct {
	Base

	MaxLength types.UInt16
	Password  types.String
	Connected types.IPAddresses
}

func (d *Hello) Parse(msg [][]byte) *types.Throw {
	verifyErr := d.Verify(msg, 3)

	if verifyErr != nil {
		return verifyErr
	}

	maxlenErr := d.MaxLength.Unserialize(msg[0])

	if maxlenErr != nil {
		return maxlenErr
	}

	d.Password = types.String(msg[1])

	connectMarErr := d.Connected.Unserialize(msg[2])

	if connectMarErr != nil {
		return connectMarErr
	}

	return nil
}

func (d *Hello) Build() ([][]byte, *types.Throw) {
	maxlen, maxlenErr := d.MaxLength.Serialize()

	if maxlenErr != nil {
		return [][]byte{}, maxlenErr
	}

	ipByte, cIPErr := d.Connected.Serialize()

	if cIPErr != nil {
		return [][]byte{}, cIPErr
	}

	return [][]byte{maxlen, d.Password.Bytes(), ipByte}, nil
}

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

type HelloConflict struct {
	Base

	Confilct types.IPAddresses
}

func (d *HelloConflict) Parse(msg [][]byte) *types.Throw {
	verifyErr := d.Verify(msg, 1)

	if verifyErr != nil {
		return verifyErr
	}

	connectedErr := d.Confilct.Unserialize(msg[0])

	if connectedErr != nil {
		return connectedErr
	}

	return nil
}

func (d *HelloConflict) Build() ([][]byte, *types.Throw) {
	ipByte, cIPErr := d.Confilct.Serialize()

	if cIPErr != nil {
		return [][]byte{}, cIPErr
	}

	return [][]byte{
		ipByte,
	}, nil
}

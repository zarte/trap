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

var (
	ErrDataInvalidLengthForIP *types.Error = types.NewError(
		"Data length is invalid for an IP address bytes")
)

type ClientUnmark struct {
	Base

	Addresses []types.IP
}

func (d *ClientUnmark) Parse(msg [][]byte) *types.Throw {
	verifyErr := d.Verify(msg, 1)

	if verifyErr != nil {
		return verifyErr
	}

	for _, data := range msg {
		if len(data) != types.IP_ADDR_SLICE_LEN {
			return ErrDataInvalidLengthForIP.Throw()
		}

		newIP := types.IP{}

		copy(newIP[:], data[:])

		d.Addresses = append(d.Addresses, newIP)
	}

	return nil
}

func (d *ClientUnmark) Build() ([][]byte, *types.Throw) {
	result := [][]byte{}

	for _, addr := range d.Addresses {
		result = append(result, addr[:])
	}

	return result, nil
}

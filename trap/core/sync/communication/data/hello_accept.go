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

	"time"
)

type HelloAccept struct {
	Base

	MaxLength      types.UInt16
	HeatBeatPeriod time.Duration
	Timeout        time.Duration
	Connected      types.IPAddresses
}

func (d *HelloAccept) Parse(msg [][]byte) *types.Throw {
	heatbeatPeriod := types.Int64(0)
	timeout := types.Int64(0)
	verifyErr := d.Verify(msg, 4)

	if verifyErr != nil {
		return verifyErr
	}

	maxlenErr := d.MaxLength.Unserialize(msg[0])

	if maxlenErr != nil {
		return maxlenErr
	}

	hbpErr := heatbeatPeriod.Unserialize(msg[1])

	if hbpErr != nil {
		return hbpErr
	}

	d.HeatBeatPeriod = time.Duration(heatbeatPeriod)

	timeoutErr := timeout.Unserialize(msg[2])

	if timeoutErr != nil {
		return timeoutErr
	}

	d.Timeout = time.Duration(timeout)

	connectedErr := d.Connected.Unserialize(msg[3])

	if connectedErr != nil {
		return connectedErr
	}

	return nil
}

func (d *HelloAccept) Build() ([][]byte, *types.Throw) {
	maxlen, maxlenErr := d.MaxLength.Serialize()

	if maxlenErr != nil {
		return [][]byte{}, maxlenErr
	}

	timeBytes, timeBErr := types.Int64(d.HeatBeatPeriod).Serialize()

	if timeBErr != nil {
		return [][]byte{}, timeBErr
	}

	timeoutBytes, timeOutBErr := types.Int64(d.Timeout).Serialize()

	if timeOutBErr != nil {
		return [][]byte{}, timeOutBErr
	}

	ipByte, cIPErr := d.Connected.Serialize()

	if cIPErr != nil {
		return [][]byte{}, cIPErr
	}

	return [][]byte{
		maxlen,
		timeBytes,
		timeoutBytes,
		ipByte,
	}, nil
}

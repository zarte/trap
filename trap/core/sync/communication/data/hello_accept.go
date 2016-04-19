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

	HeatBeatPeriod time.Duration
	Timeout        time.Duration
	Connected      types.IPAddresses
}

func (d *HelloAccept) Parse(msg [][]byte) *types.Throw {
	heatbeatPeriod := types.Int64(0)
	timeout := types.Int64(0)
	verifyErr := d.Verify(msg, 3)

	if verifyErr != nil {
		return verifyErr
	}

	connectedErr := d.Connected.Unserialize(msg[2])

	if connectedErr != nil {
		return connectedErr
	}

	hbpErr := heatbeatPeriod.Unserialize(msg[0])

	if hbpErr != nil {
		return hbpErr
	}

	d.HeatBeatPeriod = time.Duration(heatbeatPeriod)

	timeoutErr := timeout.Unserialize(msg[1])

	if timeoutErr != nil {
		return timeoutErr
	}

	d.Timeout = time.Duration(timeout)

	return nil
}

func (d *HelloAccept) Build() ([][]byte, *types.Throw) {
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
		timeBytes,
		timeoutBytes,
		ipByte,
	}, nil
}

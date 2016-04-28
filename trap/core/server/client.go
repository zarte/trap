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

var (
	ErrClientInfoInvalidLength *types.Error = types.NewError(
		"Invalid length for unserialization of the Client Info")
)

const (
	CLIENT_INFO_SEGMENT_CLIENT_BEGIN = 0
	CLIENT_INFO_SEGMENT_CLIENT_END   = types.IP_ADDR_SLICE_LEN

	CLIENT_INFO_SEGMENT_SERVER_BEGIN = types.IP_ADDR_SLICE_LEN
	CLIENT_INFO_SEGMENT_SERVER_END   = CLIENT_INFO_SEGMENT_SERVER_BEGIN +
		types.IP_ADDR_SLICE_LEN +
		types.IP_PORT_LEN

	CLIENT_INFO_SEGMENT_INDEX = CLIENT_INFO_SEGMENT_SERVER_END + 1
)

type ClientInfo struct {
	Client types.IP
	Server types.IPAddress
	Type   types.String
	Marked bool
}

func (c *ClientInfo) Serialize() ([]byte, *types.Throw) {
	result := []byte{}

	serverByte, serErr := c.Server.Serialize()

	if serErr != nil {
		return []byte{}, serErr
	}

	result = append(result, c.Client[:]...)
	result = append(result, serverByte...)

	if c.Marked {
		result = append(result, ^byte(0))
	} else {
		result = append(result, byte(0))
	}

	result = append(result, []byte(c.Type)...)

	return result, nil
}

func (c *ClientInfo) Unserialize(data []byte) *types.Throw {
	// types.IP + types.IPAddress + Marked (1 byte)
	if len(data) < CLIENT_INFO_SEGMENT_INDEX {
		return ErrClientInfoInvalidLength.Throw()
	}

	copy(c.Client[:],
		data[CLIENT_INFO_SEGMENT_CLIENT_BEGIN:CLIENT_INFO_SEGMENT_CLIENT_END])

	serverErr := c.Server.Unserialize(
		data[CLIENT_INFO_SEGMENT_SERVER_BEGIN:CLIENT_INFO_SEGMENT_SERVER_END])

	if serverErr != nil {
		return serverErr
	}

	if data[CLIENT_INFO_SEGMENT_INDEX] != byte(0) {
		c.Marked = true
	} else {
		c.Marked = false
	}

	c.Type = types.String(data[CLIENT_INFO_SEGMENT_INDEX:])

	return nil
}

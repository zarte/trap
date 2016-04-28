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

package net

import (
	"github.com/raincious/trap/trap/core/types"

	"net"
)

var (
	ErrInvalidIPAddress *types.Error = types.NewError(
		"'%s' is not an IP address")

	ErrInvalidPort *types.Error = types.NewError(
		"'%s' is not an invalid port")
)

type Net struct {
}

func (n *Net) ParseConfig(
	cfg types.String) (net.IP, types.UInt16, *types.Throw) {
	var ip net.IP

	portStr, ipAddrStr := cfg.SpiltWith("@")

	port := portStr.UInt16()

	if port <= 0 {
		return net.IP{}, 0, ErrInvalidPort.Throw(portStr)
	}

	if ipAddrStr != "" {
		ip = net.ParseIP(ipAddrStr.String())

		if ip == nil {
			return net.IP{}, 0, ErrInvalidIPAddress.Throw(ipAddrStr)
		}
	} else {
		ip = net.ParseIP("0.0.0.0")
	}

	return ip, port, nil
}

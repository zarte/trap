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

package conn

import (
	"net"
	"time"
)

type Conn struct {
	net.Conn

	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (c *Conn) SetReadTimeout(timeout time.Duration) {
	c.readTimeout = timeout
}

func (c *Conn) SetWriteTimeout(timeout time.Duration) {
	c.writeTimeout = timeout
}

func (c *Conn) SetTimeout(timeout time.Duration) {
	c.readTimeout = timeout
	c.writeTimeout = timeout
}

func (c *Conn) Read(buf []byte) (int, error) {
	e := c.Conn.SetReadDeadline(
		time.Now().Add(c.readTimeout))

	if e != nil {
		return 0, e
	}

	return c.Conn.Read(buf)
}

func (c *Conn) Write(buf []byte) (int, error) {
	e := c.Conn.SetReadDeadline(
		time.Now().Add(c.writeTimeout))

	if e != nil {
		return 0, e
	}

	return c.Conn.Write(buf)
}

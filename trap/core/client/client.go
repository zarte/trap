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

package client

import (
	"github.com/raincious/trap/trap/core/types"

	"net"
	"time"
)

const (
	CLIENT_EXPIRED_NO = iota
	CLIENT_EXPIRED_YES
	CLIENT_EXPIRED_RESTRICTED
)

type Client struct {
	address        net.IP
	firstSeen      time.Time
	lastSeen       time.Time
	count          types.UInt32
	records        []Record
	lastRecord     *Record
	marked         bool
	onMark         func(*Client, MarkType)
	onUnmark       func(*Client, UnmarkType)
	onRecord       func(*Client, Record)
	tolerateCount  types.UInt32
	tolerateExpire time.Duration
	restrictExpire time.Duration
}

func (c *Client) Address() net.IP {
	return c.address
}

func (c *Client) FirstSeen() time.Time {
	return c.firstSeen
}

func (c *Client) LastSeen() time.Time {
	return c.lastSeen
}

func (c *Client) Count() types.UInt32 {
	return c.count
}

func (c *Client) Marked() bool {
	return c.marked
}

func (c *Client) Record(record Record, maxLen types.UInt16) {
	dataLen := types.Int32(len(c.records)).UInt16()

	c.records = append(c.records, record)

	if dataLen > maxLen {
		c.records = c.records[dataLen-maxLen:]
	}

	c.lastRecord = &c.records[len(c.records)-1]

	c.onRecord(c, record)
}

func (c *Client) Records() []Record {
	return c.records
}

func (c *Client) LastRecord() *Record {
	return c.lastRecord
}

func (c *Client) Mark(ty MarkType) {
	oldMarkStatus := c.marked

	c.marked = true

	if !oldMarkStatus {
		c.onMark(c, ty)
	}
}

func (c *Client) Unmark(ty UnmarkType) {
	oldUnmarkStatus := c.marked

	c.marked = false

	if oldUnmarkStatus {
		c.onUnmark(c, ty)
	}
}

func (c *Client) Bump() {
	c.lastSeen = time.Now()

	if c.count+1 > types.UINT32_MAX_UINT32 {
		return
	}

	c.count += 1
}

func (c *Client) Rebump() {
	c.count = 1
	c.lastSeen = time.Now()
}

func (c *Client) Tolerate(count types.UInt32, expire time.Duration,
	restrict time.Duration) {
	c.tolerateCount = count
	c.tolerateExpire = expire
	c.restrictExpire = restrict
}

func (c *Client) Expired(now time.Time) int {
	expireTime := c.lastSeen.Add(c.tolerateExpire)

	if !now.After(expireTime) {
		return CLIENT_EXPIRED_NO
	}

	restrictTime := expireTime.Add(c.restrictExpire)

	if c.count >= c.tolerateCount && !now.After(restrictTime) {
		return CLIENT_EXPIRED_RESTRICTED
	}

	return CLIENT_EXPIRED_YES
}

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

type Client struct {
    Address                 net.IP

    LastSeen                time.Time
    Count                   types.UInt32

    Data                    []Data

    Marked                  bool
}

func (c *Client) Bump() {
    if c.Count + 1 > types.UINT32_MAX_UINT32 {
        return
    }

    c.Count                 +=  1
    c.LastSeen              =   time.Now()
}

func (c *Client) AppendData(data Data, maxLen types.UInt16) {
    dataLen                 :=  types.UInt16(len(c.Data))

    c.Data                  =   append(c.Data, data)

    if dataLen > maxLen {
        c.Data              =   c.Data[dataLen - maxLen:]
    }
}

type Clients map[types.IP]*Client

func (c Clients) Get(ip types.IP) (*Client) {
    if _, ok := c[ip]; !ok {
        c[ip]               =   &Client{
            Address:            ip.IP(),
            LastSeen:           time.Now(),
            Count:              0,
            Data:               []Data{},
            Marked:             false,
        }
    }

    return c[ip]
}

func (c Clients) Has(ip types.IP) (bool) {
    if _, ok := c[ip]; !ok {
        return false
    }

    return true
}

func (c Clients) Delete(ip types.IP) (*types.Throw) {
    if !c.Has(ip) {
        return ErrClientNotFound.Throw(ip)
    }

    delete(c, ip)

    return nil
}

func (c Clients) Export() ([]Client) {
    clients                 :=  []Client{}

    for _, clientInfo := range c {
        clients             =   append(clients, *clientInfo)
    }

    return clients
}
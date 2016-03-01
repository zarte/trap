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

    "time"
    "net"
)

type Client struct {
    address                 net.IP

    firstSeen               time.Time
    lastSeen                time.Time

    count                   types.UInt32

    data                    []Data

    marked                  bool

    onMark                  func(*Client)
    onUnmark                func(*Client)
    onRecord                func(*Client, Data)
}

func (c *Client) Address() (net.IP) {
    return c.address
}

func (c *Client) FirstSeen() (time.Time) {
    return c.firstSeen
}

func (c *Client) LastSeen() (time.Time) {
    return c.lastSeen
}

func (c *Client) Count() (types.UInt32) {
    return c.count
}

func (c *Client) Data() ([]Data) {
    return c.data
}

func (c *Client) Marked() (bool) {
    return c.marked
}

func (c *Client) Mark() {
    oldMarkStatus           :=  c.marked

    c.marked                =   true

    if !oldMarkStatus {
        c.onMark(c)
    }
}

func (c *Client) Unmark() {
    oldUnmarkStatus         :=  c.marked

    c.marked                =   false

    if oldUnmarkStatus {
        c.onUnmark(c)
    }
}

func (c *Client) Bump() {
    if c.count + 1 > types.UINT32_MAX_UINT32 {
        return
    }

    c.count                 +=  1
    c.lastSeen              =   time.Now()
}

func (c *Client) AppendData(data Data, maxLen types.UInt16) {
    dataLen                 :=  types.Int32(len(c.data)).UInt16()

    c.data                  =   append(c.data, data)

    if dataLen > maxLen {
        c.data              =   c.data[dataLen - maxLen:]
    }

    c.onRecord(c, data)
}

type Clients struct {
    clients                 map[types.IP]*Client

    onMark                  func(*Client)
    onUnmark                func(*Client)
    onRecord                func(*Client, Data)
}

func NewClients(config *Config) (*Clients) {
    return &Clients{
        clients:            map[types.IP]*Client{},
        onMark:             config.OnMark,
        onUnmark:           config.OnUnmark,
        onRecord:           config.OnRecord,
    }
}

func (c *Clients) Get(ip types.IP) (*Client, bool) {
    isNew                   :=  false

    if _, ok := c.clients[ip]; !ok {
        c.clients[ip]       =   &Client{
            address:            ip.IP(),
            firstSeen:          time.Now(),
            lastSeen:           time.Now(),
            count:              0,
            data:               []Data{},
            marked:             false,
            onMark:             c.onMark,
            onUnmark:           c.onUnmark,
            onRecord:           c.onRecord,
        }

        isNew               =   true
    }

    return c.clients[ip], isNew
}

func (c *Clients) Has(ip types.IP) (bool) {
    if _, ok := c.clients[ip]; !ok {
        return false
    }

    return true
}

func (c *Clients) Len() (int) {
    return len(c.clients)
}

func (c *Clients) Delete(ip types.IP) (*types.Throw) {
    if !c.Has(ip) {
        return ErrClientNotFound.Throw(ip)
    }

    if c.clients[ip].marked {
        c.clients[ip].Unmark()
    }

    delete(c.clients, ip)

    return nil
}

func (c *Clients) Scan(callback func(types.IP, *Client) (*types.Throw)) (*types.Throw) {
    var err *types.Throw    =   nil

    for key, val := range c.clients {
        callbackErr         :=  callback(key, val)

        if callbackErr == nil {
            continue
        }

        err = callbackErr

        break
    }

    return err
}

func (c *Clients) Clear() (*types.Throw) {
    var err *types.Throw    =   nil

    for key, _ := range c.clients {
        deleteErr           :=  c.Delete(key)

        if deleteErr != nil {
            err = deleteErr
        }
    }

    return err
}

func (c *Clients) Export() ([]ClientExport) {
    clients                 :=  []ClientExport{}

    for _, clientInfo := range c.clients {
        clients             =   append(clients, ClientExport{
            Address:            clientInfo.Address(),
            FirstSeen:          clientInfo.FirstSeen(),
            LastSeen:           clientInfo.LastSeen(),
            Count:              clientInfo.Count(),
            Data:               clientInfo.Data(),
            Marked:             clientInfo.Marked(),
        })
    }

    return clients
}
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

    "testing"
    "net"
    "time"
)

func TestClientAssign(t *testing.T) {
    nw      :=  time.Now()

    client  :=  Client{
        address:            net.ParseIP("127.0.0.1"),
        firstSeen:          nw,
        lastSeen:           nw.Add(10 * time.Second),
        count:              0,
        records:            []Record{},
        marked:             false,
        onMark:             func(c *Client) {

        },
        onUnmark:           func(c *Client) {

        },
        onRecord:           func(c *Client, d Record) {

        },
        tolerateCount:      0,
        tolerateExpire:     time.Duration(0),
        restrictExpire:     time.Duration(0),
    }

    if client.Address().String() != "127.0.0.1" {
        t.Error("Unexpected Address")

        return
    }

    if client.FirstSeen() != nw {
        t.Error("Unexpected FirstSeen")

        return
    }

    if client.LastSeen() != nw.Add(10 * time.Second) {
        t.Error("Unexpected LastSeen")

        return
    }

    if client.Count() != 0 {
        t.Error("Unexpected Count")

        return
    }

    if len(client.Records()) != 0 {
        t.Error("Unexpected Data")

        return
    }

    if client.Marked() != false {
        t.Error("Unexpected Data")

        return
    }
}

func TestClientDataRecord(t *testing.T) {
    now             :=  time.Now()
    recordCalled    :=  false
    client          :=  Client{
        address:            net.ParseIP("127.0.0.1"),
        firstSeen:          now,
        lastSeen:           now.Add(10 * time.Second),
        count:              0,
        records:            []Record{},
        marked:             false,
        onMark:             func(c *Client) {

        },
        onUnmark:           func(c *Client) {

        },
        onRecord:           func(c *Client, d Record) {
            recordCalled    =   true
        },
        tolerateCount:      0,
        tolerateExpire:     time.Duration(0),
        restrictExpire:     time.Duration(0),
    }

    ip, ipErr   :=  types.ConvertIPFromString("127.0.0.1")

    if ipErr != nil {
        t.Error("Unexpected error happened when trying to " +
            "convert IP address from a string")

        return
    }

    client.Record(Record{
        Inbound:        []byte("TEST BYTES Inbound"),
        Outbound:       []byte("TEST BYTES Outbound"),
        Hitting:        Hitting{
            IPAddress:  types.IPAddress{
                IP:     ip,
                Port:   8080,
            },
            Type:       "Test",
        },
        Time:           now,
    }, 128)

    if !recordCalled {
        t.Error("`onRecord` callback hasn't been call")

        return
    }

    records := client.Records()

    if len(records) != 1 {
        t.Error("Total amount of exported records is not as expected")

        return
    }

    if string(records[0].Inbound) != "TEST BYTES Inbound" ||
        string(records[0].Outbound) != "TEST BYTES Outbound" ||
        !records[0].Hitting.IPAddress.IP.IsEqual(&ip) ||
        records[0].Hitting.IPAddress.Port != 8080 ||
        records[0].Hitting.Type != "Test" ||
        records[0].Time != now {
        t.Error("Exported data is not the original")

        return
    }
}

func TestClientMarkUnmark(t *testing.T) {
    now             :=  time.Now()
    markCalled      :=  false
    unmarkCalled    :=  false
    client          :=  Client{
        address:            net.ParseIP("127.0.0.1"),
        firstSeen:          now,
        lastSeen:           now.Add(10 * time.Second),
        count:              0,
        records:            []Record{},
        marked:             false,
        onMark:             func(c *Client) {
            markCalled      =   true
        },
        onUnmark:           func(c *Client) {
            unmarkCalled    =   true
        },
        onRecord:           func(c *Client, d Record) {},
        tolerateCount:      0,
        tolerateExpire:     time.Duration(0),
        restrictExpire:     time.Duration(0),
    }

    if client.Marked() {
        t.Error("The default status of `Marked` must be 'false'")

        return
    }

    client.Mark()

    if !markCalled {
        t.Error("`onMark` callback hasn't been call")

        return
    }

    if !client.Marked() {
        t.Error("The client should be marked by now, but it's not " +
            "what's happened")

        return
    }

    client.Unmark()

    if !unmarkCalled {
        t.Error("`onUnmark` callback hasn't been call")

        return
    }

    if client.Marked() {
        t.Error("The client should be unmarked by now, but it's not " +
            "what's happened")

        return
    }
}

func TestClientBumpRebump(t *testing.T) {
    now             :=  time.Now()
    client          :=  Client{
        address:            net.ParseIP("127.0.0.1"),
        firstSeen:          now,
        lastSeen:           now.Add(1 * time.Second),
        count:              0,
        records:            []Record{},
        marked:             false,
        onMark:             func(c *Client) {},
        onUnmark:           func(c *Client) {},
        onRecord:           func(c *Client, d Record) {},
        tolerateCount:      0,
        tolerateExpire:     time.Duration(0),
        restrictExpire:     time.Duration(0),
    }

    if client.Count() != 0 {
        t.Error("The initial value of `Count` should be '0'")

        return
    }

    if client.LastSeen() != now.Add(1 * time.Second) {
        t.Error("The initial value of `LastSeen` is invalid")

        return
    }

    time.Sleep(3 * time.Second)

    client.Bump()

    if client.Count() != 1 {
        t.Error("The new value of `Count` should be '1'")

        return
    }

    if !client.LastSeen().After(now.Add(1 * time.Second)) {
        t.Errorf("The new value of `LastSeen` should be laster than %s",
            now.Add(1 * time.Second))

        return
    }

    time.Sleep(3 * time.Second)

    client.Bump()

    if client.Count() != 2 {
        t.Error("The new value of `Count` should be '2'")

        return
    }

    if !client.LastSeen().After(now.Add(6 * time.Second)) {
        t.Errorf("The new value of `LastSeen` should be laster than %s",
            now.Add(6 * time.Second))

        return
    }

    client.Rebump()

    if client.Count() != 1 {
        t.Error("The new value of `Count` should be '1'")

        return
    }
}

func TestClientTolerateExpired(t *testing.T) {
    now             :=  time.Now()
    client          :=  Client{
        address:            net.ParseIP("127.0.0.1"),
        firstSeen:          now,
        lastSeen:           now,
        count:              0,
        records:            []Record{},
        marked:             false,
        onMark:             func(c *Client) {},
        onUnmark:           func(c *Client) {},
        onRecord:           func(c *Client, d Record) {},
        tolerateCount:      0,
        tolerateExpire:     time.Duration(0),
        restrictExpire:     time.Duration(0),
    }

    client.Tolerate(3, 3 * time.Second, 3 * time.Second)

    time.Sleep(1 * time.Second)

    client.Bump()

    if client.Expired(time.Now()) {
        t.Error("The client should not be expired for now")

        return
    }

    if !client.Expired(time.Now().Add(3 * time.Second)) {
        t.Error("The client should be expired by now")

        return
    }

    client.Bump()
    client.Bump()

    if client.Expired(time.Now().Add(5 * time.Second)) {
        t.Error("The client should not be expired for now")

        return
    }
}
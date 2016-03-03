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
    //"net"
    //"time"
)

var (
    errFakeError *types.Error = types.NewError("This is a fake error")
)

func TestClientsNewClients(t *testing.T) {
    clients := NewClients(Config{
        OnMark:             func(*Client) {},
        OnUnmark:           func(*Client) {},
        OnRecord:           func(*Client, Record) {},
    })

    if clients.Len() != 0 {
        t.Error("Invalid number of a completely new Client Map")

        return
    }
}

func TestClientsGet(t *testing.T) {
    clients         :=  NewClients(Config{
        OnMark:         func(*Client) {},
        OnUnmark:       func(*Client) {},
        OnRecord:       func(*Client, Record) {},
    })
    ip, ipErr       :=  types.ConvertIPFromString("127.0.0.1")

    if ipErr != nil {
        t.Error("Failed to convert IP address from string")

        return
    }

    client, isNew   :=  clients.Get(ip)

    if !isNew {
        t.Error("Failed asserting a new client is new")

        return
    }

    if !client.Address().Equal(ip.IP()) {
        t.Error("The new client contains invalid IP address")

        return
    }

    if clients.Len() != 1 {
        t.Error("Failed asserting '1' client in the list")

        return
    }

    if !clients.Has(ip) {
        t.Error("Can't found a client which is existed")

        return
    }
}

func TestClientsHas(t *testing.T) {
    clients         :=  NewClients(Config{
        OnMark:         func(*Client) {},
        OnUnmark:       func(*Client) {},
        OnRecord:       func(*Client, Record) {},
    })
    ip1, ipErr1     :=  types.ConvertIPFromString("127.0.0.1")
    ip2, ipErr2     :=  types.ConvertIPFromString("127.0.0.2")

    if ipErr1 != nil || ipErr2 != nil {
        t.Error("Failed to convert IP address from string")

        return
    }

    clients.Get(ip1)

    if !clients.Has(ip1) {
        t.Error("Failed asserting an existed client is existed")

        return
    }

    if clients.Has(ip2) {
        t.Error("Failed asserting an not existed client is not existed")

        return
    }

    if clients.Len() != 1 {
        t.Error("Has method can't change length of a client map")

        return
    }
}

func TestClientsDelete(t *testing.T) {
    clients         :=  NewClients(Config{
        OnMark:         func(*Client) {},
        OnUnmark:       func(*Client) {},
        OnRecord:       func(*Client, Record) {},
    })
    ip, ipErr       :=  types.ConvertIPFromString("127.0.0.1")

    if ipErr != nil {
        t.Error("Failed to convert IP address from string")

        return
    }

    clients.Get(ip)

    if !clients.Has(ip) {
        t.Error("Failed asserting an existed client is existed")

        return
    }

    deleteErr       :=  clients.Delete(ip)

    if deleteErr != nil {
        t.Errorf("Can't delete client due to error: %s", deleteErr)

        return
    }

    deleteErr       =   clients.Delete(ip)

    if deleteErr == nil || !deleteErr.Is(ErrClientNotFound) {
        t.Errorf("Unexpected error when trying to delete a non-existed client: %s",
            deleteErr)

        return
    }
}

func TestClientsScan(t *testing.T) {
    scaned          :=  0
    clients         :=  NewClients(Config{
        OnMark:         func(*Client) {},
        OnUnmark:       func(*Client) {},
        OnRecord:       func(*Client, Record) {},
    })

    clients.Scan(func(ip types.IP, client *Client) (*types.Throw) {
        scaned++

        return nil
    })

    if scaned != 0 {
        t.Error("Unexpected scanning count")

        return
    }

    ip1, ipErr1     :=  types.ConvertIPFromString("127.0.0.1")

    if ipErr1 != nil {
        t.Errorf("Failed to convert string to IP due to error: %s", ipErr1)

        return
    }

    clients.Get(ip1)

    scaned = 0

    clients.Scan(func(ip types.IP, client *Client) (*types.Throw) {
        scaned++

        return nil
    })

    if scaned != 1 {
        t.Error("Unexpected scanning count")

        return
    }

    ip2, ipErr2     :=  types.ConvertIPFromString("127.0.0.2")

    if ipErr2 != nil {
        t.Errorf("Failed to convert string to IP due to error: %s", ipErr2)

        return
    }

    clients.Get(ip2)

    scaned = 0

    clients.Scan(func(ip types.IP, client *Client) (*types.Throw) {
        scaned++

        return nil
    })

    if scaned != 2 {
        t.Error("Unexpected scanning count")

        return
    }

    scaned = 0

    scanErr := clients.Scan(func(ip types.IP, client *Client) (*types.Throw) {
        scaned++

        if scaned >= 1 {
            return errFakeError.Throw()
        }

        return nil
    })

    if scaned != 1 {
        t.Error("Unexpected scanning count")

        return
    }

    if scanErr == nil || !scanErr.Is(errFakeError) {
        t.Error("Unexpected error when scanning clients map")

        return
    }
}

func TestClientsClear(t *testing.T) {
    markCalled      :=  false
    unmarkCalled    :=  false
    clients         :=  NewClients(Config{
        OnMark:         func(*Client) {
            markCalled      = true
        },
        OnUnmark:       func(*Client) {
            unmarkCalled    = true
        },
        OnRecord:       func(*Client, Record) {},
    })

    ip1, ipErr1     :=  types.ConvertIPFromString("127.0.0.1")
    ip2, ipErr2     :=  types.ConvertIPFromString("127.0.0.2")

    if ipErr1 != nil || ipErr2 != nil {
        t.Error("Failed to convert string to IP due to error")

        return
    }

    clients.Get(ip1)
    client2, _      :=  clients.Get(ip2)

    client2.Mark()

    if clients.Len() != 2 {
        t.Error("Unexpected amount of clients")

        return
    }

    if clients.Len() != 2 {
        t.Error("Unexpected amount of clients")

        return
    }

    clearErr        :=  clients.Clear()

    if clearErr != nil {
        t.Errorf("Failed clear clients due to error: %s", clearErr)

        return
    }

    if clients.Len() != 0 {
        t.Error("Unexpected amount of clients")

        return
    }

    if !markCalled || !unmarkCalled {
        t.Error("Expected callback hasn't been call")

        return
    }
}

func TestClientsExport(t *testing.T) {
    clients         :=  NewClients(Config{
        OnMark:         func(*Client) {},
        OnUnmark:       func(*Client) {},
        OnRecord:       func(*Client, Record) {},
    })

    ip1, ipErr1     :=  types.ConvertIPFromString("127.0.0.1")
    ip2, ipErr2     :=  types.ConvertIPFromString("127.0.0.2")

    if ipErr1 != nil || ipErr2 != nil {
        t.Error("Failed to convert string to IP due to error")

        return
    }

    clients.Get(ip1)
    clients.Get(ip2)

    exported := clients.Export()

    if len(exported) != 2 {
        t.Error("Unexpected amount of exported clients")

        return
    }

    if !exported[0].Address.Equal(ip1.IP()) ||
        !exported[1].Address.Equal(ip2.IP()) {
        t.Error("Failed to convert string to IP due to error")

        return
    }
}
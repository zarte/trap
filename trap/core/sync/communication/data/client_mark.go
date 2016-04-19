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
    "github.com/raincious/trap/trap/core/server"
)

type ClientMark struct {
    Base

    Addresses           []server.ClientInfo
}

func (d *ClientMark) Parse(msg [][]byte) (*types.Throw) {
    verifyErr           :=  d.Verify(msg, 1)

    if verifyErr != nil {
        return verifyErr
    }

    for _, data := range msg {
        clientInfo      :=  server.ClientInfo{}

        clientSerErr    :=  clientInfo.Unserialize(data)

        if clientSerErr != nil {
            return clientSerErr
        }

        d.Addresses     =   append(d.Addresses, clientInfo)
    }

    return nil
}

func (d *ClientMark) Build() ([][]byte, *types.Throw) {
    result              :=  [][]byte{}

    for _, addr := range d.Addresses {
        clientBy, cErr  :=  addr.Serialize()

        if cErr != nil {
            return [][]byte{}, cErr
        }

        result          =   append(result, clientBy)
    }

    return result, nil
}

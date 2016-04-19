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

package messager

import (
    "github.com/raincious/trap/trap/core/sync/communication/conn"
    "github.com/raincious/trap/trap/core/types"

    "time"
    "net"
)

var (
    ErrRequestNotRespondable *types.Error =
        types.NewError("`Sync` Request for '%d' is not respondable")
)

type Request struct {
    conn            *conn.Conn
    messager        *Messager

    data            [][]byte
    dataLen         uint
    code            byte
    id              byte

    isReplyable     bool
}

func (r *Request) Length() (uint) {
    return r.dataLen
}

func (r *Request) RemoteAddr() (net.Addr) {
    return r.conn.RemoteAddr()
}

func (r *Request) LocalAddr() (net.Addr) {
    return r.conn.LocalAddr()
}

func (r *Request) Conn() (*conn.Conn) {
    return r.conn
}

func (r *Request) ID() (byte) {
    return r.id
}

func (r *Request) Code() (byte) {
    return r.code
}

func (r *Request) Data() ([][]byte) {
    return r.data
}

func (r *Request) Stats() (Stats) {
    return r.messager.Stats()
}

func (r *Request) Close() (*types.Throw) {
    err :=  r.conn.Close()

    if err != nil {
        return types.ConvertError(err)
    }

    return nil
}

func (r *Request) Reply(code byte, data Data) (*types.Throw) {
    if !r.isReplyable {
        return ErrRequestNotRespondable.Throw(r.code)
    }

    return r.messager.Reply(r.id, code, data)
}

func (r *Request) Query(code byte, data Data, responds Callbacks,
    maxRespondLen uint, waitTime time.Duration) (*types.Throw) {
    return r.messager.Query(code, data, responds, maxRespondLen, waitTime)
}
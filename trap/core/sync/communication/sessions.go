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

package communication

import (
    "github.com/raincious/trap/trap/core/sync/communication/messager"
    "github.com/raincious/trap/trap/core/sync/communication/conn"
    "github.com/raincious/trap/trap/core/types"

    "time"
    "net"
)

type Sessions struct {
    defaultResponders       messager.Callbacks

    reqTimeout              time.Duration
    sessions                map[string]*Session
    sessionLock             types.Mutex

    onRegister              func()
    onUnregister            func()
}

func NewSessions(defaultResponders messager.Callbacks, reqTimeout time.Duration,
    onRegister, onUnregister func()) (*Sessions) {
    return &Sessions{
        defaultResponders:  defaultResponders,

        reqTimeout:         reqTimeout,
        sessions:           map[string]*Session{},
        sessionLock:        types.Mutex{},

        onRegister:         onRegister,
        onUnregister:       onUnregister,
    }
}

func (s *Sessions) Has(conn net.Conn) (bool) {
    var hasIt bool      =   false

    s.sessionLock.Exec(func() {
        hasIt           =   s.hasByKey(conn.RemoteAddr().String())
    })

    return hasIt
}

func (s *Sessions) hasByKey(key string) (bool) {
    if _, ok := s.sessions[key]; !ok {
        return false
    }

    return true
}

func (s *Sessions) Register(connection *conn.Conn) (*Session, *types.Throw) {
    var er *types.Throw =   nil

    session             :=  &Session{
        conn:               connection,
        messager:           messager.NewMessager(s.defaultResponders),
        requestTimeout:     s.reqTimeout,
    }

    regErr              :=  session.registering()

    if regErr != nil {
        return nil, regErr
    }

    add                 :=  connection.RemoteAddr().String()

    s.sessionLock.Exec(func() {
        if s.hasByKey(add) {
            er          =   ErrSessionAlreadyRegistered.Throw(add)

            return
        }

        s.sessions[add] =   session

        s.onRegister()
    })

    if er != nil {
        return nil, er
    }

    return session, nil
}

func (s *Sessions) Unregister(connection *conn.Conn) (*types.Throw) {
    var er *types.Throw =   nil

    s.sessionLock.Exec(func() {
        addr            :=  connection.RemoteAddr().String()

        if !s.hasByKey(addr) {
            er          =   ErrSessionNotRegistered.Throw(addr)

            return
        }

        s.sessions[addr].unregistering()

        delete(s.sessions, addr)

        s.onUnregister()
    })

    return er
}

func (s *Sessions) Scan(excludedConns []*conn.Conn,
    callback func(string, *Session) (*types.Throw)) (*types.Throw) {
    var err *types.Throw = nil
    var excludes []string = []string{}

    for _, conn := range excludedConns {
        excludes        =   append(excludes, conn.RemoteAddr().String())
    }

    for sessKey, session := range s.sessions {
        skipThis        :=  false

        for _, excludedKey := range excludes {
            if excludedKey != sessKey {
                continue
            }

            skipThis    =   true
        }

        if skipThis {
            continue
        }

        err             =   callback(sessKey, session)

        if err != nil {
            return err
        }
    }

    return nil
}

func (s *Sessions) Clear() (*types.Throw) {
    var e *types.Throw  =   nil

    for _, sess := range s.sessions {
        delErr          :=  s.Unregister(sess.conn)

        if delErr == nil {
            continue
        }

        e               =   delErr
    }

    return e
}

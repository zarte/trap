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

package tcp

import (
    "github.com/raincious/trap/trap/core/types"
    "github.com/raincious/trap/trap/core/listen"
    "github.com/raincious/trap/trap/core/logger"

    "time"
    "net"
    "math/rand"
)

type TCP struct {
    responders      []Responder

    onError         func(listen.ConnectionInfo, *types.Throw)
    onPick          func(listen.ConnectionInfo, listen.RespondedResult)

    readTimeout     time.Duration
    writeTimeout    time.Duration
    totalTimeout    time.Duration

    inited          bool

    logger          *logger.Logger
    concurrent      uint

    rand            *rand.Rand
}

func (t *TCP) Init(c *listen.ProtocolConfig) (*types.Throw) {
    if t.inited {
        return listen.ErrProtocolAlreadyInited.Throw()
    }

    t.inited = true

    t.logger            = c.Logger.NewContext("TCP")

    t.onError           = c.OnError
    t.onPick            = c.OnPick

    t.readTimeout       = c.ReadTimeout
    t.writeTimeout      = c.WriteTimeout
    t.totalTimeout      = c.TotalTimeout
    t.concurrent        = uint(c.Concurrent.UInt16())

    t.rand              = rand.New(rand.NewSource(time.Now().UnixNano()))

    return nil
}

func (t *TCP) Responder(resp Responder) (listen.Protocol) {
    t.responders    = append(t.responders, resp)

    return t
}

func (t *TCP) getRandomResponder() (Responder, *types.Throw) {
    totalLen        := len(t.responders)

    if totalLen <= 0 {
        return nil, ErrNoAnyResponder.Throw()
    }

    randKey         := t.rand.Intn(totalLen)

    return t.responders[randKey], nil
}

func (t *TCP) Spawn(ip net.IP, port types.UInt16,
    setting types.String) (listen.Listener, *types.Throw) {
    listener := &Listener{}

    resp, rspErr := t.getRandomResponder()

    if rspErr != nil {
        t.logger.Warningf("Can't spawn the new TCP `Listener` due to error: %s",
            rspErr)

        return nil, rspErr
    }

    listener.Init(ListenerConfig{
        listen.ListenerConfig{
            Logger:         t.logger,
            Concurrent:     t.concurrent,

            OnError:        t.onError,
            OnPick:         t.onPick,

            ReadTimeout:    t.readTimeout,
            WriteTimeout:   t.writeTimeout,
            TotalTimeout:   t.totalTimeout,

            IP:             ip,
            Port:           port,
        },
        resp,
    })

    t.logger.Debugf("New TCP `Listener` has been spawned")

    return listener, nil
}
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

package udp

import (
    "github.com/raincious/trap/trap/core/types"
    "github.com/raincious/trap/trap/core/listen"
    "github.com/raincious/trap/trap/core/logger"

    "time"
    "net"
)

type UDP struct {
    onError         func(listen.ConnectionInfo, *types.Throw)
    onPick          func(listen.ConnectionInfo, listen.RespondedResult)

    readTimeout     time.Duration
    writeTimeout    time.Duration
    totalTimeout    time.Duration

    inited          bool

    logger          *logger.Logger
    concurrent      uint
}

func (t *UDP) Init(c *listen.ProtocolConfig) (*types.Throw) {
    if t.inited {
        return listen.ErrProtocolAlreadyInited.Throw()
    }

    t.inited = true

    t.logger            = c.Logger.NewContext("UDP")

    t.onError           = c.OnError
    t.onPick            = c.OnPick

    t.readTimeout       = c.ReadTimeout
    t.writeTimeout      = c.WriteTimeout
    t.totalTimeout      = c.TotalTimeout
    t.concurrent        = uint(c.Concurrent.UInt16())

    return nil
}

func (t *UDP) Spawn(ip net.IP, port types.UInt16,
    setting types.String) (listen.Listener, *types.Throw) {
    listener := &Listener{}

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
    })

    t.logger.Debugf("New UDP `Listener` has been spawned")

    return listener, nil
}
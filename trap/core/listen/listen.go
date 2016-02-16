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

package listen

import (
    "github.com/raincious/trap/trap/core/types"
    "github.com/raincious/trap/trap/core/logger"

    "time"
    "net"
)

type Port struct {
    Type            types.String
    IP              net.IP
    Port            types.UInt16
}

type Protocols map[types.String]Protocol

type Listen struct {
    inited          bool

    timeout         time.Duration
    logger          *logger.Logger

    listeners       []Listener
    randomPorts     types.UInt16

    error           *types.Throw

    protocols       Protocols

    onError         func(ConnectionInfo, *types.Throw)
    onPick          func(ConnectionInfo, RespondedResult)

    onListened      func(*ListeningInfo)
    onUnListened    func(*ListeningInfo)

    concurrent      types.UInt16
}

func (this *Listen) Init(cfg *Config) {
    if this.inited {
        return
    }

    this.inited         = true

    this.protocols      = Protocols{}

    this.timeout        = cfg.Timeout
    this.logger         = cfg.Logger.NewContext("Listen")

    this.onError        = cfg.OnError
    this.onPick         = cfg.OnPick

    this.onListened     = cfg.OnListened
    this.onUnListened   = cfg.OnUnListened

    this.concurrent     = cfg.Concurrent
}

func (this *Listen) GetLastError() (*types.Throw) {
    return this.error
}

func (this *Listen) Register(pType types.String, protocol Protocol) (*types.Throw) {
    if _, ok := this.protocols[pType]; ok {
        this.error = ErrProtocalRegistered.Throw(pType)

        this.logger.Warningf("Can't register another '%s' protocol", pType)

        return this.error
    }

    protocol.Init(&ProtocolConfig{
        OnError:        this.onError,
        OnPick:         this.onPick,

        ReadTimeout:    this.timeout,
        WriteTimeout:   this.timeout,
        TotalTimeout:   this.timeout,

        Logger:         this.logger,
        Concurrent:     this.concurrent,
    })

    this.protocols[pType] = protocol

    this.logger.Debugf("`Protocol` '%s' has been registered", pType)

    return nil
}

func (this *Listen) Add(pType types.String, ip net.IP, port types.UInt16,
    setting types.String) (*types.Throw) {
    if _, ok := this.protocols[pType]; !ok {
        this.error  =  ErrProtocalNotSupported.Throw(pType)

        this.logger.Warningf("Can't add `Listener`: %s", this.error)

        return this.error
    }

    listener, lErr  := this.protocols[pType].Spawn(ip, port, setting)

    if lErr != nil {
        this.logger.Warningf("Spawn `Listener` with error: %s", lErr)

        return lErr
    }

    this.listeners  =  append(this.listeners, listener)

    this.logger.Debugf("New `Listener` '%d' has been added at '%s:%d'",
        len(this.listeners) - 1, ip, port)

    return nil
}

func (this *Listen) Serv() (*types.Throw) {
    var lastErr *types.Throw = nil

    if this.error != nil {
        this.logger.Warningf("Can't serv due to previous error: %s", this.error)

        return this.error
    }

    for idx, listener := range this.listeners {
        upInfo, upErr := listener.Up()

        if upErr != nil {
            lastErr = upErr

            this.logger.Debugf("Can't bring up `Listener` '%d' due " +
                "to error: %s", idx, upErr)
        } else {
            this.onListened(upInfo)

            this.logger.Debugf("`Listener` '%d' is up", idx)
        }
    }

    if lastErr != nil {
        return lastErr
    }

    return nil
}

func (this *Listen) Down() (*types.Throw) {
    var lastErr *types.Throw = nil

    for idx, listener := range this.listeners {
        downInfo, downErr := listener.Down()

        if downErr != nil {
            lastErr = downErr

            this.logger.Debugf("Can't bring down `Listener` '%d' without " +
                "error: %s", idx, downErr)
        } else {
            this.onUnListened(downInfo)

            this.logger.Debugf("`Listener` '%d' is down", idx)
        }
    }

    if lastErr != nil {
        return lastErr
    }

    return nil
}
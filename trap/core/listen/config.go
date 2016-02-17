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

const (
    _                       = iota

    RESPOND_SUGGEST_SKIP
    RESPOND_SUGGEST_MARK
)

type Config struct {
    OnError         func(ConnectionInfo, *types.Throw)
    OnPick          func(ConnectionInfo, RespondedResult)

    OnListened      func(*ListeningInfo)
    OnUnListened    func(*ListeningInfo)

    Logger          *logger.Logger
    Concurrent      types.UInt16

    Timeout         time.Duration
}

type ListenerConfig struct {
    Logger          *logger.Logger
    Concurrent      uint

    OnError         func(ConnectionInfo, *types.Throw)
    OnPick          func(ConnectionInfo, RespondedResult)

    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    TotalTimeout    time.Duration

    IP              net.IP
    Port            types.UInt16
}

type ProtocolConfig struct {
    Logger          *logger.Logger
    Concurrent      types.UInt16

    OnError         func(ConnectionInfo, *types.Throw)
    OnPick          func(ConnectionInfo, RespondedResult)

    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    TotalTimeout    time.Duration
}

type ListeningInfo struct {
    Port            int
    IP              net.IP
    Protocol        string          // TCP or UDP in lowercase
}

type ConnectionInfo struct {
    ClientIP        types.IP
    ServerAddress   types.IPAddress
    Type            types.String    // Friendly name of the port which current
                                    //     connection attached like:
                                    //     http_proxy, tcp_filter, fake_ssh etc
}

type RespondedResult struct {
    ReceivedSample  [512]byte
    ReceivedLen     int

    RespondedData   [512]byte
    RespondedLen    int

    Suggestion      int
}
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

    "net"
    "time"
)

type Listener struct {
    inited          bool
    upped           bool
    closeable       bool

    listener        *net.UDPConn

    logger          *logger.Logger
    concurrent      int
    maxBytes        types.UInt32

    timeoutRead     time.Duration
    timeoutWrite    time.Duration
    timeoutTotal    time.Duration

    onError         func(listen.ConnectionInfo, *types.Throw)
    onPick          func(listen.ConnectionInfo, listen.RespondedResult)

    downChan        chan bool
    upwaitChan      chan bool

    listenOn        *net.UDPAddr
}

func (this *Listener) Init(cfg ListenerConfig) (*types.Throw) {
    if this.inited {
        return listen.ErrListenerAlreadyInited.Throw()
    }

    strIP                   :=  cfg.IP.String()

    this.inited             =   true
    this.upped              =   false
    this.closeable          =   false

    this.logger             =   cfg.Logger.NewContext("Listener").
                                    NewContext(types.String(cfg.IP.String())).
                                    NewContext(cfg.Port.String())

    this.concurrent         =   int(cfg.Concurrent.Int32()) // UInt16 to Int32

    if strIP == "0.0.0.0" {
        strIP = ""
    }

    this.timeoutRead        =   cfg.ReadTimeout
    this.timeoutWrite       =   cfg.WriteTimeout
    this.timeoutTotal       =   cfg.TotalTimeout

    this.maxBytes           =   cfg.MaxBytes

    this.onError            =   cfg.OnError
    this.onPick             =   cfg.OnPick

    this.downChan           =   make(chan bool, 1)
    this.upwaitChan         =   make(chan bool)

    this.listenOn           =   &net.UDPAddr{
                                    IP: cfg.IP,
                                    Port: int(cfg.Port.Int16()),
                                }

    return nil
}

func (this *Listener) Up() (*listen.ListeningInfo, *types.Throw) {
    if this.upped {
        return nil, listen.ErrListenerAlreadyUp.Throw(this.listenOn)
    }

    udpConn, udpConErr := net.ListenUDP("udp", this.listenOn)

    if udpConErr != nil {
        return nil, types.ConvertError(udpConErr)
    }

    this.upped      = true
    this.listener   = udpConn

    go func() {
        conChan := make(chan bool, this.concurrent)

        defer func() {
            if !this.upped {
                return
            }

            this.listener.Close()
        }()

        this.closeable = true
        this.upwaitChan <- true

        go func() {
            for {
                time.Sleep(1 * time.Second)

                select {
                    case <- this.downChan:
                        return

                    case <- conChan:
                        curChanLen := len(conChan)

                        if curChanLen < this.concurrent {
                            for i := curChanLen; i < this.concurrent; i++ {
                                conChan <- true
                            }
                        }

                    default:
                        for i := 0; i < this.concurrent; i++ {
                            conChan <- true
                        }
                }
            }
        }()

        totalbuffer := make([]byte, this.maxBytes)

        this.logger.Debugf("Waiting for connection. Maximum rate is " +
            "'%d' bytes per second",
            len(totalbuffer) * this.concurrent)

        for {
            select {
                case <- this.downChan:
                    return

                case <- conChan:
                    this.listener.SetDeadline(time.Now().Add(this.timeoutTotal))
                    this.listener.SetReadDeadline(time.Now().Add(this.timeoutRead))
                    this.listener.SetWriteDeadline(time.Now().Add(this.timeoutWrite))

                    length, srcAddr, conErr :=  this.listener.ReadFromUDP(totalbuffer)

                    if conErr != nil {
                        continue
                    }

                    clientAddr, cAddrErr    :=  types.ConvertIPAddress(srcAddr)
                    serverAddr, sAddrErr    :=  types.ConvertIPAddress(this.listenOn)

                    if cAddrErr != nil || sAddrErr != nil {
                        continue
                    }

                    connection              :=  listen.ConnectionInfo{
                        ClientIP:               clientAddr.IP,
                        ServerAddress:          serverAddr,
                        Type:                   "udp",
                    }

                    result                  :=  listen.RespondedResult{
                        Suggestion:             listen.RESPOND_SUGGEST_MARK,
                    }

                    result.ReceivedSample   =   totalbuffer[:length]

                    this.onPick(connection, result)
            }
        }
    }()

    <- this.upwaitChan

    return &listen.ListeningInfo{
        Port: this.listenOn.Port,
        IP: this.listenOn.IP,
        Protocol: "udp",
    }, nil
}

func (this *Listener) Down() (*listen.ListeningInfo, *types.Throw) {
    if !this.closeable {
        return nil, listen.ErrListenerNotCloseable.Throw(this.listenOn)
    }

    this.upped      =   false
    this.closeable  =   false

    this.downChan   <-  true
    this.downChan   <-  true

    lnErr := this.listener.Close()

    if lnErr != nil {
        return nil, types.ConvertError(lnErr)
    }

    return &listen.ListeningInfo{
        Port: this.listenOn.Port,
        IP: this.listenOn.IP,
        Protocol: "udp",
    }, nil
}
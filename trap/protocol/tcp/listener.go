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

    "net"
    "time"
    "sync"
)

type Listener struct {
    inited          bool
    upped           bool
    closeable       bool

    listener        *net.TCPListener
    responder       Responder

    logger          *logger.Logger
    concurrent      uint

    timeoutRead     time.Duration
    timeoutWrite    time.Duration
    timeoutTotal    time.Duration

    onError         func(listen.ConnectionInfo, *types.Throw)
    onPick          func(listen.ConnectionInfo, listen.RespondedResult)

    listenOn        *net.TCPAddr

    downChan        chan bool
    upwaitChan      chan bool

    waitingGroup    sync.WaitGroup
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

    this.concurrent         =   cfg.Concurrent

    if strIP == "0.0.0.0" {
        strIP = ""
    }

    this.timeoutRead        =   cfg.ReadTimeout
    this.timeoutWrite       =   cfg.WriteTimeout
    this.timeoutTotal       =   cfg.TotalTimeout

    this.onError            =   cfg.OnError
    this.onPick             =   cfg.OnPick

    this.responder          =   cfg.Responder

    this.listenOn           =   &net.TCPAddr{
                                    IP:     cfg.IP,
                                    Port:   int(cfg.Port.Int16()),
                                }

    this.waitingGroup       =   sync.WaitGroup{}
    this.downChan           =   make(chan bool, 2)
    this.upwaitChan         =   make(chan bool)

    return nil
}

func (this *Listener) Up() (*listen.ListeningInfo, *types.Throw) {
    if this.upped {
        return nil, listen.ErrListenerAlreadyUp.Throw(this.listenOn)
    }

    listener, lErr          :=  net.ListenTCP("tcp", this.listenOn)

    if lErr != nil {
        return nil, types.ConvertError(lErr)
    }

    this.listener           =  listener
    this.upped              =  true

    // Main loop waiter
    this.waitingGroup.Add(1)

    go func() {
        defer this.waitingGroup.Done()

        conChan             :=  make(chan bool, this.concurrent)

        defer func() {
            if !this.upped {
                return
            }

            this.listener.Close()

            this.logger.Debugf("Defered connection close executed")
        }()

        this.waitingGroup.Add(1) // con chan waiter

        go func() {
            defer this.waitingGroup.Done()

            for {
                time.Sleep(1 * time.Second)

                select {
                    case <- this.downChan:
                        return

                    case <- conChan:
                        curChanLen := uint(len(conChan))

                        if curChanLen < this.concurrent {
                            for i := curChanLen; i < this.concurrent; i++ {
                                conChan <- true
                            }
                        }

                    default:
                        for i := uint(0); i < this.concurrent; i++ {
                            conChan <- true
                        }
                }
            }
        }()

        this.closeable = true

        this.logger.Debugf("Waiting for connection. Maximum concurrent is '%d'",
            this.concurrent)

        this.upwaitChan <- true

        // Listen loop
        for {
            select {
                case <- this.downChan:
                    return

                case <- conChan:
                    conn, err       := this.listener.AcceptTCP()

                    if err != nil {
                        continue
                    }

                    this.waitingGroup.Add(1)

                    // Dispatch a routine to serve this com
                    go func(conn *net.TCPConn) {
                        defer this.waitingGroup.Done()


                        clientAddr, cAddrErr    :=  types.ConvertIPAddress(
                                                        conn.RemoteAddr())

                        serverAddr, sAddrErr    :=  types.ConvertIPAddress(
                                                        conn.LocalAddr())

                        if cAddrErr != nil || sAddrErr != nil{
                            return
                        }

                        connection  :=  listen.ConnectionInfo{
                            ClientIP:               clientAddr.IP,
                            ServerAddress:          serverAddr,
                            Type:                   "tcp",
                        }

                        conn.SetDeadline(time.Now().Add(this.timeoutTotal))
                        conn.SetReadDeadline(time.Now().Add(this.timeoutRead))
                        conn.SetWriteDeadline(time.Now().Add(this.timeoutWrite))

                        result, err := this.responder.Handle(conn)

                        if err != nil {
                            this.onError(connection, err)
                        }

                        closeErr    := conn.Close()

                        if closeErr != nil {
                            this.onError(connection,
                                types.ConvertError(closeErr))
                        }

                        this.onPick(connection, result)
                    }(conn)
            }
        }
    }()

    // wait for listen
    <-this.upwaitChan

    return &listen.ListeningInfo{
        Port: this.listenOn.Port,
        IP: this.listenOn.IP,
        Protocol: "tcp",
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

    closeErr        :=  this.listener.Close()

    this.waitingGroup.Wait()

    if closeErr != nil {
        return nil, types.ConvertError(closeErr)
    }

    return &listen.ListeningInfo{
        Port: this.listenOn.Port,
        IP: this.listenOn.IP,
        Protocol: "tcp",
    }, nil
}
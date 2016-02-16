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

package trap

import (
    "github.com/raincious/trap/trap/core/types"
    "github.com/raincious/trap/trap/core/logger"
    "github.com/raincious/trap/trap/core/client"
    "github.com/raincious/trap/trap/core/server"
    "github.com/raincious/trap/trap/core/status"
    "github.com/raincious/trap/trap/core/status/controller"

    "net"
    "net/http"
    "sync"
    "time"
)

type Status struct {
    ip                  types.IP
    host                types.String
    port                types.UInt16

    accounts            status.Accounts

    logger              *logger.Logger

    status              *http.Server
    server              *Server
    serverRWLock        types.Mutex

    sessions            status.Sessions

    sessionRWLock       types.Mutex

    statusListener      *net.TCPListener
    statusDownWait      sync.WaitGroup
}

func NewStatus() (*Status) {
    ip, ipErr           :=  types.ConvertIPFromString("127.0.0.1")

    if ipErr != nil {
        ip              =   types.IP{}
    }

    return &Status{
        ip:             ip,
        host:           "127.0.0.1",
        port:           1793,
        accounts:       status.Accounts{},
        sessions:       status.Sessions{},
        serverRWLock:   types.Mutex{},
        sessionRWLock:  types.Mutex{},
        statusDownWait: sync.WaitGroup{},
    }
}

func (this *Status) SetLogger(l *logger.Logger) {
    this.logger         =   l.NewContext("Status")
}

func (this *Status) SetServer(s *Server) {
    this.server         =   s
}

func (this *Status) IP(ip types.IP) {
    this.ip             =   ip
}

func (this *Status) Host(host types.String) {
    this.host           =   host
}

func (this *Status) Port(port types.UInt16) {
    this.port           =   port
}

func (this *Status) Account(pass types.String,
    permissions []types.String) (*status.Account, *types.Throw) {
    return this.accounts.Register(pass, permissions)
}

func (this *Status) verifyUser(ip net.IP,
    sessionKey types.String) (*status.Session, *types.Throw) {
    var sess    *status.Session =   nil
    var err     *types.Throw    =   nil

    this.sessionRWLock.Exec(func() {
        sess, err               =   this.sessions.Verify(ip, sessionKey)
    })

    return sess, err
}

func (this *Status) authUser(ip net.IP,
    pass types.String) (*status.Session, *types.Throw) {
    var result                  *status.Session = nil
    var resultErr               *types.Throw = nil

    account, accountErr         :=  this.accounts.Get(pass)

    if accountErr != nil {
        this.server.AddClient(types.ConvertIP(ip), server.ClientConInfo{
            Server:             types.IPAddress{
                                    IP:             types.IP{},
                                    Port:           this.port,
                                },
            Type:               "status_ui",
            Marked:             true,
        })

        this.logger.Warningf("Bad authorization attempt from '%s'", ip)

        return nil, accountErr
    }

    this.logger.Infof("A new session has been binded with '%s'", ip)

    this.sessionRWLock.Exec(func() {
        result, resultErr       =   this.sessions.Add(ip, pass,
                                        account, 12 * time.Hour)
    })

    return result, resultErr
}

func (this *Status) getAllSessions() ([]status.SessionDump) {
    var dump []status.SessionDump

    this.sessionRWLock.Exec(func() {
        dump                    =   this.sessions.Dump()
    })

    return dump
}

func (this *Status) getNewServer(httpHost types.String) (*http.Server, *types.Throw) {
    httpMux         :=  status.NewMux()

    httpMux.HandleController("/", &controller.Home{
        StaticPage:         status.StaticClientPage,
    })

    httpMux.HandleController("/api/auth", &controller.Auth{
        Verify:             this.verifyUser,
        Auth:               this.authUser,
    })

    httpMux.HandleController("/api/status", &controller.Status{
        controller.SessionedJSON{
            Verify:         this.verifyUser,
        },
        func() (server.Status) {
            return this.server.Status()
        },
    })

    httpMux.HandleController("/api/clients", &controller.Clients{
        controller.SessionedJSON{
            Verify:         this.verifyUser,
        },
        func() ([]client.Client) {
            return this.server.Clients()
        },
    })
    httpMux.HandleController("/api/client", &controller.Client{
        controller.SessionedJSON{
            Verify:         this.verifyUser,
        },
        func(addr types.IP) (*client.Client, *types.Throw) {
            return this.server.Client(addr)
        },
        func(addr types.IP, cCon server.ClientConInfo) (*client.Client, *types.Throw) {
            return this.server.AddClient(addr, cCon)
        },
        func(addr types.IP) (*types.Throw) {
            return this.server.RemoveClient(addr)
        },
    })

    httpMux.HandleController("/api/logs", &controller.Logs{
        controller.SessionedJSON{
            Verify:         this.verifyUser,
        },
        func() ([]logger.LogExport) {
            return this.logger.Dump()
        },
    })

    httpMux.HandleController("/api/sessions", &controller.Sessions{
        controller.SessionedJSON{
            Verify:         this.verifyUser,
        },
        func() ([]status.SessionDump) {
            return this.getAllSessions()
        },
    })

    return &http.Server{
        Addr:               httpHost.String(),
        Handler:            httpMux,
        WriteTimeout:       2 * time.Second,
        ReadTimeout:        2 * time.Second,
    }, nil
}

func (this *Status) up() (*types.Throw) {
    listenOn                :=  &net.TCPAddr{
        IP:                 this.ip.IP(),
        Port:               int(this.port.Int32()),
    }

    // Check if the server is down before start a new one
    if this.status != nil || this.statusListener != nil {
        return status.ErrServerAlreadyUp.Throw()
    }

    // No server? no up
    if this.server == nil {
        return status.ErrServerNotSet.Throw()
    }

    // Create an instance of the http status server
    sServer, sErr       :=  this.getNewServer(this.host)

    if sErr != nil {
        return types.ConvertError(sErr)
    }

    this.status         =   sServer

    listener, lErr      :=  net.ListenTCP("tcp", listenOn)

    if lErr != nil {
        return types.ConvertError(lErr)
    }

    go func() {
        this.statusDownWait.Add(1)

        defer this.statusDownWait.Done()

        this.logger.Infof("Serving `Status` server at: %s", listenOn)

        this.statusListener = listener

        servErr         :=  this.status.Serve(this.statusListener)

        if servErr != nil {
            this.logger.Infof("`Status` server closed due to error: %s",
                servErr)

            return
        }
    }()

    return nil
}

func (this *Status) down() (*types.Throw) {
    if this.status == nil {
        return status.ErrServerNotDownable.Throw()
    }

    if this.statusListener == nil {
        return status.ErrServerNotDownable.Throw()
    }

    this.statusListener.Close()
    this.statusDownWait.Wait()

    this.status         =   nil
    this.statusListener =   nil

    return nil
}

func (this *Status) Reset() (*types.Throw) {
    var e *types.Throw  =   nil

    this.serverRWLock.Exec(func() {
        e               =   this.down()

        if e != nil {
            return
        }

        this.accounts   =   status.Accounts{}
        this.sessions   =   status.Sessions{}

        // Needs manual up
    })

    return e
}

func (this *Status) Serv() (*types.Throw) {
    var e *types.Throw  =   nil

    this.serverRWLock.Exec(func() {
        e               =   this.up()
    })

    return e
}

func (this *Status) Down() (*types.Throw) {
    var e *types.Throw  =   nil

    this.serverRWLock.Exec(func() {
        e               =   this.down()
    })

    return e
}
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

package main

import (
    "github.com/raincious/trap/trap"
    "github.com/raincious/trap/trap/core"

    "github.com/raincious/trap/trap/core/event"
    "github.com/raincious/trap/trap/core/types"
    "github.com/raincious/trap/trap/core/logger"

    "github.com/raincious/trap/trap/protocol/tcp"
    tcpResponder "github.com/raincious/trap/trap/protocol/tcp/responder"

    "github.com/raincious/trap/trap/protocol/udp"

    logPrinter "github.com/raincious/trap/trap/logger"
    "github.com/raincious/trap/trap/config"

    "os/signal"
    "os/exec"
    "os"
    "flag"
    "syscall"
    "time"
    "net"
    "fmt"
    "bufio"
    "path/filepath"
)

var (
    logFile     =   ""
    cfgFile     =   ""
)

func init() {
    flag.StringVar(&logFile, "log", "", "Path of log file")
    flag.StringVar(&cfgFile, "config", "", "Config file")

    flag.Parse()
}

func initConfig(server *trap.Server, status *trap.Status) {
    if cfgFile == "" {
        panic(fmt.Errorf("Configuration is not specified. " +
            "Please use command `%s -help` for more information",
            filepath.Base(os.Args[0])))
    }

    cfg, err    := config.Load(cfgFile)

    if err != nil {
        panic(fmt.Errorf("Can't load configuration file '%s' under error: %s",
            cfgFile, err))
    }

    server.SetTimeout(1 * time.Second)
    server.SetTolerate(1, 3600 * time.Second)
    server.SetConcurrentLimit(100)

    // Init TCP Protocol
    tcpProtocol := &tcp.TCP{}

    tcpProtocol.Responder(&tcpResponder.Echo{})
    tcpProtocol.Responder(&tcpResponder.Empty{})

    server.Listen().Register("tcp", tcpProtocol)

    // Init UDP Protocol
    udpProtocol := &udp.UDP{}

    server.Listen().Register("udp", udpProtocol)

    // Register ports
    for _, listenPort := range cfg.Listens {
        sAddErr := server.Listen().Add(listenPort.Type,
            net.ParseIP(listenPort.IP.String()), listenPort.Port,
            listenPort.Additional)

        if sAddErr == nil {
            continue
        }

        panic(fmt.Errorf("Error registering '%s' listener '%s,%d': %s",
            listenPort.Type, listenPort.IP, listenPort.Port, sAddErr))
    }

    // Register events
    for eventName, eventCommands := range cfg.Commands {
        for _, eventCommand := range eventCommands {
            func(eName types.String, eCmd config.Command) {
                server.Event().Register(eName,
                    func(p *event.Parameters) (*types.Throw) {
                        var params []string

                        for _, cmdParam := range eCmd.Parameters {
                            params = append(params, p.Parse(cmdParam.Format,
                                cmdParam.Labels).String())
                        }

                        cmd := exec.Command(eCmd.Command.String(),
                            params...)

                        err := cmd.Run()

                        if err != nil {
                            return types.ConvertError(err)
                        }

                        return nil
                    })
            }(eventName, eventCommand)
        }
    }

    // Start `Status` Server to display some of the status of the server
    if cfg.StatusPort > 0 {
        status.SetServer(server)

        status.Port(cfg.StatusPort)

        if !cfg.StatusInterface.IsEmpty() {
            status.IP(cfg.StatusInterface)
        }

        if cfg.StatusHost != "" {
            status.Host(cfg.StatusHost)
        }

        for account, permissions := range cfg.StatusAccounts {
            status.Account(account, permissions)
        }

        server.OnUpDown(func() (*types.Throw) {
            return status.Serv()
        }, func() (*types.Throw) {
            return status.Down()
        })
    }
}

func main() {
    fmt.Printf(core.TRAP_COMMAND_BANNDER, core.TRAP_DESCRIPTION,
        core.TRAP_COPYRIGHT, core.TRAP_VERSION, core.TRAP_LICENSE,
        core.TRAP_LICENSEURL, core.TRAP_PROJECTURL, core.TRAP_SOURCEURL)

    log         :=  logger.NewLogger()

    // Register file logger if variable `logFile` is filled
    if logFile != "" {
        logFileHander, logFErr  :=  os.OpenFile(logFile,
            os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY,
            0600)

        if logFErr != nil {
            panic(fmt.Errorf("Can't open log file '%s' due to error: %s",
                logFile, logFErr))
        }

        loggerBuffer            :=  bufio.NewWriter(logFileHander)

        defer func() {
            loggerBuffer.Flush()

            logFileHander.Close()
        }()

        fileLogPrinter, fLPErr  :=  logPrinter.NewFilePrinter(loggerBuffer)

        if fLPErr != nil {
            panic(fmt.Errorf("Can't create File logger due to error: %s",
                fLPErr))
        }

        log.Register(fileLogPrinter)
    } else {
        // Or, register screen logger instead
        log.Register(logPrinter.NewScreenPrinter())
    }

    signalCall  :=  make(chan os.Signal, 1)

    defer close(signalCall)

    server      :=  trap.NewServer()

    defer server.Down()

    server.SetLogger(log)

    // Init Status server
    status      :=  trap.NewStatus()

    status.SetLogger(log)

    initConfig(server, status)

    servErr     :=  server.Serv()

    if servErr != nil {
        server.Down()

        panic(servErr)
    }

    // Register system signal handlers
    signal.Notify(signalCall,
        syscall.SIGHUP,     // For Reload
        syscall.SIGINT)    // Unstopable shutdown

    // Loop for system signal
    for {
        signalExit := false

        switch <-signalCall {
            case syscall.SIGHUP:
                server.Reload(func(s *trap.Server) (*types.Throw) {
                    status.Reset()

                    initConfig(s, status)

                    return nil
                })

            case syscall.SIGINT:
                log.Infof("Exit signal picked up")

                signalExit = true

            default:
                log.Infof("Unknown signal")
        }

        if signalExit {
            break
        }
    }
}
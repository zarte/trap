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
    statusPkg "github.com/raincious/trap/trap/core/status"
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
    "log"
    "fmt"
    "bufio"
    "path/filepath"
    "runtime/pprof"
)

var (
    logFile         =   ""
    silentRun       =   false
    cfgFile         =   ""
    cpuPrfFile      =   ""
    memPrfFile      =   ""
)

func init() {
    flag.StringVar(&cfgFile, "config", "",
        "Load configuration from specified file, must be defined.")

    flag.StringVar(&logFile, "log", "",
        "Save log data to specified file, " +
        "keep it default to disable file logger.")

    flag.StringVar(&cpuPrfFile, "profiling-cpu", "",
        "Dump CPU profile data to specified file, " +
        "keep it blank to disable profiling.")

    flag.StringVar(&memPrfFile, "profiling-mem", "",
        "Dump memory profile data to specified file, " +
        "keep it blank to disable profiling.")

    flag.BoolVar(&silentRun, "silent", false,
        "Do not generate any output")

    flag.Parse()
}

func initConfig(server *trap.Server, status *trap.Status) {
    if cfgFile == "" {
        panic(fmt.Errorf("Configuration is not specified. " +
            "Please use command `%s -help` for more information",
            filepath.Base(os.Args[0])))
    }

    cfg, err    :=  config.Load(cfgFile)

    if err != nil {
        panic(fmt.Errorf("Can't load configuration file '%s' under error: %s",
            cfgFile, err))
    }

    if cfg.AttemptTimeout > 0 {
        server.SetTimeout(time.Duration(cfg.AttemptTimeout) * time.Second)
    }

    if cfg.AttemptMaxBytes > 0 {
        server.SetClientRecordDataLimit(cfg.AttemptMaxBytes)
    }

    if cfg.AttemptThershold > 0 && cfg.AttemptExpire > 0 {
        server.SetTolerate(cfg.AttemptThershold,
            time.Duration(cfg.AttemptExpire.Int64()) * time.Second,
            time.Duration(cfg.AttemptRestrict.Int64()) * time.Second)
    }

    server.SetConcurrentLimit(100)

    // Init TCP Protocol
    tcpProtocol :=  &tcp.TCP{}

    tcpProtocol.Responder(&tcpResponder.Echo{})
    tcpProtocol.Responder(&tcpResponder.Empty{})

    server.Listen().Register("tcp", tcpProtocol)

    // Init UDP Protocol
    server.Listen().Register("udp", &udp.UDP{})

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

        if cfg.StatusTLSCert != "" && cfg.StatusTLSCertKey != "" {
            status.LoadCert(cfg.StatusTLSCert, cfg.StatusTLSCertKey)
        }

        for account, permissions := range cfg.StatusAccounts {
            _, sAccErr := status.Account(account, permissions)

            if sAccErr == nil {
                continue
            }

            panic(fmt.Errorf("Error registering status account '%s' due to error: %s",
                account, sAccErr))
        }

        server.OnUpDown(func() (*types.Throw) {
            return status.Serv()
        }, func() (*types.Throw) {
            return status.Down()
        })
    }
}

func main() {
    // Enable CPU profiling
    if cpuPrfFile != "" {
        cpuPrfFileFile, cpuPrfFileErr :=  os.OpenFile(cpuPrfFile,
            os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY,
            0600)

        if cpuPrfFileErr != nil {
            panic(fmt.Errorf("Can't create CPU profiling file due to error: %s",
                cpuPrfFileErr))
        }

        pprof.StartCPUProfile(cpuPrfFileFile)

        defer func() {
            pprof.StopCPUProfile()

            cpuPrfFileFile.Close()
        }()
    }

    // Enable Memory profiling
    if memPrfFile != "" {
        memPrfFileFile, memFileErr :=  os.OpenFile(memPrfFile,
            os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY,
            0600)

        if memFileErr != nil {
            panic(fmt.Errorf("Can't create Mem profiling file due to error: %s",
                memFileErr))
        }

        defer func() {
            pprof.WriteHeapProfile(memPrfFileFile)

            memPrfFileFile.Close()
        }()
    }

    if !silentRun {
        fmt.Printf(core.TRAP_COMMAND_BANNDER, core.TRAP_DESCRIPTION,
            core.TRAP_COPYRIGHT, core.TRAP_VERSION, core.TRAP_LICENSE,
            core.TRAP_LICENSEURL, core.TRAP_PROJECTURL, core.TRAP_SOURCEURL)
    }

    logging         :=  logger.NewLogger()

    log.SetOutput(logging.NewContext("System"))

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

        logging.Register(fileLogPrinter)
    } else if (!silentRun) {
        // Or, register screen logger instead
        logging.Register(logPrinter.NewScreenPrinter())
    }

    // Start booting
    server      :=  trap.NewServer()

    defer server.Down()

    server.SetLogger(logging)

    // Init Status server
    status      :=  trap.NewStatus()

    status.SetLogger(logging)

    initConfig(server, status)

    servErr     :=  server.Serv()

    if servErr != nil {
        server.Down()

        panic(fmt.Errorf("Encountered at least one error while " +
            "server is booting up: %s", servErr))
    }

    // Catch system signals
    signalCall  :=  make(chan os.Signal, 1)

    defer close(signalCall)

    // Register system signal handlers
    signal.Notify(signalCall,
        syscall.SIGHUP,    // For Reload
        syscall.SIGINT,
        syscall.SIGTERM)    // Control + C

    // Loop for system signal
    for {
        signalExit := false

        callSignal := <-signalCall

        switch {
            case callSignal == syscall.SIGHUP:
                server.Reload(func(s *trap.Server) (*types.Throw) {
                    statusErr := status.Reset()

                    if statusErr != nil &&
                    !statusErr.Is(statusPkg.ErrServerNotDownable) {
                        return statusErr
                    }

                    initConfig(s, status)

                    return nil
                })

            case callSignal == syscall.SIGINT || callSignal == syscall.SIGTERM:
                logging.Infof("Exit signal picked up")

                signalExit = true

            default:
                logging.Infof("Unknown signal")
        }

        if signalExit {
            break
        }
    }
}
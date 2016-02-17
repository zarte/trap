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

package logger

import (
    "github.com/raincious/trap/trap/core/types"

    "fmt"
    "time"
)

type Logger struct {
    printers            *Printers
    logs                *Logs
    mutex               *types.Mutex

    context             types.String
}

func NewLogger() (*Logger) {
    return &Logger{
        printers:       &Printers{},
        logs:           &Logs{},
        mutex:          &types.Mutex{},
        context:        "Trap",
    }
}

func (l *Logger) Register(printer Printer) {
    l.printers.Add(printer)
}

func (l *Logger) NewContext(s types.String) (*Logger) {
    return &Logger{
        printers:       l.printers,
        logs:           l.logs,
        mutex:          l.mutex,
        context:        l.context + ":" + s,
    }
}

func (l *Logger) append(log Log) {
    l.mutex.Exec(func() {
        l.logs.Append(log, 256)

        switch log.Type {
            case LOG_TYPE_DEBUG:
                l.printers.Debug(log.Context, log.Time, log.Message)

            case LOG_TYPE_INFO:
                l.printers.Info(log.Context, log.Time, log.Message)

            case LOG_TYPE_WARNING:
                l.printers.Warning(log.Context, log.Time, log.Message)

            case LOG_TYPE_ERROR:
                l.printers.Error(log.Context, log.Time, log.Message)

            default:
                l.printers.Print(log.Context, log.Time, log.Message)
        }
    })
}

func (l *Logger) Debugf(s string, v ...interface{}) {
    /* Disable debug totally after finish develpment
    l.append(Log{
        Time:           time.Now(),
        Type:           LOG_TYPE_DEBUG,
        Context:        l.context,
        Message:        types.String(fmt.Sprintf(s, v...)),
    })*/
}

func (l *Logger) Infof(s string, v ...interface{}) {
    l.append(Log{
        Time:           time.Now(),
        Type:           LOG_TYPE_INFO,
        Context:        l.context,
        Message:        types.String(fmt.Sprintf(s, v...)),
    })
}

func (l *Logger) Warningf(s string, v ...interface{}) {
    l.append(Log{
        Time:           time.Now(),
        Type:           LOG_TYPE_WARNING,
        Context:        l.context,
        Message:        types.String(fmt.Sprintf(s, v...)),
    })
}

func (l *Logger) Errorf(s string, v ...interface{}) {
    l.append(Log{
        Time:           time.Now(),
        Type:           LOG_TYPE_ERROR,
        Context:        l.context,
        Message:        types.String(fmt.Sprintf(s, v...)),
    })
}

func (l *Logger) Dump() ([]LogExport) {
    logs                :=  []LogExport{}

    l.mutex.Exec(func() {
        logs            =   l.logs.Export()
    })

    return logs
}
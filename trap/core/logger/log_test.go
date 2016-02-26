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
    "testing"
    "time"
)

func TestLogsAppend(t *testing.T) {
    logs := Logs{}

    defaultLogItem := Log{
        Time:       time.Now(),
        Type:       LOG_TYPE_DEFAULT,
        Context:    "Test default context",
        Message:    "Test default message",
    }

    if len(logs) != 0 {
        t.Error("Unexpected initial data")

        return
    }

    logs.Append(defaultLogItem, 2)

    if len(logs) != 1 {
        t.Error("Unexpected amount of log items, expecting '1', got '%d'",
            len(logs))

        return
    }

    logs.Append(defaultLogItem, 2)

    if len(logs) != 2 {
        t.Error("Logs.Append() failed to append log items")

        return
    }

    logs.Append(defaultLogItem, 2)
    logs.Append(defaultLogItem, 2)

    if len(logs) != 2 {
        t.Error("Logs.Append() failed to remove old log items")

        return
    }

    logs.Append(defaultLogItem, 1)

    if len(logs) != 1 {
        t.Error("Logs.Append() failed to remove old log items")

        return
    }

    // Log can be all removed by this way
    logs.Append(defaultLogItem, 0)

    if len(logs) != 0 {
        t.Error("Logs.Append() failed to remove old log items")

        return
    }
}

func TestLogsExport(t *testing.T) {
    logs        :=  Logs{}
    logItems    :=  []Log{}

    logItems    =   append(logItems, Log{
        Time:       time.Now(),
        Type:       LOG_TYPE_DEFAULT,
        Context:    "Test default context",
        Message:    "Test default message",
    })

    logItems    =   append(logItems, Log{
        Time:       time.Now().Add(100 * time.Millisecond),
        Type:       LOG_TYPE_DEBUG,
        Context:    "Test debug context",
        Message:    "Test debug message",
    })

    logItems    =   append(logItems, Log{
        Time:       time.Now().Add(200 * time.Millisecond),
        Type:       LOG_TYPE_INFO,
        Context:    "Test info context",
        Message:    "Test info message",
    })

    logItems    =   append(logItems, Log{
        Time:       time.Now().Add(300 * time.Millisecond),
        Type:       LOG_TYPE_WARNING,
        Context:    "Test warning context",
        Message:    "Test warning message",
    })

    logItems    =   append(logItems, Log{
        Time:       time.Now().Add(400 * time.Millisecond),
        Type:       LOG_TYPE_ERROR,
        Context:    "Test error context",
        Message:    "Test error message",
    })

    logs.Append(logItems[4], 64)
    logs.Append(logItems[3], 64)
    logs.Append(logItems[2], 64)
    logs.Append(logItems[1], 64)
    logs.Append(logItems[0], 64)

    for idx, exportedLog := range logs.Export() {
        if exportedLog.Time != logItems[idx].Time ||
            exportedLog.Context != logItems[idx].Context ||
            exportedLog.Message != logItems[idx].Message {
            t.Errorf("Logs.Export() exports invalid log data")

            return
        }

        switch logItems[idx].Type {
            case LOG_TYPE_DEBUG:
                if exportedLog.Type != "Debug" {
                    t.Errorf("Logs.Export() exports invalid log data")

                    return
                }

            case LOG_TYPE_INFO:
                if exportedLog.Type != "Information" {
                    t.Errorf("Logs.Export() exports invalid log data")

                    return
                }

            case LOG_TYPE_WARNING:
                if exportedLog.Type != "Warning" {
                    t.Errorf("Logs.Export() exports invalid log data")

                    return
                }

            case LOG_TYPE_ERROR:
                if exportedLog.Type != "Error" {
                    t.Errorf("Logs.Export() exports invalid log data")

                    return
                }
        }
    }
}
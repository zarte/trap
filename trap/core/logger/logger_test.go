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

    "testing"
    "time"
    "sync"
)

type fakeLogPrinterForMutex struct {
    order           []types.String
}

// Make sure the timeline is sorted:
// 600 = 300 + 200 + 100
// 300 = 200 + 100
// 200 = 100 * 2
func (f *fakeLogPrinterForMutex) Info(c types.String, t time.Time, m types.String) {
    time.Sleep(600 * time.Millisecond)

    f.order = append(f.order, "Info")
}

func (f *fakeLogPrinterForMutex) Debug(c types.String, t time.Time, m types.String) {
    // Normally not call
    time.Sleep(300 * time.Millisecond)

    f.order = append(f.order, "Debug")
}

func (f *fakeLogPrinterForMutex) Warning(c types.String, t time.Time, m types.String) {
    time.Sleep(200 * time.Millisecond)

    f.order = append(f.order, "Warning")
}

func (f *fakeLogPrinterForMutex) Error(c types.String, t time.Time, m types.String) {
    time.Sleep(100 * time.Millisecond)

    f.order = append(f.order, "Error")
}

func (f *fakeLogPrinterForMutex) Print(c types.String, t time.Time, m types.String) {
    f.order = append(f.order, "Print")
}

func TestLoggerNewContext(t *testing.T) {
    logger := NewLogger()

    if logger.context != "Trap" {
        t.Error("The context name of root Logger must be 'Trap'")

        return
    }

    newCtx := logger.NewContext("Sub context").NewContext("Subagain")

    if newCtx.context != "Trap:Sub context:Subagain" {
        t.Error("Unexpected context title for the new context")

        return
    }

    if newCtx.logs != logger.logs ||
        newCtx.printers != logger.printers ||
        newCtx.mutex != logger.mutex {
        t.Error("The new context is not using inherited properties")

        return
    }
}

func testLoggerAppend(t *testing.T, logger *Logger, baseN int) {
    if len(logger.Dump()) != (baseN - 1) * 4{
        t.Error("Unexpected initial amount of log items")

        return
    }

    logger.Debugf("Formated %s", "String")
    logger.Infof("Formated integer %d", 10)
    logger.Warningf("Formated integer %d", 10)
    logger.Errorf("Formated integer %d", 10)
    logger.Write([]byte("This is a slice of bytes"))

    dumpped := logger.Dump()

    // Notice the `Debugf` should be commented out when not in development
    // So it wouldn't be count
    if len(dumpped) != 4 * baseN {
        t.Error("Unexpected amount of log items")

        return
    }
}

func TestLoggerAppend(t *testing.T) {
    logger := NewLogger()

    testLoggerAppend(t, logger, 1)
    testLoggerAppend(t, logger.NewContext("Another context"), 2)
    testLoggerAppend(t, logger.NewContext("Another context").NewContext("Sub"), 3)
}

func testLoggerRoutineAppend(t *testing.T, logger *Logger, baseN int) {
    mwg := sync.WaitGroup{}
    pwg := sync.WaitGroup{}

    mwg.Add(4)
    pwg.Add(1)

    go func() {
        defer mwg.Done()

        pwg.Done()

        logger.Infof("Formated integer %d", 10)
    }()

    pwg.Wait()
    pwg.Add(1)

    go func() {
        defer mwg.Done()

        pwg.Done()

        time.Sleep(10 * time.Millisecond)

        logger.Warningf("Formated integer %d", 10)
    }()

    pwg.Wait()
    pwg.Add(1)

    go func() {
        defer mwg.Done()

        pwg.Done()

        time.Sleep(20 * time.Millisecond)

        logger.Errorf("Formated integer %d", 10)
    }()

    pwg.Wait()

    go func() {
        defer mwg.Done()

        time.Sleep(30 * time.Millisecond)

        logger.Write([]byte("This is a slice of bytes"))
    }()

    mwg.Wait()

    dumpped := logger.Dump()

    if len(dumpped) != 4 * baseN {
        t.Error("Unexpected amount of log items")

        return
    }
}

func testLoggerRegisterNAppendMutex(t *testing.T) {
    logger := NewLogger()
    printer := &fakeLogPrinterForMutex{}

    logger.Register(printer)

    testLoggerRoutineAppend(t, logger, 1)

    if printer.order[0] != "Info" ||
        printer.order[1] != "Warning" ||
        printer.order[2] != "Error" ||
        printer.order[3] != "Print" {
        t.Error("Invalid log order")

        return
    }
}

func TestLoggerAppendMutex(t *testing.T) {
    testLoggerRegisterNAppendMutex(t)
    testLoggerRegisterNAppendMutex(t)
    testLoggerRegisterNAppendMutex(t)
    testLoggerRegisterNAppendMutex(t)
    testLoggerRegisterNAppendMutex(t)
}
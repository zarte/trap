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

package types

import (
    "testing"
    "time"
    "sync"
)

func testMutexExec(t *testing.T) {
    mx      := Mutex{}
    tWG1    := sync.WaitGroup{}
    tWG2    := sync.WaitGroup{}
    val     := []byte{}

    tWG1.Add(1)
    tWG2.Add(1)

    go mx.Exec(func() { // Callback 1
        defer tWG2.Done()

        tWG1.Done()

        time.Sleep(100 * time.Millisecond)

        val = append(val, 'A')
    })

    tWG1.Wait() // We can confirm the callback 1 is running

    tWG2.Add(1)

    go mx.Exec(func() { // Callback 2
        defer tWG2.Done()

        val = append(val, 'B')
    })

    tWG2.Wait()

    if string(val) != "AB" {
        t.Error("Mutex failed to control the running order")
    }
}

func TestMutexExec(t *testing.T) {
    // Test 3 times to avoid any false negative
    testMutexExec(t)
    testMutexExec(t)
    testMutexExec(t)
}

func testMutexRoutineExec(t *testing.T) {
    mx      := Mutex{}
    tWG1    := sync.WaitGroup{}
    tWG2    := sync.WaitGroup{}
    val     := []byte{}

    tWG1.Add(1)
    tWG2.Add(1)

    mx.RoutineExec(func() { // Callback 1
        defer tWG2.Done()

        tWG1.Done()

        time.Sleep(100 * time.Millisecond)

        val = append(val, 'A')
    })

    tWG1.Wait() // We can confirm the callback 1 is running

    tWG2.Add(1)

    mx.RoutineExec(func() { // Callback 2
        defer tWG2.Done()

        val = append(val, 'B')
    })

    tWG2.Wait()

    if string(val) != "AB" {
        t.Error("RoutineExec failed to control the running order")
    }
}

func TestMutexRoutineExec(t *testing.T) {
    // Test 3 times to avoid any false negative
    testMutexRoutineExec(t)
    testMutexRoutineExec(t)
    testMutexRoutineExec(t)
}
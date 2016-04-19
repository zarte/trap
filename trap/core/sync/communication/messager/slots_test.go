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

package messager

import (
	"github.com/raincious/trap/trap/core/types"

	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestslotsHold(t *testing.T) {
	wait := sync.WaitGroup{}
	messages := slots{}
	cleared := uint64(0)
	expired := uint64(0)
	taken := uint64(0)
	errored := uint64(0)
	unknown := uint64(0)
	total := uint64(0)

	messages.Init()

	for i := uint(0); i < uint(6000); i++ {
		wait.Add(1)

		go func(idx uint) {
			defer func() {
				wait.Done()

				atomic.AddUint64(&total, 1)
			}()

			message := &message{}

			err := messages.Hold(message, 2*time.Second,
				func(reason MessageDeleteReason, err *types.Throw) {
					switch reason {
					case MESSAGE_DELETE_REASON_CLEAR:
						atomic.AddUint64(&cleared, 1)

					case MESSAGE_DELETE_REASON_EXPIRE:
						atomic.AddUint64(&expired, 1)

					case MESSAGE_DELETE_REASON_TAKEN:
						atomic.AddUint64(&taken, 1)

					default:
						atomic.AddUint64(&unknown, 1)
					}
				})

			if err != nil {
				atomic.AddUint64(&errored, 1)

				return
			}

			wait.Add(1)

			go func(msgID byte) {
				defer wait.Done()

				time.Sleep(1 * time.Second)

				messages.Take(msgID)
			}(message.ID)
		}(i)
	}

	<-time.After(10 * time.Second)

	messages.Deinit()

	wait.Wait()

	if total != (errored + cleared + expired + taken + unknown) {
		t.Errorf("Total: %d, Handled: %d (Errored: %d, Cleared: %d, "+
			"Expired: %d, Taken: %d, Unknown: %d)",
			total, (errored + cleared + expired + taken + unknown),
			errored, cleared, expired, taken, unknown)
	}
}

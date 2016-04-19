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
)

type message struct {
    ID                          byte
    Code                        byte
    Held                        bool
    Message                     [][]byte
    Responds                    Callbacks
    ResultChan                  chan *types.Throw
    MaxRespondLen               uint
    Ready                       bool
    StatusLock                  types.Mutex
}
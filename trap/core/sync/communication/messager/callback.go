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

var (
	ErrCallbackAlreadyExisted *types.Error = types.NewError(
		"`Sync` Callback for signal '%d' is already existed")

	ErrCallbackNotFound *types.Error = types.NewError(
		"`Sync` Callback for signal '%s' is not found")
)

type Callback func(Request) *types.Throw
type Callbacks [uint16(^byte(0)) + 1]Callback

func (c *Callbacks) Has(code byte) bool {
	if c[code] == nil {
		return false
	}

	return true
}

func (c *Callbacks) Register(code byte, callback Callback) *types.Throw {
	if c.Has(code) {
		return ErrCallbackAlreadyExisted.Throw(code)
	}

	c[code] = callback

	return nil
}

func (c *Callbacks) Call(code byte, request Request) *types.Throw {
	if !c.Has(code) {
		return ErrCallbackNotFound.Throw(code)
	}

	return c[code](request)
}

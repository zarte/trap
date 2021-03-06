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

package listen

import (
	"github.com/raincious/trap/trap/core/types"
)

var (
	ErrProtocolAlreadyRegistered *types.Error = types.NewError(
		"Protocol '%s' already been registered")

	ErrProtocolNotSupported *types.Error = types.NewError(
		"Protocol '%s' is not supported")

	ErrProtocolPrototypeNotImplmented *types.Error = types.NewError(
		"You can't use a prototype protocol")

	ErrListenerAlreadyInited *types.Error = types.NewError(
		"Listener already initialized")

	ErrListenerAlreadyUp *types.Error = types.NewError(
		"Listener '%s' is already up")

	ErrListenerNotCloseable *types.Error = types.NewError(
		"Listener which listening '%s' is not closeable")

	ErrProtocolAlreadyInited *types.Error = types.NewError(
		"This protocol already beed initialized")
)

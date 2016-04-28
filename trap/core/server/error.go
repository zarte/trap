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

package server

import (
	"github.com/raincious/trap/trap/core/types"
)

var (
	ErrServerNotYetStarted *types.Error = types.NewError(
		"`Server` is not yet started")

	ErrServerAlreadyUp *types.Error = types.NewError(
		"`Server` is already up")

	ErrServerIsBooting *types.Error = types.NewError(
		"`Server` is currently booting up")

	ErrClientNotFound *types.Error = types.NewError(
		"Can't found client '%s'")

	ErrClientAlreadyExisted *types.Error = types.NewError(
		"Client '%s' already existed")

	ErrInvalidClientAddress *types.Error = types.NewError(
		"Client Address '%s' is invalid")

	ErrInvalidServerAddress *types.Error = types.NewError(
		"Server Address '%s:%d' is invalid")

	ErrInvalidConnectionType *types.Error = types.NewError(
		"Connection Type for client '%s' is invalid")
)

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

package status

import (
	"github.com/raincious/trap/trap/core/types"
)

var (
	ErrServerNotSet *types.Error = types.NewError(
		"Trap `Server` must be set before start a `Status` server")

	ErrServerAlreadyUp *types.Error = types.NewError(
		"`Status` server is already up")

	ErrServerNotDownable *types.Error = types.NewError(
		"`Status` can't be down at this moment")

	ErrUIDirNotFound *types.Error = types.NewError(
		"'%d' is not a directory")

	ErrFailedGenerateSessionKey *types.Error = types.NewError(
		"Can't generate session key for user '%s'")

	ErrRequestedURLIsNotImplemented *types.Error = types.NewError(
		"The handler for URL '%s' is not implemented")

	ErrUnsupportedClient *types.Error = types.NewError(
		"Trap `Server` can't handle request sent by this client")

	ErrRequestMethodNotImplemented *types.Error = types.NewError(
		"Request method '%s' is not implemented")

	ErrAccountNotFound *types.Error = types.NewError(
		"Account '%s' is not found")

	ErrAccountAlreadyExisted *types.Error = types.NewError(
		"Account '%s' already existed")

	ErrSessionNotFound *types.Error = types.NewError(
		"Session key '%s' does not binded with user '%s'")

	ErrSessionKeyNotFound *types.Error = types.NewError(
		"Can't found Session with key '%s'")

	ErrSessionExpired *types.Error = types.NewError(
		"The Session key '%s' for user '%s' has expired")

	ErrSessionInvalid *types.Error = types.NewError(
		"The session key for user '%s' is invalid")

	ErrInvalidServerPassword *types.Error = types.NewError(
		"User '%s' provided an invalid server password '%s'")

	ErrInvalidUserIPAddress *types.Error = types.NewError(
		"'%s' is not a ip address")

	ErrSessionLoginReqiured *types.Error = types.NewError(
		"Login is required for access this function")

	ErrSessionNoPermission *types.Error = types.NewError(
		"Current session doesn't have permission to perform this operation")

	ErrStatusControllerInvalidParameter *types.Error = types.NewError(
		"Invalid parameter")
)

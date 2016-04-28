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

package controller

import (
	"github.com/raincious/trap/trap/core/status"
	"github.com/raincious/trap/trap/core/types"

	"encoding/json"
	"net"
	"net/http"
)

type Auth struct {
	JSON

	Verify func(net.IP, types.String) (*status.Session, *types.Throw)
	Auth   func(net.IP, types.String) (*status.Session, *types.Throw)
}

func (a *Auth) Get(w http.ResponseWriter, r *http.Request) {
	userIP, userIPErr := a.GetIPFormString(r.RemoteAddr)

	if userIPErr != nil {
		a.Error(status.ErrorRespond{
			Code:  400,
			Error: userIPErr,
		}, w, r)

		return
	}

	sessionInfo, vErr := a.Verify(userIP, types.String(
		r.Header.Get(status.STATUS_SERVER_SESSION_KEY_HEADER)))

	if vErr != nil {
		a.Error(status.ErrorRespond{
			Code:  403,
			Error: vErr,
		}, w, r)

		return
	}

	responseData, jResErr := json.Marshal(status.AuthRespond{
		Token:       sessionInfo.Key,
		Permissions: sessionInfo.Account().Permissions(),
	})

	if jResErr != nil {
		a.Error(status.ErrorRespond{
			Code:  500,
			Error: types.ConvertError(jResErr),
		}, w, r)

		return
	}

	a.WriteGZIP(200, responseData, w, r)
}

func (a *Auth) Post(w http.ResponseWriter, r *http.Request) {
	var authField status.AuthRequest

	decoder := json.NewDecoder(r.Body)

	decodeErr := decoder.Decode(&authField)

	if decodeErr != nil {
		a.Error(status.ErrorRespond{
			Code:  400,
			Error: types.ConvertError(decodeErr),
		}, w, r)

		return
	}

	userIP, userIPErr := a.GetIPFormString(r.RemoteAddr)

	if userIPErr != nil {
		a.Error(status.ErrorRespond{
			Code:  400,
			Error: userIPErr,
		}, w, r)

		return
	}

	sessionData, vErr := a.Auth(userIP, authField.Password)

	if vErr != nil {
		a.Error(status.ErrorRespond{
			Code:  403,
			Error: vErr,
		}, w, r)

		return
	}

	responseData, jResErr := json.Marshal(status.AuthRespond{
		Token:       sessionData.Key,
		Permissions: sessionData.Account().Permissions(),
	})

	if jResErr != nil {
		a.Error(status.ErrorRespond{
			Code:  500,
			Error: types.ConvertError(jResErr),
		}, w, r)

		return
	}

	a.WriteGZIP(200, responseData, w, r)
}

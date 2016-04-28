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

	"net"
	"net/http"
)

type Sessioned struct {
	Default

	Verify func(net.IP, types.String) (*status.Session, *types.Throw)

	session *status.Session
}

func (s *Sessioned) Session() *status.Session {
	return s.session
}

func (s *Sessioned) GetSessionID(r *http.Request) types.String {
	headerID := r.Header.Get(
		status.STATUS_SERVER_SESSION_KEY_HEADER)

	if headerID != "" {
		return types.String(headerID)
	}

	return types.String("")
}

func (s *Sessioned) Before(w http.ResponseWriter,
	r *http.Request) *types.Throw {
	userIP, userIPErr := s.GetIPFormString(r.RemoteAddr)

	if userIPErr != nil {
		s.Error(status.ErrorRespond{
			Code:  400,
			Error: userIPErr,
		}, w, r)

		return userIPErr
	}

	session, vErr := s.Verify(userIP, s.GetSessionID(r))

	if vErr != nil {
		s.Error(status.ErrorRespond{
			Code:  403,
			Error: vErr,
		}, w, r)

		return vErr
	}

	s.session = session

	w.Header().Set("Cache-Control", "private, max-age=0, no-cache")

	return nil
}

type SessionedJSON struct {
	JSON

	Verify func(net.IP, types.String) (*status.Session, *types.Throw)

	session *status.Session
}

func (s *SessionedJSON) Session() *status.Session {
	return s.session
}

func (s *SessionedJSON) GetSessionID(r *http.Request) types.String {
	headerID := r.Header.Get(
		status.STATUS_SERVER_SESSION_KEY_HEADER)

	if headerID != "" {
		return types.String(headerID)
	}

	return types.String("")
}

func (s *SessionedJSON) Before(w http.ResponseWriter,
	r *http.Request) *types.Throw {
	beforeErr := s.JSON.Before(w, r)

	if beforeErr != nil {
		s.Error(status.ErrorRespond{
			Code:  500,
			Error: beforeErr,
		}, w, r)

		return beforeErr
	}

	userIP, userIPErr := s.GetIPFormString(r.RemoteAddr)

	if userIPErr != nil {
		s.Error(status.ErrorRespond{
			Code:  400,
			Error: userIPErr,
		}, w, r)

		return userIPErr
	}

	session, vErr := s.Verify(userIP, s.GetSessionID(r))

	if vErr != nil {
		s.Error(status.ErrorRespond{
			Code:  403,
			Error: vErr,
		}, w, r)

		return vErr
	}

	s.session = session

	w.Header().Set("Cache-Control", "private, max-age=0, no-cache")

	return nil
}

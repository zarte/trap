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
	"github.com/raincious/trap/trap/core/logger"
	"github.com/raincious/trap/trap/core/status"
	"github.com/raincious/trap/trap/core/types"

	"encoding/json"
	"net/http"
)

type Logs struct {
	SessionedJSON

	GetLogs func() []logger.LogExport
}

func (l *Logs) Before(w http.ResponseWriter,
	r *http.Request) *types.Throw {
	parentBefore := l.SessionedJSON.Before(w, r)

	if parentBefore != nil {
		return parentBefore
	}

	session := l.Session()

	if session == nil {
		l.Error(status.ErrorRespond{
			Code:  401,
			Error: status.ErrSessionLoginReqiured.Throw(),
		}, w, r)

		return status.ErrSessionLoginReqiured.Throw()
	}

	if !session.Account().Allowed("logs") {
		l.Error(status.ErrorRespond{
			Code:  403,
			Error: status.ErrSessionNoPermission.Throw(),
		}, w, r)

		return status.ErrSessionNoPermission.Throw()
	}

	return nil
}

func (l *Logs) Get(w http.ResponseWriter, r *http.Request) {
	jsonData, jsonErr := json.Marshal(l.GetLogs())

	if jsonErr != nil {
		l.Error(status.ErrorRespond{
			Code:  500,
			Error: types.ConvertError(jsonErr),
		}, w, r)

		return
	}

	l.WriteGZIP(200, jsonData, w, r)
}

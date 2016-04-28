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
	"net/http"
)

type JSON struct {
	Default
}

func (j *JSON) Before(w http.ResponseWriter,
	r *http.Request) *types.Throw {
	w.Header().Set("Access-Control-Allow-Origin", "Origin")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	return nil
}

func (j *JSON) Error(err status.ErrorRespond, w http.ResponseWriter,
	r *http.Request) {
	jsonData, jsonErr := json.Marshal(err)

	if jsonErr != nil {
		j.WriteGZIP(500,
			[]byte("{\"error\": \"Can't parse JSON data\"}"), w, r)

		return
	}

	j.WriteGZIP(err.Code, jsonData, w, r)
}

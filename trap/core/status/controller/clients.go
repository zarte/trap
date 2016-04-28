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
	"github.com/raincious/trap/trap/core/client"
	"github.com/raincious/trap/trap/core/status"
	"github.com/raincious/trap/trap/core/types"

	"encoding/json"
	"net/http"
)

type Clients struct {
	SessionedJSON

	GetClients func() []client.ClientExport
}

func (c *Clients) Get(w http.ResponseWriter, r *http.Request) {
	jsonData, jsonErr := json.Marshal(c.GetClients())

	if jsonErr != nil {
		c.Error(status.ErrorRespond{
			Code:  500,
			Error: types.ConvertError(jsonErr),
		}, w, r)

		return
	}

	c.WriteGZIP(200, jsonData, w, r)
}

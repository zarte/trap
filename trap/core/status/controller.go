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

    "net/http"
)

type Controller interface {
    Get(http.ResponseWriter, *http.Request)
    Post(http.ResponseWriter, *http.Request)
    Put(http.ResponseWriter, *http.Request)
    Delete(http.ResponseWriter, *http.Request)
    Head(http.ResponseWriter, *http.Request)
    Options(http.ResponseWriter, *http.Request)

    Init() (*types.Throw)
    Before(http.ResponseWriter, *http.Request) (*types.Throw)
    Error(ErrorRespond, http.ResponseWriter, *http.Request)
}
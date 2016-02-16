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

    "net/http"
)

type Home struct {
    Default

    StaticPage      []byte
}

func (s *Home) Get(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        s.Error(status.ErrorRespond{
            Code:   404,
            Error:  status.ErrRequestedURLIsNotImplemented.Throw(
                        r.URL.Path),
        }, w, r)

        return
    }

    if !s.IsGZIPSupported(r) {
        s.Error(status.ErrorRespond{
            Code:   406,
            Error:  status.ErrUnsupportedClient.Throw(),
        }, w, r)

        return
    }

    w.Header().Set("Content-Encoding", "gzip")
    w.Write(s.StaticPage)
}
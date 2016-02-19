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
    "github.com/raincious/trap/trap/core/types"
    "github.com/raincious/trap/trap/core/status"

    "time"
    "net/http"
)

type Home struct {
    Default

    bootTime            time.Time
    formatedBootTime    string

    staticExpired       time.Duration

    staticPage          []byte
}

func (s *Home) Init() (*types.Throw) {
    initErr     := s.Default.Init()

    if initErr != nil {
        return initErr
    }

    s.bootTime              =   time.Now()
    s.formatedBootTime      =   s.bootTime.Format(time.RFC1123)

    s.staticExpired         =   86400 * time.Second

    s.staticPage            =   []byte(status.STATUS_HOME_STATIC_PAGE)

    return nil
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

    if s.IsUnmodified(s.bootTime, r) {
        w.WriteHeader(304)

        return
    }

    w.Header().Set("Content-Encoding",  "gzip")
    w.Header().Set("Pragma",            "cache")
    w.Header().Set("Cache-Control",     "private, max-age=86400")
    w.Header().Set("Last-Modified",     s.formatedBootTime)
    w.Header().Set("Expires",           time.Now().
                                            Add(s.staticExpired).
                                            Format(time.RFC1123))

    w.Write(s.staticPage)
}
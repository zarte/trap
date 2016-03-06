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

type Mux struct {
    server          *http.ServeMux
}

func NewMux() (*Mux) {
    return &Mux{
        server:     http.NewServeMux(),
    }
}

func (mux *Mux) Handle(pattern string, handler http.Handler) {
    mux.server.Handle(pattern, handler)
}

func (mux *Mux) HandleFunc(pattern string,
    handler func(http.ResponseWriter, *http.Request)) {
    mux.server.HandleFunc(pattern, handler)
}

func (mux *Mux) HandleController(pattern string,
    c Controller) (*types.Throw) {
    err     :=  c.Init()

    if err != nil {
        return err
    }

    mux.server.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
        mux.dispatchController(w, r, c)
    })

    return nil
}

func (mux *Mux) dispatchController(w http.ResponseWriter, r *http.Request,
    c Controller) {

    err := c.Before(w, r)

    if err != nil {
        return
    }

    switch r.Method {
        case "GET":
            c.Get(w, r)

        case "POST":
            c.Post(w, r)

        case "PUT":
            c.Put(w, r)

        case "DELETE":
            c.Delete(w, r)

        case "HEAD":
            c.Head(w, r)

        case "OPTIONS":
            c.Options(w, r)
    }
}

func (mux *Mux) Handler(r *http.Request) (h http.Handler, pattern string) {
    return mux.server.Handler(r)
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    mux.server.ServeHTTP(w, r)
}
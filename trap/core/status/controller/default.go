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

    "net"
    "net/http"
    "compress/gzip"
)

type Default struct {

}

func (d *Default) Init() (*types.Throw) {
    return nil
}

func (d *Default) Get(w http.ResponseWriter, r *http.Request) {
    d.Error(status.ErrorRespond{
        Code:       405,
        Error:      status.ErrRequestMethodNotImplemented.Throw(r.Method),
    }, w, r)
}

func (d *Default) Post(w http.ResponseWriter, r *http.Request) {
    d.Error(status.ErrorRespond{
        Code:       405,
        Error:      status.ErrRequestMethodNotImplemented.Throw(r.Method),
    }, w, r)
}

func (d *Default) Put(w http.ResponseWriter, r *http.Request) {
    d.Error(status.ErrorRespond{
        Code:       405,
        Error:      status.ErrRequestMethodNotImplemented.Throw(r.Method),
    }, w, r)
}

func (d *Default) Delete(w http.ResponseWriter, r *http.Request) {
    d.Error(status.ErrorRespond{
        Code:       405,
        Error:      status.ErrRequestMethodNotImplemented.Throw(r.Method),
    }, w, r)
}

func (d *Default) Head(w http.ResponseWriter, r *http.Request) {
    d.Error(status.ErrorRespond{
        Code:       405,
        Error:      status.ErrRequestMethodNotImplemented.Throw(r.Method),
    }, w, r)
}

func (d *Default) Options(w http.ResponseWriter, r *http.Request) {
    d.Error(status.ErrorRespond{
        Code:       405,
        Error:      status.ErrRequestMethodNotImplemented.Throw(r.Method),
    }, w, r)
}

func (d *Default) Before(w http.ResponseWriter,
    r *http.Request) (*types.Throw) {
    w.Header().Set("Content-Type", "text/html; charset=UTF-8")

    return nil
}

func (d *Default) Error(err status.ErrorRespond, w http.ResponseWriter,
    r *http.Request) {
    d.WriteGZIP(err.Code, []byte(err.Error.Error()), w, r)
}

func (d *Default) GetIPFormString(addr string) (net.IP, *types.Throw) {
    userHost, _, ipSplitErr := net.SplitHostPort(addr)

    if ipSplitErr != nil {
        return nil, types.ConvertError(ipSplitErr)
    }

    userIP := net.ParseIP(userHost)

    if userIP == nil {
        return nil, status.ErrInvalidUserIPAddress.Throw(userHost)
    }

    return userIP, nil
}

func (d *Default) IsGZIPSupported(r *http.Request) (bool) {
    clientEncodings :=  types.String(r.Header.Get("Accept-Encoding")).Lower()

    if !clientEncodings.Contains("gzip") {
        return false
    }

    return true
}

func (d *Default) WriteGZIP(code int, data []byte, w http.ResponseWriter,
    r *http.Request) (*types.Throw) {
    if !d.IsGZIPSupported(r) {
        w.WriteHeader(code)

        w.Write(data)

        return nil
    }

    w.Header().Set("Content-Encoding", "gzip")

    w.WriteHeader(code)

    gzipWriter      :=  gzip.NewWriter(w)

    defer gzipWriter.Close()

    _, wError       :=  gzipWriter.Write(data)

    if wError != nil {
        return types.ConvertError(wError)
    }

    return nil
}
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
    "testing"
    "bytes"
    "errors"
)

type dummyHttpHandler struct {
    called          bool
}

func (d *dummyHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    d.called        =   true
}

type dummyResponseWriter struct {
    headers         http.Header
}

func (d dummyResponseWriter) Header() (http.Header) {
    return d.headers
}

func (d dummyResponseWriter) Write(data []byte) (int, error) {
    return 0, errors.New("Can't write dummy responder")
}

func (d dummyResponseWriter) WriteHeader(code int) {}

type dummyController struct {
    initCalled      bool
    beforeCalled    bool

    getCalled       bool
    postCalled      bool
    putCalled       bool
    deleteCalled    bool
    headCalled      bool
    optionsCalled   bool
}

func (d *dummyController) Init() (*types.Throw) {
    d.initCalled    =   true

    return nil
}

func (d *dummyController) Get(w http.ResponseWriter, r *http.Request) {
    d.getCalled     =   true
}

func (d *dummyController) Post(w http.ResponseWriter, r *http.Request) {
    d.postCalled    =   true
}

func (d *dummyController) Put(w http.ResponseWriter, r *http.Request) {
    d.putCalled     =   true
}

func (d *dummyController) Delete(w http.ResponseWriter, r *http.Request) {
    d.deleteCalled  =   true
}

func (d *dummyController) Head(w http.ResponseWriter, r *http.Request) {
    d.headCalled    =   true
}

func (d *dummyController) Options(w http.ResponseWriter, r *http.Request) {
    d.optionsCalled =   true
}

func (d *dummyController) Before(w http.ResponseWriter,
    r *http.Request) (*types.Throw) {
    d.beforeCalled  =   true

    return nil
}

func (d *dummyController) Error(err ErrorRespond, w http.ResponseWriter,
    r *http.Request) {}

func TestMuxNewMux(t *testing.T) {
    dummyHandler    :=  &dummyHttpHandler{
        called:         false,
    }
    dummyCrtl       :=  &dummyController{

    }

    resp            :=  dummyResponseWriter{
        headers:        http.Header{},
    }
    req, reqErr     :=  http.NewRequest("GET", "/handler",
                            bytes.NewBufferString(""))

    if reqErr != nil {
        t.Errorf("Failed to create request struct for this test due to error: %s",
            reqErr)

        return
    }

    mux             :=  NewMux()

    mux.Handle("/handler", dummyHandler)
    mux.HandleFunc("/handlefunc", func(w http.ResponseWriter, r *http.Request) {

    })
    mux.HandleController("/handlecontroller", dummyCrtl)

    mux.ServeHTTP(resp, req)
}

func getNewHTTPRequest(method string, url string) (*http.Request) {
    req, _          :=  http.NewRequest(method, url,
                            bytes.NewBufferString(""))

    return req
}

func TestMuxHandle(t *testing.T) {
    resp            :=  dummyResponseWriter{
        headers:        http.Header{},
    }

    dummyHandler    :=  &dummyHttpHandler{
        called:         false,
    }

    mux             :=  NewMux()

    mux.Handle("/handler", dummyHandler)

    for _, method := range []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"} {
        dummyHandler.called     =   false

        mux.ServeHTTP(resp, getNewHTTPRequest(method, "/handler"))

        if !dummyHandler.called {
            t.Error("Failed asserting that the handler is called")

            return
        }
    }
}

func TestMuxHandleFunc(t *testing.T) {
    called          :=  false
    resp            :=  dummyResponseWriter{
        headers:        http.Header{},
    }

    mux             :=  NewMux()

    mux.HandleFunc("/handler", func(w http.ResponseWriter, r *http.Request) {
        called      =   true
    })

    for _, method := range []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"} {
        called      =   false

        mux.ServeHTTP(resp, getNewHTTPRequest(method, "/handler"))

        if !called {
            t.Error("Failed asserting that the handler is called")

            return
        }
    }
}

func TestMuxHandleController(t *testing.T) {
    resp            :=  dummyResponseWriter{
        headers:        http.Header{},
    }
    dummyCrtl       :=  &dummyController{
        initCalled:     false,
        beforeCalled:   false,
        getCalled:      false,
        postCalled:     false,
        putCalled:      false,
        deleteCalled:   false,
        headCalled:     false,
        optionsCalled:  false,
    }

    mux             :=  NewMux()

    mux.HandleController("/handler", dummyCrtl)

    for _, method := range []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"} {
        mux.ServeHTTP(resp, getNewHTTPRequest(method, "/handler"))
    }

    if !dummyCrtl.initCalled || !dummyCrtl.getCalled ||
        !dummyCrtl.postCalled || !dummyCrtl.putCalled ||
        !dummyCrtl.deleteCalled || !dummyCrtl.headCalled ||
        !dummyCrtl.optionsCalled {
            t.Error("Controller call is incompleted")

            return
        }
}
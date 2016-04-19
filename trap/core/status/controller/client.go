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
    "github.com/raincious/trap/trap/core/types"
    "github.com/raincious/trap/trap/core/status"
    "github.com/raincious/trap/trap/core/server"

    "net/http"
    "encoding/json"
)

type Client struct {
    SessionedJSON

    GetClient   func(addr types.IP) (*client.Client, *types.Throw)
    AddClient   func(cCon server.ClientInfo) (*client.Client, *types.Throw)
    DelClient   func(addr types.IP) (*types.Throw)
}

func (c *Client) Before(w http.ResponseWriter,
    r *http.Request) (*types.Throw) {
    parentBefore                    :=  c.SessionedJSON.Before(w, r)

    if parentBefore != nil {
        return parentBefore
    }

    session                         :=  c.Session()

    if session == nil {
        c.Error(status.ErrorRespond{
            Code:   401,
            Error:  status.ErrSessionLoginReqiured.Throw(),
        }, w, r)

        return status.ErrSessionLoginReqiured.Throw()
    }

    if !session.Account().Allowed("clients") {
        c.Error(status.ErrorRespond{
            Code:   403,
            Error:  status.ErrSessionNoPermission.Throw(),
        }, w, r)

        return status.ErrSessionNoPermission.Throw()
    }

    return nil
}

func (c *Client) Get(w http.ResponseWriter, r *http.Request) {
    clientIP, clientIPErr           :=  types.ConvertIPFromString(
                                            types.String(
                                                r.URL.Query().Get("client")))

    if clientIPErr != nil {
        c.Error(status.ErrorRespond{
            Code:   400,
            Error:  status.ErrStatusControllerInvalidParameter.Throw(),
        }, w, r)

        return
    }

    client, clientErr               :=  c.GetClient(clientIP)

    if clientErr != nil {
        c.Error(status.ErrorRespond{
            Code:   404,
            Error:  clientErr,
        }, w, r)

        return
    }

    jsonData, jsonErr   :=  json.Marshal(client)

    if jsonErr != nil {
        c.Error(status.ErrorRespond{
            Code: 500,
            Error: types.ConvertError(jsonErr),
        }, w, r)

        return
    }

    c.WriteGZIP(200, jsonData, w, r)
}

func (c *Client) Post(w http.ResponseWriter, r *http.Request) {
    var clientConField              server.ClientInfo

    clientIP, clientIPErr           :=  types.ConvertIPFromString(
                                            types.String(
                                                r.URL.Query().Get("client")))

    if clientIPErr != nil {
        c.Error(status.ErrorRespond{
            Code:   400,
            Error:  status.ErrStatusControllerInvalidParameter.Throw(),
        }, w, r)

        return
    }

    decoder                         :=  json.NewDecoder(r.Body)

    decodeErr                       :=  decoder.Decode(&clientConField)

    if decodeErr != nil {
        c.Error(status.ErrorRespond{
            Code:   400,
            Error:  types.ConvertError(decodeErr),
        }, w, r)

        return
    }

    clientConField.Client           =   clientIP

    clientInfo, clientErr           :=  c.AddClient(clientConField)

    if clientErr != nil {
        code                        :=  500

        if clientErr.Is(server.ErrClientAlreadyExisted) {
            code                    =   409
        } else if clientErr.Is(server.ErrClientNotFound) {
            code                    =   404
        } else if clientErr.Is(server.ErrInvalidConnectionType) {
            code                    =   400
        } else if clientErr.Is(server.ErrInvalidServerAddress) {
            code                    =   400
        } else if clientErr.Is(server.ErrInvalidClientAddress) {
            code                    =   400
        }

        c.Error(status.ErrorRespond{
            Code:   code,
            Error:  clientErr,
        }, w, r)

        return
    }

    jsonData, jsonErr               :=  json.Marshal(clientInfo)

    if jsonErr != nil {
        c.Error(status.ErrorRespond{
            Code: 500,
            Error: types.ConvertError(jsonErr),
        }, w, r)

        return
    }

    c.WriteGZIP(201, jsonData, w, r)
}

func (c *Client) Delete(w http.ResponseWriter, r *http.Request) {
    clientIP, clientIPErr           :=  types.ConvertIPFromString(
                                            types.String(
                                                r.URL.Query().Get("client")))

    if clientIPErr != nil {
        c.Error(status.ErrorRespond{
            Code:   400,
            Error:  status.ErrStatusControllerInvalidParameter.Throw(),
        }, w, r)

        return
    }

    clientDeleteErr                 :=  c.DelClient(clientIP)

    if clientDeleteErr != nil {
        c.Error(status.ErrorRespond{
            Code:   404,
            Error:  clientDeleteErr,
        }, w, r)

        return
    }

    jsonData, jsonErr               :=  json.Marshal(status.ClientDeletedRespond{
        Result:                         true,
    })

    if jsonErr != nil {
        c.Error(status.ErrorRespond{
            Code: 500,
            Error: types.ConvertError(jsonErr),
        }, w, r)

        return
    }

    c.WriteGZIP(200, jsonData, w, r)
}

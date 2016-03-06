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

    "net"
    "time"
    "crypto/rand"
    "encoding/base64"
)

type Session struct {
    IP                              net.IP
    Created                         time.Time
    Key                             types.String

    LastSeen                        time.Time
    Expire                          time.Duration

    account                         *Account
}

type SessionDump struct {
    IP                              net.IP
    Created                         time.Time

    LastSeen                        time.Time
    Expire                          time.Duration

    Permissions                     map[types.String]bool
}

func (s *Session) Bump() {
    s.LastSeen                      =   time.Now()
}

func (s *Session) Account() (*Account) {
    return s.account
}

func (s *Session) Expired() (bool) {
    if s.LastSeen.Add(s.Expire).After(time.Now()) {
        return false
    }

    return true
}

type Sessions map[types.String]*Session

func (s Sessions) getRandomKey() (types.String) {
    rBytes                          :=  make([]byte, 32)

    _, rErr                         :=  rand.Read(rBytes)

    if rErr != nil {
        return types.String("")
    }

    return types.String(base64.URLEncoding.EncodeToString(rBytes))
}

func (s Sessions) scanExpired() {
    for key, val := range s {
        if !val.Expired() {
            continue
        }

        delete(s, key)
    }
}

func (s Sessions) Add(ip net.IP, account *Account,
    expire time.Duration) (*Session, *types.Throw) {
    maxRetry                        :=  3
    ipStr                           :=  types.String(ip.String())

    // Add new session for user
    newSession                      :=  &Session{
        IP:                         ip,
        Created:                    time.Now(),
        Key:                        "",

        LastSeen:                   time.Now(),
        Expire:                     expire,

        account:                    account,
    }

    for {
        newSession.Key              =   s.getRandomKey()

        if _, ok := s[ipStr + ":" + newSession.Key]; !ok {
            break
        }

        newSession.Key              =   ""

        if maxRetry <= 0 {
            break
        }

        maxRetry                    -=  1
    }

    if newSession.Key == "" {
        return nil, ErrFailedGenerateSessionKey.Throw(ip)
    }

    s.scanExpired()

    s[ipStr + ":" + newSession.Key] =   newSession

    return s[ipStr + ":" + newSession.Key], nil
}

func (s Sessions) Delete(sessionKey types.String) (*types.Throw) {
    for key, val := range s {
        if val.Key != sessionKey {
            continue
        }

        delete(s, key)

        return nil
    }

    return ErrSessionKeyNotFound.Throw(sessionKey)
}

func (s Sessions) Verify(ip net.IP,
    sessionKey types.String) (*Session, *types.Throw) {
    ipStr                           :=  types.String(ip.String())
    keyName                         :=  ipStr + ":" + sessionKey

    if _, ok := s[keyName]; !ok {
        return nil, ErrSessionNotFound.Throw(sessionKey, ipStr)
    }

    if s[keyName].Expired() {
        delete(s, keyName)

        return nil, ErrSessionExpired.Throw(sessionKey, ipStr)
    }

    s[keyName].Bump()

    return s[keyName], nil
}

func (s Sessions) Dump() ([]SessionDump) {
    result := []SessionDump{}

    for _, sess := range s {
        result = append(result, SessionDump{
            IP:                     sess.IP,
            Created:                sess.Created,

            LastSeen:               sess.LastSeen,
            Expire:                 sess.Expire,

            Permissions:            sess.Account().Permissions(),
        })
    }

    return result
}
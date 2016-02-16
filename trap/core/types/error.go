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

package types

import (
    "fmt"
    "encoding/json"
)

type Error struct {
    ref     *Error
    error   error
    text    string
}

type Throw struct {
    parentRef  *Error
    parentErr  error
    message    string
}

func NewError(text string) (*Error) {
    e           := &Error{}

    e.ref       =  e
    e.text      =  text
    e.error     =  nil

    return e
}

func ConvertError(err error) (*Throw) {
    e           := &Error{}

    e.ref       =  e
    e.text      =  err.Error()
    e.error     =  err

    return e.Throw()
}

func (e *Error) Throw(formats ...interface{}) (*Throw) {
    t           := &Throw{}

    t.parentRef =  e.ref
    t.parentErr =  e.error
    t.message   =  fmt.Sprintf(e.text, formats...)

    return t
}

func (t *Throw) Is(e *Error) (bool) {
    if e.ref != t.parentRef {
        return false
    }

    return true
}

func (t *Throw) IsError(e error) (bool) {
    if e != t.parentErr {
        return false
    }

    return true
}

func (t *Throw) Error() (string) {
    return t.message
}

func (t *Throw) MarshalJSON() ([]byte, error) {
    return json.Marshal(t.Error())
}
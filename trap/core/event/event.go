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

package event

import (
    "github.com/raincious/trap/trap/core/logger"
    "github.com/raincious/trap/trap/core/types"
)

type Callbacks  map[types.String][]func(*Parameters) (*types.Throw)

type Event struct {
    inited              bool

    error               *types.Throw

    logger              *logger.Logger
    events              Callbacks
}

func (this *Event) Init(cfg *Config) {
    this.logger         = cfg.Logger.NewContext("Event")

    this.events         = Callbacks{}
}

func (this *Event) Register(name types.String,
    callback func(*Parameters) (*types.Throw)) {
    this.events[name]   = append(this.events[name], callback)

    this.logger.Debugf("New `Event` handler has been registered to '%s' event",
        name)
}

func (this *Event) Trigger(name types.String,
    params Parameters) (*types.Throw) {
    var e *types.Throw = nil

    if _, ok := this.events[name]; !ok {
        this.error = ErrNoEvent.Throw(name)

        this.logger.Warningf("Can't trigger event due to error: %s", this.error)

        return this.error
    }

    this.logger.Debugf("The event '%s' has been triggered", name)

    for _, eventHandler := range this.events[name] {
        e = eventHandler(&params)

        if e != nil {
            this.logger.Errorf("An error happed when " +
                "run handler for event '%s': %s", name, e)
        }
    }

    return e
}
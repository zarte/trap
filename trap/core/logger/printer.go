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

package logger

import (
	"github.com/raincious/trap/trap/core/types"

	"time"
)

type Printer interface {
	Debug(types.String, time.Time, types.String)
	Info(types.String, time.Time, types.String)
	Warning(types.String, time.Time, types.String)
	Error(types.String, time.Time, types.String)

	Print(types.String, time.Time, types.String)
}

type Printers []Printer

func (l *Printers) Add(newPrinter Printer) {
	oldPrinters := append(*l, newPrinter)

	*l = oldPrinters
}

func (l Printers) Debug(context types.String,
	now time.Time, msg types.String) {
	for _, printer := range l {
		printer.Debug(context, now, msg)
	}
}

func (l Printers) Info(context types.String,
	now time.Time, msg types.String) {
	for _, printer := range l {
		printer.Info(context, now, msg)
	}
}

func (l Printers) Warning(context types.String,
	now time.Time, msg types.String) {
	for _, printer := range l {
		printer.Warning(context, now, msg)
	}
}

func (l Printers) Error(context types.String,
	now time.Time, msg types.String) {
	for _, printer := range l {
		printer.Error(context, now, msg)
	}
}

func (l Printers) Print(context types.String,
	now time.Time, msg types.String) {
	for _, printer := range l {
		printer.Print(context, now, msg)
	}
}

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

	"fmt"
	"time"
)

type ScreenPrinter struct{}

func NewScreenPrinter() *ScreenPrinter {
	return &ScreenPrinter{}
}

func (l *ScreenPrinter) print(w types.String, c types.String,
	t time.Time, m types.String) {
	fmt.Printf("  <%3.3s> %-30.30s [%19.19s]: %s\r\n", w, c,
		t.Format(time.StampMilli), m)
}

func (l *ScreenPrinter) Info(c types.String, t time.Time, m types.String) {
	l.print("INF", c, t, m)
}

func (l *ScreenPrinter) Debug(c types.String, t time.Time, m types.String) {
	l.print("DBG", c, t, m)
}

func (l *ScreenPrinter) Warning(c types.String, t time.Time, m types.String) {
	l.print("WRN", c, t, m)
}

func (l *ScreenPrinter) Error(c types.String, t time.Time, m types.String) {
	l.print("ERR", c, t, m)
}

func (l *ScreenPrinter) Print(c types.String, t time.Time, m types.String) {
	l.print("DEF", c, t, m)
}

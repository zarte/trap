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

	"bufio"
	"fmt"
	"time"
)

type FilePrinter struct {
	writer *bufio.Writer

	writeCounts uint16
}

func NewFilePrinter(w *bufio.Writer) (*FilePrinter, *types.Throw) {
	_, writeErr := w.Write([]byte(""))

	if writeErr != nil {
		return nil, types.ConvertError(writeErr)
	}

	return &FilePrinter{
		writer: w,
	}, nil
}

func (l *FilePrinter) save(w types.String, c types.String,
	t time.Time, m types.String) {

	_, err := l.writer.WriteString(fmt.Sprintf("<%s> %s [%s]: %s\r\n",
		w, c, t.Format(time.StampMilli), m))

	if err != nil {
		panic(fmt.Errorf("Can't write log file due to error: %s", err))
	}

	l.writeCounts += 1

	if l.writeCounts > 10 {
		l.writer.Flush()

		l.writeCounts = 0
	}
}

func (l *FilePrinter) Info(c types.String, t time.Time, m types.String) {
	l.save("INF", c, t, m)
}

func (l *FilePrinter) Debug(c types.String, t time.Time, m types.String) {
	l.save("DBG", c, t, m)
}

func (l *FilePrinter) Warning(c types.String, t time.Time, m types.String) {
	l.save("WRN", c, t, m)
}

func (l *FilePrinter) Error(c types.String, t time.Time, m types.String) {
	l.save("ERR", c, t, m)
}

func (l *FilePrinter) Print(c types.String, t time.Time, m types.String) {
	l.save("DEF", c, t, m)
}

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

type Log struct {
	Time    time.Time
	Type    int
	Context types.String
	Message types.String
}

type LogExport struct {
	Time    time.Time
	Type    types.String
	Context types.String
	Message types.String
}

type Logs []Log

func (l *Logs) Append(log Log, maxLen int) {
	oldLogs := append(*l, log)

	totalLen := len(oldLogs)

	if totalLen > maxLen {
		oldLogs = oldLogs[totalLen-maxLen:]
	}

	*l = oldLogs
}

func (l *Logs) Export() []LogExport {
	exported := []LogExport{}

	oldLogs := *l

	for i := len(oldLogs) - 1; i >= 0; i-- {
		lType := types.String("Default")

		switch oldLogs[i].Type {
		case LOG_TYPE_DEBUG:
			lType = "Debug"

		case LOG_TYPE_INFO:
			lType = "Information"

		case LOG_TYPE_WARNING:
			lType = "Warning"

		case LOG_TYPE_ERROR:
			lType = "Error"
		}

		exported = append(exported, LogExport{
			Time:    oldLogs[i].Time,
			Type:    lType,
			Context: oldLogs[i].Context,
			Message: oldLogs[i].Message,
		})
	}

	return exported
}

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

package responder

import (
	"github.com/raincious/trap/trap/core/listen"
	"github.com/raincious/trap/trap/core/types"
	"github.com/raincious/trap/trap/protocol/tcp"

	"io"
	"net"
)

type Empty struct {
}

func (e *Empty) Handle(conn *net.TCPConn,
	config *tcp.ResponderConfig) (listen.RespondedResult, *types.Throw) {
	readLen := uint(256)

	result := listen.RespondedResult{
		Suggestion:     listen.RESPOND_SUGGEST_MARK,
		ReceivedSample: []byte{},
		RespondedData:  []byte{},
	}

	totalLen := uint(0)
	maxLen := config.MaxBytes

	if maxLen < 256 {
		readLen = maxLen
	}

	for {
		buffer := make([]byte, readLen)

		rLen, rErr := conn.Read(buffer)

		if rErr == io.EOF {
			break
		} else if rErr != nil {
			if totalLen > 0 {
				break
			}

			return result, types.ConvertError(rErr)
		}

		totalLen += uint(rLen)

		if totalLen > maxLen {
			break
		}

		result.ReceivedSample = append(result.ReceivedSample,
			buffer[:rLen]...)
	}

	return result, nil
}

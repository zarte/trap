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
    "github.com/raincious/trap/trap/core/types"
    "github.com/raincious/trap/trap/core/listen"

    "io"
    "net"
)

type Empty struct {

}

func (e *Empty) Handle(conn *net.TCPConn) (listen.RespondedResult, *types.Throw) {
    var totalbuffer []byte

    result                  :=  listen.RespondedResult{
                                    Suggestion: listen.RESPOND_SUGGEST_MARK,
                                }

    totalLen                :=  0
    maxLen                  :=  len(result.ReceivedSample)

    for {
        buffer              :=  make([]byte, 256)

        rLen, rErr          :=  conn.Read(buffer)

        if rErr == io.EOF {
            break
        } else if rErr != nil {
            if totalLen > 0 {
                break;
            }

            return result, types.ConvertError(rErr)
        }

        totalLen            +=  rLen

        if totalLen > maxLen {
            break
        }

        totalbuffer         =   append(totalbuffer, buffer[:rLen]...)
    }

    result.ReceivedLen      =   totalLen

    copy(result.ReceivedSample[:], totalbuffer)

    return result, nil
}
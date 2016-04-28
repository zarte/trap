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

package communication
/*
import (
    "testing"

    "bytes"
)


func TestCommonCombine(t *testing.T) {
    common              :=      Common{}
    data                :=      [][]byte{
        []byte("This is some test data"),
        []byte{SYNC_CONTROLCHAR_TRANSEND, 'A', SYNC_CONTROLCHAR_SEPARATOR},
    }

    combined            :=      common.Combine(data)

    if !bytes.Equal(combined, []byte{84, 104, 105, 115, 32, 105, 115, 32, 115,
        111, 109, 101, 32, 116, 101, 115, 116, 32, 100, 97, 116, 97, 2, 0, 1,
        65, 0, 2, 1}) {
        t.Errorf("Common.Combine() failed to product expected result")

        return
    }
}*/
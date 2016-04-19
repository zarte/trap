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

package data

import (
    "github.com/raincious/trap/trap/core/types"
)

type Partner struct {
    Base

    Removed             types.IPAddresses
    Added               types.IPAddresses
}

func (d *Partner) Parse(msg [][]byte) (*types.Throw) {
    verifyErr           :=  d.Verify(msg, 2)

    if verifyErr != nil {
        return verifyErr
    }

    removedPartnerErr   :=  d.Removed.Unserialize(msg[1])

    if removedPartnerErr != nil {
        return removedPartnerErr
    }

    addedPartnerErr     :=  d.Added.Unserialize(msg[1])

    if addedPartnerErr != nil {
        return addedPartnerErr
    }

    return nil
}

func (d *Partner) Build() ([][]byte, *types.Throw) {
    removedBy, rmBErr   :=  d.Removed.Serialize()

    if rmBErr != nil {
        return [][]byte{}, rmBErr
    }

    addedBy, adBErr     :=  d.Added.Serialize()

    if adBErr != nil {
        return [][]byte{}, adBErr
    }

    return [][]byte{
        removedBy,
        addedBy,
    }, nil
}
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

package status

import (
    //"github.com/raincious/trap/trap/core/types"

    "testing"
)

func TestPermissionAuthorize(t *testing.T) {
    permission := Permission{}

    permission.Authorize("Test permission")

    if !permission.Allowed("Test permission") {
        t.Error("Failed asserting an authorized permission is authorized")

        return
    }
}

func TestPermissionAllowed(t *testing.T) {
    permission := Permission{}

    permission.Authorize("Test permission")

    if !permission.Allowed("Test permission") {
        t.Error("Failed asserting an authorized permission is authorized")

        return
    }

    if permission.Allowed("Another permission") {
        t.Error("Failed asserting an unauthorized permission is unauthorized")

        return
    }
}

func TestPermissionAll(t *testing.T) {
    permission := Permission{}

    permission.Authorize("Test permission 1")
    permission.Authorize("Test permission 2")
    permission.Authorize("Test permission 3")
    permission.Authorize("Test permission 4")

    allPermissions := permission.All()

    // Why 6 permissions?
    // Because we record all permissions in a global shared variable which
    // share the same value across different structs.
    // So the permission we actually have in this table is:
    // Test permission, Another permission and Test permission 1 - 4
    if len(allPermissions) != 6 {
        t.Error("Invalid amount of permissions")

        return
    }

    for pKey, pVal := range allPermissions {
        if permission.Allowed(pKey) == pVal {
           continue
        }

        t.Error("Permission.All() exports an invalid permission")

        return
    }
}
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
	"github.com/raincious/trap/trap/core/status/permission"
	"github.com/raincious/trap/trap/core/types"
)

var (
	permissionRegister = permission.NewRegister()
)

type Permission struct {
	permission types.UInt64
}

func (p *Permission) Authorize(name types.String) {
	permissionVal := permissionRegister.Get(name)

	p.permission = p.permission | permissionVal
}

func (p *Permission) Allowed(name types.String) bool {
	permissionVal := permissionRegister.Get(name)

	if (p.permission & permissionVal) == 0 {
		return false
	}

	return true
}

func (p *Permission) All() map[types.String]bool {
	permissions := map[types.String]bool{}

	for name, permissionVal := range permissionRegister.All() {
		permissions[name] = (p.permission & permissionVal) != 0
	}

	return permissions
}

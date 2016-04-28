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
	"github.com/raincious/trap/trap/core/types"
)

type Account struct {
	permission Permission
}

func (a Account) Allowed(name types.String) bool {
	return a.permission.Allowed(name)
}

func (a Account) Permissions() map[types.String]bool {
	return a.permission.All()
}

type Accounts map[types.String]*Account

func (a Accounts) Get(pass types.String) (*Account, *types.Throw) {
	if _, ok := a[pass]; !ok {
		return nil, ErrAccountNotFound.Throw(pass)
	}

	return a[pass], nil
}

func (a Accounts) Register(
	pass types.String, permissions []types.String) (*Account, *types.Throw) {
	testAccount, testErr := a.Get(pass)

	if testErr != nil && !testErr.Is(ErrAccountNotFound) {
		return nil, testErr
	}

	if testAccount != nil {
		return nil, ErrAccountAlreadyExisted.Throw(pass)
	}

	newAccount := &Account{
		permission: Permission{},
	}

	for _, permissionName := range permissions {
		newAccount.permission.Authorize(permissionName)
	}

	a[pass] = newAccount

	return a[pass], nil
}

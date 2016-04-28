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

	"testing"
)

func TestAccountsRegister(t *testing.T) {
	accounts := Accounts{}

	acc, regErr := accounts.Register("The test pass",
		[]types.String{"Test permission 1", "Test permission 2"})

	if regErr != nil {
		t.Errorf("Accounts.Register() failed register account due to error: %s",
			regErr)

		return
	}

	// Let's by the way test the Account a little bit here
	if !acc.Allowed("Test permission 1") ||
		!acc.Allowed("Test permission 2") ||
		len(acc.Permissions()) == 0 {
		t.Error("Accounts.Register() didn't succefully initialize the Account")

		return
	}

	acc, regErr = accounts.Register("The test pass",
		[]types.String{"Test permission 3", "Test permission 4"})

	if regErr == nil || !regErr.Is(ErrAccountAlreadyExisted) {
		t.Errorf("Accounts.Register() failed register account due to error: %s",
			regErr)

		return
	}

	if acc != nil {
		t.Error("Accounts.Register() mistakenly filled account")

		return
	}
}

func TestAccountsGet(t *testing.T) {
	accounts := Accounts{}

	regAcc, regErr := accounts.Register("The test pass",
		[]types.String{"Test permission 1", "Test permission 2"})

	if regErr != nil {
		t.Errorf("Accounts.Register() failed to register account due to error: %s",
			regErr)

		return
	}

	getAcc, getErr := accounts.Get("The test pass")

	if getErr != nil {
		t.Errorf("Accounts.Get() failed to get the account due to error: %s",
			getErr)

		return
	}

	if getAcc != regAcc {
		t.Error("Accounts.Get() failed to get the correct Account")

		return
	}

	// Try get an account which is not existed
	getAcc, getErr = accounts.Get("An account which not existed")

	if getErr == nil || !getErr.Is(ErrAccountNotFound) {
		t.Error("Accounts.Get() got an account which is not existed, impossible")

		return
	}

	if getAcc != nil {
		t.Error("Accounts.Get() mistakenly filled the Account result")

		return
	}
}

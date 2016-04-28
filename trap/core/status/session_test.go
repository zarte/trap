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

	"net"
	"testing"
	"time"
)

func getEmptyAccount() *Account {
	account := &Account{
		permission: Permission{},
	}

	return account
}

func getTestSession() *Session {
	now := time.Now()

	return &Session{
		IP:      net.ParseIP("127.0.0.1"),
		Created: now,
		Key:     "TEST_RANDOM_KEY",

		LastSeen: now,
		Expire:   time.Duration(2) * time.Second,

		account: getEmptyAccount(),
	}
}

func TestSessionBump(t *testing.T) {
	session := getTestSession()

	time.Sleep(1 * time.Second)

	session.Bump()

	if !session.LastSeen.After(session.Created) {
		t.Error("Session.Bump() failed up renew the LastSeen field")

		return
	}
}

func TestSessionAccount(t *testing.T) {
	session := getTestSession()

	if session.Account() == nil {
		t.Error("Session.Account() failed to pass the Account information")

		return
	}
}

func TestSessionExpired(t *testing.T) {
	session := getTestSession()

	if session.Expired() {
		t.Error("Session can't be expired right now")

		return
	}

	time.Sleep(3 * time.Second)

	if !session.Expired() {
		t.Error("Session should be expired right now")

		return
	}
}

func TestSessionsGetRandomKey(t *testing.T) {
	maxKeys := 1000
	randonKeyMap := map[types.String]bool{}

	sessions := Sessions{}

	for i := 0; i < maxKeys; i++ {
		randomKeys := sessions.getRandomKey()

		if _, ok := randonKeyMap[randomKeys]; !ok {
			randonKeyMap[randomKeys] = true

			continue
		}

		t.Error("Random key conflicted even in a small test")

		return
	}
}

func TestSessionsAdd(t *testing.T) {
	sessions := Sessions{}

	newSession, newSessErr := sessions.Add(net.ParseIP("127.0.0.1"),
		getEmptyAccount(), 12*time.Second)

	if newSessErr != nil {
		t.Errorf("Failed to create session due to error: %s", newSessErr)

		return
	}

	if !newSession.IP.Equal(net.ParseIP("127.0.0.1")) {
		t.Error("Failed to assign session IP address")

		return
	}
}

func TestSessionsDelete(t *testing.T) {
	sessions := Sessions{}

	newSession, newSessErr := sessions.Add(net.ParseIP("127.0.0.1"),
		getEmptyAccount(), 12*time.Second)

	if newSessErr != nil {
		t.Errorf("Failed to create dummy session due to error: %s", newSessErr)

		return
	}

	deleteErr := sessions.Delete(newSession.Key)

	if deleteErr != nil {
		t.Errorf("Failed to delete session due to error: %s", deleteErr)

		return
	}

	deleteErr = sessions.Delete("NOT_EXISTED_SESSION")

	if deleteErr == nil || !deleteErr.Is(ErrSessionKeyNotFound) {
		t.Error("Expected error does not happen")

		return
	}
}

func TestSessionsVerify(t *testing.T) {
	sessions := Sessions{}

	newSession, newSessErr := sessions.Add(net.ParseIP("127.0.0.1"),
		getEmptyAccount(), 12*time.Second)

	if newSessErr != nil {
		t.Errorf("Failed to create dummy session due to error: %s", newSessErr)

		return
	}

	verifySess, verifyErr := sessions.Verify(net.ParseIP("127.0.0.1"),
		newSession.Key)

	if verifyErr != nil {
		t.Errorf("Failed to verify session due to error: %s", verifyErr)

		return
	}

	if verifySess == nil {
		t.Error("Failed to output session information")

		return
	}

	// Test not found
	verifySess, verifyErr = sessions.Verify(net.ParseIP("127.0.0.2"),
		newSession.Key)

	if verifyErr == nil || !verifyErr.Is(ErrSessionNotFound) {
		t.Error("Expected error does not happen")

		return
	}

	if verifySess != nil {
		t.Error("Session information should not be exported here")

		return
	}

	// Test expired
	newSession.LastSeen = time.Now().Add(-(13 * time.Hour))

	verifySess, verifyErr = sessions.Verify(net.ParseIP("127.0.0.1"),
		newSession.Key)

	if verifySess != nil {
		t.Error("Session information should not be exported here")

		return
	}

	if verifyErr == nil || !verifyErr.Is(ErrSessionExpired) {
		t.Error("Expected error does not happen")

		return
	}
}

func TestSessionsDump(t *testing.T) {
	sessions := Sessions{}

	sessions.Add(net.ParseIP("127.0.0.1"), getEmptyAccount(), 12*time.Second)
	sessions.Add(net.ParseIP("127.0.0.2"), getEmptyAccount(), 12*time.Second)
	sessions.Add(net.ParseIP("127.0.0.3"), getEmptyAccount(), 12*time.Second)

	if len(sessions.Dump()) != 3 {
		t.Error("Dumpped an invalid amount of sessions")

		return
	}
}

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

package types

import (
    "testing"
)

func TestSecretBytes(t *testing.T) {
    s := Secret("Data for test here")

    if string(s.Bytes()) != "Data for test here" {
        t.Error("Secret.Bytes() can't convert bytes back to byte slice correctly")
    }
}

func TestSecretString(t *testing.T) {
    s := Secret("Data for test here")

    if s.String() != "Data for test here" {
        t.Error("Secret.String() can't convert Secret to string correctly")
    }
}

func TestSecretLen(t *testing.T) {
    s := Secret("Data for test here")

    if s.Len() != 18 {
        t.Error("Secret.Len() can't get Secret length correctly")
    }
}

func TestSecretSHA256(t *testing.T) {
    s := Secret("Data for test here")
    v := Secret{
        '\x9c', '\xef', '\xd8', '\x5b', '\x0f', '\x44', '\x26', '\x6a', '\x1c',
        '\xfb', '\x2b', '\xa5', '\x7c', '\xac', '\x98', '\x27', '\x29', '\xca',
        '\xe5', '\xb1', '\x4d', '\xd5', '\xd6', '\x5b', '\xc0', '\x8d', '\x90',
        '\x61', '\xb7', '\xc2', '\x5e', '\xf6',
    }

    if !s.SHA256().IsEqual(&v) {
        t.Errorf("Secret can't sum SHA256 correctly, excepting '%d', got '%d'",
            v, s.SHA256())
    }
}

func TestSecretSHA512(t *testing.T) {
    s := Secret("Data for test here")
    v := Secret{
        '\x11', '\x13', '\xED', '\x95', '\xD7', '\xEE', '\x06', '\x9D', '\x00',
        '\x19', '\xFE', '\x80', '\x5B', '\xEF', '\x7B', '\x31', '\xFF', '\x2A',
        '\x53', '\x41', '\x58', '\x29', '\x22', '\xA4', '\xE9', '\xAA', '\xF0',
        '\xD2', '\x9D', '\xD3', '\x8C', '\x7B', '\x81', '\x3D', '\x41', '\xAD',
        '\xC3', '\xD2', '\x96', '\x07', '\xE9', '\x19', '\x67', '\x9F', '\x1E',
        '\x1A', '\xAC', '\x53', '\x71', '\x05', '\x01', '\x80', '\x0E', '\x77',
        '\x8A', '\x9F', '\xFC', '\xB7', '\x90', '\x7C', '\xD5', '\x84', '\xCE',
        '\x1B',
    }

    if !s.SHA512().IsEqual(&v) {
        t.Errorf("Secret can't sum SHA512 correctly, excepting '%d', got '%d'",
            v, s.SHA512())
    }
}

func TestSecretEncryptDecrypt(t *testing.T) {
    passwd  := Secret("This is the password")
    passwd2 := Secret("This is another password")
    raw     := Secret("This is some data which will be encrypted")

    encrypted, enErr    := raw.Encrypt(passwd)
    encrypted2, enErr2  := raw.Encrypt(passwd)

    if enErr != nil {
        t.Errorf("Secret can't encrypt data due to error: %s", enErr)

        return
    }

    decrypted, deErr    := encrypted.Decrypt(passwd)
    decrypted2, deErr2  := encrypted.Decrypt(passwd2)

    if deErr != nil {
        t.Errorf("Secret can't decrypt data due to error: %s", deErr)

        return
    }

    if !decrypted.IsEqual(&raw) {
        t.Error("Secret can't decrypt data which encrypted by itself with same passcode")

        return
    }

    if enErr2 != nil {
        t.Errorf("Secret can't encrypt data due to error: %s", enErr2)

        return
    }

    if encrypted2.IsEqual(&encrypted) {
        t.Error("Two encrypted data is equal, which is impossible as we are using random IV")

        return
    }

    if deErr2 != nil {
        t.Errorf("Secret can't decrypt data due to error: %s", deErr2)

        return
    }

    if encrypted2.IsEqual(&encrypted) {
        t.Error("Two encrypted data is equal, which is impossible as we are using random IV")

        return
    }

    if decrypted2.IsEqual(&raw) {
        t.Error("This is impossible to decrypt an encrypted data with wrong passcode")

        return
    }
}
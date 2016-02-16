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
    "crypto/rand"
    "crypto/sha512"
    "crypto/sha256"
    "crypto/aes"
    "crypto/cipher"
)

type Secret []byte

var (
    ErrSecretDecryptionContentTooShort =
        NewError("Content is too short to be decrypted")
)

func (s Secret) Byte() ([]byte) {
    return s
}

func (s Secret) SHA256() (Secret) {
    sum := sha256.Sum256(s)

    return Secret(sum[:])
}

func (s Secret) SHA512() (Secret) {
    sum := sha512.Sum512(s)

    return Secret(sum[:])
}

/**
 * Following two functions is largely steal from:
 * https://gist.github.com/josephspurrier/8304f09562d81babb494
 *
 * Thank you Joseph Spurrier!
 */
func (s Secret) Encrypt(key Secret) (Secret, *Throw) {
    aesBlock, blockErr  := aes.NewCipher(key.SHA256())

    if blockErr != nil {
        return Secret(""), ConvertError(blockErr)
    }

    encrypted           := make([]byte, aes.BlockSize + len(s))

    // Refer the first `aes.BlockSize` bytes of `encrypted`, magic!
    iv                  := encrypted[:aes.BlockSize]

    _, rErr             := rand.Read(iv)

    if rErr != nil {
        return Secret(""), ConvertError(rErr)
    }

    stream := cipher.NewCFBEncrypter(aesBlock, iv)

    stream.XORKeyStream(encrypted[aes.BlockSize:], s)

    return Secret(encrypted), nil
}

func (s Secret) Decrypt(key Secret) (Secret, *Throw) {
    aesBlock, blockErr  := aes.NewCipher(key.SHA256())

    if blockErr != nil {
        return Secret(""), ConvertError(blockErr)
    }

    if len(s) < aes.BlockSize {
        return Secret(""), ErrSecretDecryptionContentTooShort.Throw()
    }

    iv                  := s[:aes.BlockSize]
    withoutIV           := s[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(aesBlock, iv)

    stream.XORKeyStream(withoutIV, withoutIV)

    return Secret(withoutIV), nil
}
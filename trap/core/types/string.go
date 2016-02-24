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
    "strings"
    "strconv"
    "bytes"
)

type String string

func (s String) String() (string) {
    return string(s)
}

func (s String) Bytes() ([]byte) {
    return []byte(s)
}

func (s String) Len() (int) {
    return len(s.String())
}

func (s String) Join(joiningStrs ...String) (String) {
    var buffer bytes.Buffer

    buffer.WriteString(s.String())

    for _, joiningStr := range joiningStrs {
        buffer.WriteString(joiningStr.String())
    }

    return String(buffer.String())
}

func (s String) Contains(contain string) (bool) {
    return strings.Contains(s.String(), contain)
}

func (s String) Cut(start, end int) (String) {
    result      :=  ""
    str         :=  s.String()
    strLen      :=  len(str)
    strLastIdx  :=  strLen - 1

    if strLen <= 0 {
        return String("")
    }

    if start > strLastIdx || end > strLastIdx {
        return String("")
    }

    if end < 0 || start < 0 {
        return String("")
    }

    if start < end {
        for i := start; i <= end; i++ {
            result += string(str[i])
        }
    } else {
        for i := start; i >= end; i-- {
            result += string(str[i])
        }
    }

    return String(result)
}

func (s String) SelectHead(selLen int) (String, String) {
    strLastIdx  :=  s.Len() - 1
    strEndIdx   :=  selLen - 1
    strTail     :=  String("")

    if strLastIdx < 0 {
        return String(""), String("")
    }

    if strEndIdx > strLastIdx {
        strEndIdx = strLastIdx
    }

    if strLastIdx > strEndIdx {
        strTail   =   s.Cut(strEndIdx + 1, strLastIdx)
    }

    return s.Cut(0, strEndIdx), strTail
}

func (s String) SelectTail(selLen int) (String, String) {
    strHead     :=  String("")
    strLen      :=  s.Len()

    strEndIdx   :=  strLen - 1

    if strEndIdx < 0 {
        return String(""), String("")
    }

    strBeginIdx :=  strEndIdx - (selLen - 1)

    if strBeginIdx < 0 {
        strBeginIdx = 0
    }

    if strBeginIdx > 0 {
        strHead =   s.Cut(0, strBeginIdx - 1)
    }

    return strHead, s.Cut(strBeginIdx, strEndIdx)
}

func (s String) SpiltWith(spiltter string) (String, String) {
    pureString := s.String()

    if !strings.Contains(pureString, spiltter) {
        return s, ""
    }

    spIndex := strings.Index(pureString, spiltter)

    return String(pureString[0:spIndex]),
        String(pureString[spIndex + 1:len(pureString)])
}

func (s String) Trim() (String) {
    return String(strings.TrimSpace(s.String()))
}

func (s String) Int16() (Int16) {
    resultUint, err := strconv.ParseInt(s.String(), 10, 16)

    if err != nil {
        return 0
    }

    return Int16(resultUint)
}

func (s String) Int32() (Int32) {
    resultUint, err := strconv.ParseInt(s.String(), 10, 32)

    if err != nil {
        return 0
    }

    return Int32(resultUint)
}

func (s String) Int64() (Int64) {
    resultUint, err := strconv.ParseInt(s.String(), 10, 64)

    if err != nil {
        return 0
    }

    return Int64(resultUint)
}


func (s String) UInt16() (UInt16) {
    resultUint, err := strconv.ParseUint(s.String(), 10, 16)

    if err != nil {
        return 0
    }

    return UInt16(resultUint)
}

func (s String) UInt32() (UInt32) {
    resultUint, err := strconv.ParseUint(s.String(), 10, 32)

    if err != nil {
        return 0
    }

    return UInt32(resultUint)
}

func (s String) UInt64() (UInt64) {
    resultUint, err := strconv.ParseUint(s.String(), 10, 64)

    if err != nil {
        return 0
    }

    return UInt64(resultUint)
}

func (s String) Lower() (String) {
    return String(strings.ToLower(s.String()))
}

func (s String) Upper() (String) {
    return String(strings.ToUpper(s.String()))
}

func (s String) ExplodeWith(exploder string) ([]String) {
    var stringSlice []String

    for _, str := range strings.Split(s.String(), exploder) {
        stringSlice = append(stringSlice, String(str))
    }

    return stringSlice
}

func (s String) Replace(find String, replace String) (String) {
    f := string(find)
    t := string(replace)

    return String(strings.Replace(string(s), f, t, -1))
}
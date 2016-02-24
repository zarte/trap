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

func TestStringToRealString(t *testing.T) {
    str := String("Test string")
    compareStr := "Test string"

    if str.String() != compareStr {
        t.Error("Can't convert the String to string correctly")
    }
}

func TestStringToBytes(t *testing.T) {
    str := String("Test string")

    if string(str.Bytes()) != "Test string" {
        t.Error("Can't convert the String to byte slice correctly")
    }
}

func TestStringLen(t *testing.T) {
    str := String("this is a test string")

    if str.Len() != len("this is a test string") {
        t.Error("Can't calc the right length of the String")
    }
}

func TestStringJoin(t *testing.T) {
    str := String("this is a test string").Join("This is ", "Another", "String")

    if str != "this is a test stringThis is AnotherString" {
        t.Error("Can't join String to String")
    }
}

func TestStringContains(t *testing.T) {
    str := String("Hello World!")

    if !str.Contains("!") {
        t.Error("String.Contains() has failed to found the substr")

        return
    }

    if !str.Contains("Hello") {
        t.Error("String.Contains() has failed to found the substr")

        return
    }

    if !str.Contains("World") {
        t.Error("String.Contains() has failed to found the substr")

        return
    }

    if !str.Contains(" ") {
        t.Error("String.Contains() has failed to found the substr")

        return
    }

    if str.Contains("?") {
        t.Error("String.Contains() found a non-exist substr which is impossible")

        return
    }
}

func TestStringCut(t *testing.T) {
    str := String("This is a string example")

    if str.Cut(0, 0) != "T" {
        t.Error("String.Cut() failed to cut out the substr correctly")

        return
    }

    if str.Cut(0, 500) != "" {
        t.Error("String.Cut() failed to cut out the substr correctly")

        return
    }

    if str.Cut(500, 0) != "" {
        t.Error("String.Cut() failed to cut out the substr correctly")

        return
    }

    if str.Cut(500, 500) != "" {
        t.Error("String.Cut() failed to cut out the substr correctly")

        return
    }

    if str.Cut(0, 3) != "This" {
        t.Error("String.Cut() failed to cut out the substr correctly")

        return
    }

    if str.Cut(5, 6) != "is" {
        t.Error("String.Cut() failed to cut out the substr correctly")

        return
    }

    if str.Cut(6, 5) != "si" {
        t.Error("String.Cut() failed to cut out the substr correctly")

        return
    }

    if str.Cut(3, 0) != "sihT" {
        t.Error("String.Cut() failed to cut out the substr correctly")

        return
    }
}

func TestStringSelectHead(t *testing.T) {
    part1, part2 := String("This is a string example").SelectHead(10)

    if part1 != "This is a " || part2 != "string example" {
        t.Error("String.SelectHead() failed to select string correctly")

        return
    }

    part1, part2 =  String("This is a string example").SelectHead(0)

    if part1 != "" || part2 != "This is a string example" {
        t.Error("String.SelectHead() failed to select string correctly")

        return
    }

    part1, part2 =  String("This is a string example").SelectHead(100)

    if part1 != "This is a string example" || part2 != "" {
        t.Error("String.SelectHead() failed to select string correctly")

        return
    }
}

func TestStringSelectTail(t *testing.T) {
    part1, part2 := String("This is a string example").SelectTail(10)

    if part1 != "This is a stri" || part2 != "ng example" {
        t.Error("String.SelectTail() failed to select string correctly")

        return
    }

    part1, part2  = String("This is a string example").SelectTail(0)

    if part1 != "This is a string example" || part2 != "" {
        t.Error("String.SelectTail() failed to select string correctly")

        return
    }

    part1, part2  = String("This is a string example").SelectTail(100)

    if part1 != "" || part2 != "This is a string example" {
        t.Error("String.SelectTail() failed to select string correctly")

        return
    }

    part1, part2  = String("This is a string example").SelectTail(-100)

    if part1 != "" || part2 != "" {
        t.Error("String.SelectTail() failed to select string correctly")

        return
    }
}

func TestStringSpiltWith(t *testing.T) {
    part1, part2 := String("This is a test string").SpiltWith(" ")

    if part1 != "This" || part2 != "is a test string" {
        t.Error("String.SpiltWith() failed to spilt string correctly")

        return
    }

    part1, part2  = part2.SpiltWith("!")

    if part1 != "is a test string" || part2 != "" {
        t.Error("String.SpiltWith() failed to spilt string correctly")

        return
    }
}

func TestStringTrim(t *testing.T) {
    trimedString := String("This string has been trimmed  ").Trim()

    if trimedString != "This string has been trimmed" {
        t.Error("String.Trim() failed to trim string correctly")

        return
    }

    trimedString =  String("  This string has been trimmed  ").Trim()

    if trimedString != "This string has been trimmed" {
        t.Error("String.Trim() failed to trim string correctly")

        return
    }

    trimedString =  String("  This string has been trimmed").Trim()

    if trimedString != "This string has been trimmed" {
        t.Error("String.Trim() failed to trim string correctly")

        return
    }
}

func TestStringInt16(t *testing.T) {
    max         := Int16(MAX_INT16)
    min         := Int16(MAX_INT16)
    maxConvert  := String(max.String()).Int16()
    minConvert  := String(min.String()).Int16()

    if max != maxConvert {
        t.Error("String.Int16() failed to convert string to integer")

        return
    }

    if min != minConvert {
        t.Error("String.Int16() failed to convert string to integer")

        return
    }
}

func TestStringInt32(t *testing.T) {
    max         := Int32(MAX_INT32)
    min         := Int32(MAX_INT32)
    maxConvert  := String(max.String()).Int32()
    minConvert  := String(min.String()).Int32()

    if max != maxConvert {
        t.Error("String.Int32() failed to convert string to integer")

        return
    }

    if min != minConvert {
        t.Error("String.Int32() failed to convert string to integer")

        return
    }
}

func TestStringInt64(t *testing.T) {
    max         := Int64(MAX_INT32)
    min         := Int64(MAX_INT32)
    maxConvert  := String(max.String()).Int64()
    minConvert  := String(min.String()).Int64()

    if max != maxConvert {
        t.Error("String.Int64() failed to convert string to integer")

        return
    }

    if min != minConvert {
        t.Error("String.Int64() failed to convert string to integer")

        return
    }
}

func TestStringUInt16(t *testing.T) {
    max         := UInt16(MAX_UINT16)
    min         := UInt16(MAX_UINT16)
    maxConvert  := String(max.String()).UInt16()
    minConvert  := String(min.String()).UInt16()

    if max != maxConvert {
        t.Error("String.UInt16() failed to convert string to integer")

        return
    }

    if min != minConvert {
        t.Error("String.UInt16() failed to convert string to integer")

        return
    }
}

func TestStringUInt32(t *testing.T) {
    max         := UInt32(MAX_UINT32)
    min         := UInt32(MAX_UINT32)
    maxConvert  := String(max.String()).UInt32()
    minConvert  := String(min.String()).UInt32()

    if max != maxConvert {
        t.Error("String.UInt32() failed to convert string to integer")

        return
    }

    if min != minConvert {
        t.Error("String.UInt32() failed to convert string to integer")

        return
    }
}

func TestStringUInt64(t *testing.T) {
    max         := UInt64(MAX_UINT64)
    min         := UInt64(MAX_UINT64)
    maxConvert  := String(max.String()).UInt64()
    minConvert  := String(min.String()).UInt64()

    if max != maxConvert {
        t.Error("String.UInt64() failed to convert string to integer")

        return
    }

    if min != minConvert {
        t.Error("String.UInt64() failed to convert string to integer")

        return
    }
}

func TestStringLower(t *testing.T) {
    str := String("this Is a TEST STRINg")

    if str.Lower() != "this is a test string" {
        t.Error("String.Lower() failed to lower the test string")

        return
    }
}

func TestStringUpper(t *testing.T) {
    str := String("this Is a TEST STRINg")

    if str.Upper() != "THIS IS A TEST STRING" {
        t.Error("String.Upper() failed to upper the test string")

        return
    }
}

func TestStringExplodeWith(t *testing.T) {
    str := String("this Is a TEST STRINg")

    exploded := str.ExplodeWith(" ")

    if exploded[0] != "this" ||
    exploded[1] != "Is" ||
    exploded[2] != "a" ||
    exploded[3] != "TEST" ||
    exploded[4] != "STRINg" {
        t.Error("String.ExplodeWith() failed to explode the test string")

        return
    }

    // Two spaces here
    str = String("this  Is a TEST STRINg")

    exploded = str.ExplodeWith(" ")

    if exploded[0] != "this" ||
    exploded[1] != "" ||
    exploded[2] != "Is" ||
    exploded[3] != "a" ||
    exploded[4] != "TEST" ||
    exploded[5] != "STRINg" {
        t.Error("String.ExplodeWith() failed to explode the test string")

        return
    }
}

func TestStringReplace(t *testing.T) {
    simple := String("This is a test xample Xample")

    replaced := simple.Replace("xample", "example")

    if replaced != "This is a test example Xample" {
        t.Error("String.Replace() failed to replace the test string")

        return
    }
}
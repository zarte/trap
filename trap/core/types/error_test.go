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
	"errors"
	"fmt"
	"testing"
)

func TestErrorThrow(t *testing.T) {
	err := NewError("This is a test error: %s %s %s")
	throw := err.Throw("some", "error", "information")

	if throw.parentRef != err.ref ||
		throw.parentErr != err.error ||
		throw.message != fmt.Sprintf(err.text, "some", "error", "information") {
		t.Error("Error can't build correct Throw info")
	}
}

func TestThrowIs(t *testing.T) {
	err1 := NewError("Test error 1")
	err2 := NewError("Test error 2")

	if !err1.Throw().Is(err1) || err1.Throw().Is(err2) {
		t.Error("Throw failed compare it's parent error")
	}
}

func TestThrowIsError(t *testing.T) {
	err1 := errors.New("Test error 1")
	err2 := errors.New("Test error 2")

	thr1 := ConvertError(err1)

	if !thr1.IsError(err1) || thr1.IsError(err2) {
		t.Error("Throw failed compare it's parent error")
	}
}

func TestThrowError(t *testing.T) {
	throw := NewError("Test error 1").Throw()

	if throw.Error() != "Test error 1" {
		t.Error("Error message is incorrect")
	}

	throw = NewError("Test error 1 %s %s %s").
		Throw("formated", "strings", "here")

	if throw.Error() != "Test error 1 formated strings here" {
		t.Error("Error message is incorrect")
	}
}

func TestThrowMarshalJSON(t *testing.T) {
	json, jsonErr := NewError("Test error 1").Throw().MarshalJSON()

	if jsonErr != nil {
		t.Error("Failed marshal Throw information to JSON")

		return
	}

	if string(json) != "\"Test error 1\"" {
		t.Error("Throw generated an invalid JSON result")

		return
	}

	json, jsonErr = NewError("Test error 1 \"%s\" \"%s\" \"%s\"").
		Throw("formated", "strings", "here").MarshalJSON()

	if jsonErr != nil {
		t.Error("Failed marshal Throw information to JSON")

		return
	}

	if string(json) != "\"Test error 1 \\\"formated\\\" \\\"strings\\\" \\\"here\\\"\"" {
		t.Error("Throw generated an invalid JSON result")

		return
	}
}

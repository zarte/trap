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

package event

import (
	"github.com/raincious/trap/trap/core/logger"
	"github.com/raincious/trap/trap/core/types"

	"testing"
)

var (
	ErrFakeFail *types.Error = types.NewError("This is a fake error")
)

func getEmptyLogger() *logger.Logger {
	return logger.NewLogger()
}

func getEmptyEvent() *Event {
	e := Event{}

	e.Init(&Config{
		Logger: getEmptyLogger(),
	})

	return &e
}

func TestEventInit(t *testing.T) {
	log := getEmptyLogger()

	e := Event{}

	e.Init(&Config{
		Logger: log,
	})

	if e.error != nil {
		t.Error("`Event` error is impossible to be un-nil")

		return
	}

	if e.events == nil {
		t.Error("`Event` forgot init it's events map")

		return
	}
}

func TestEventRegisterNTrigger(t *testing.T) {
	emptyParam := Parameters{}
	filledParam := Parameters{}.AddString("VALUE", "String").AddString("VALUE2", "String2")
	filledParam2 := Parameters{}.AddString("VALUE3", "String3")
	resultString := types.String("")

	e := getEmptyEvent()

	// Test normal
	e.Register("test.succeed.callback", func(params *Parameters) *types.Throw {
		resultString = params.Parse("This $((VALUE)) $((VALUE)) $((VALUE2)) $((VALUEX)) is parsed", []types.String{
			"$((VALUE))",
			"$((VALUE2))",
		})

		return nil
	})

	e.Register("test.failed.callback", func(params *Parameters) *types.Throw {
		return ErrFakeFail.Throw()
	})

	tError1 := e.Trigger("test.succeed.callback", filledParam)

	if tError1 != nil {
		t.Errorf("Can't trigger test event due to error: %s", tError1)

		return
	}

	if resultString != "This String String String2 $((VALUEX)) is parsed" {
		t.Error("Parsed string is not correct")

		return
	}

	// Test failed callback
	tError2 := e.Trigger("test.failed.callback", emptyParam)

	if tError2 == nil || !tError2.Is(ErrFakeFail) {
		t.Errorf("Can't trigger test event due to error: %s", tError2)

		return
	}

	// Test unfilled Parameters in a callback that parses format
	tError3 := e.Trigger("test.succeed.callback", emptyParam)

	if tError3 != nil {
		t.Errorf("Can't trigger test event due to error: %s", tError3)

		return
	}

	if resultString != "This    $((VALUEX)) is parsed" {
		t.Error("Parsed string is not correct")

		return
	}

	// Test a mistakenly filled Parameters in a callback that parses format
	tError4 := e.Trigger("test.succeed.callback", filledParam2)

	if tError4 != nil {
		t.Errorf("Can't trigger test event due to error: %s", tError4)

		return
	}

	if resultString != "This    $((VALUEX)) is parsed" {
		t.Error("Parsed string is not correct")

		return
	}
}

func TestEventRegisterNTriggerMultiCallbacks(t *testing.T) {
	filledParam := Parameters{}.
		AddString("STRING", "String")

	e := getEmptyEvent()

	results := []types.String{}

	e.Register("test.callback", func(params *Parameters) *types.Throw {
		return ErrFakeFail.Throw()
	})

	e.Register("test.callback", func(params *Parameters) *types.Throw {
		results = append(results,
			params.Parse("This is a result $((STRING))", []types.String{"$((STRING))"}))

		return nil
	})

	e.Register("test.callback", func(params *Parameters) *types.Throw {
		results = append(results,
			params.Parse("This is a result $((STRING)) 2", []types.String{"$((STRING))"}))

		return nil
	})

	tError := e.Trigger("test.callback", filledParam)

	if tError == nil || !tError.Is(ErrFakeFail) {
		t.Errorf("Can't trigger test event due to error: %s", tError)

		return
	}

	if len(results) != 2 {
		t.Errorf("Excepting '2' results, got '%d'", len(results))

		return
	}

	if results[0] != "This is a result String" {
		t.Error("Callback result is invalid")

		return
	}

	if results[1] != "This is a result String 2" {
		t.Error("Callback result is invalid")

		return
	}
}

func TestEventTriggerNonExisted(t *testing.T) {
	e := getEmptyEvent()

	tError := e.Trigger("test.undefined", Parameters{})

	if tError == nil || !tError.Is(ErrNoEvent) {
		t.Errorf("Unexpected error: %s", tError)

		return
	}
}

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

package logger

import (
	"github.com/raincious/trap/trap/core/types"

	"testing"
	"time"
)

type fakeLogPrinterForPrinterTest struct {
	debugHasBeenCalled   bool
	infoHasBeenCalled    bool
	warningHasBeenCalled bool
	errorHasBeenCalled   bool
	printHasBeenCalled   bool
}

func (f *fakeLogPrinterForPrinterTest) Info(
	c types.String, t time.Time, m types.String) {
	f.infoHasBeenCalled = true
}

func (f *fakeLogPrinterForPrinterTest) Debug(
	c types.String, t time.Time, m types.String) {
	f.debugHasBeenCalled = true
}

func (f *fakeLogPrinterForPrinterTest) Warning(
	c types.String, t time.Time, m types.String) {
	f.warningHasBeenCalled = true
}

func (f *fakeLogPrinterForPrinterTest) Error(
	c types.String, t time.Time, m types.String) {
	f.errorHasBeenCalled = true
}

func (f *fakeLogPrinterForPrinterTest) Print(
	c types.String, t time.Time, m types.String) {
	f.printHasBeenCalled = true
}

func TestPrintersAdd(t *testing.T) {
	printers := Printers{}
	fakePrinter1 := &fakeLogPrinterForPrinterTest{
		debugHasBeenCalled:   false,
		infoHasBeenCalled:    false,
		warningHasBeenCalled: false,
		errorHasBeenCalled:   false,
		printHasBeenCalled:   false,
	}
	fakePrinter2 := &fakeLogPrinterForPrinterTest{
		debugHasBeenCalled:   false,
		infoHasBeenCalled:    false,
		warningHasBeenCalled: false,
		errorHasBeenCalled:   false,
		printHasBeenCalled:   false,
	}

	if len(printers) != 0 {
		t.Error("Unexpected printers size")

		return
	}

	printers.Add(fakePrinter1)

	if len(printers) != 1 {
		t.Error("Unexpected printers size")

		return
	}

	printers.Add(fakePrinter2)

	if len(printers) != 2 {
		t.Error("Unexpected printers size")

		return
	}
}

func getNewTestPrinters() (
	Printers, *fakeLogPrinterForPrinterTest, *fakeLogPrinterForPrinterTest) {
	printers := Printers{}
	fakePrinter1 := &fakeLogPrinterForPrinterTest{
		debugHasBeenCalled:   false,
		infoHasBeenCalled:    false,
		warningHasBeenCalled: false,
		errorHasBeenCalled:   false,
		printHasBeenCalled:   false,
	}
	fakePrinter2 := &fakeLogPrinterForPrinterTest{
		debugHasBeenCalled:   false,
		infoHasBeenCalled:    false,
		warningHasBeenCalled: false,
		errorHasBeenCalled:   false,
		printHasBeenCalled:   false,
	}

	printers.Add(fakePrinter1)
	printers.Add(fakePrinter2)

	return printers, fakePrinter1, fakePrinter2
}

func TestPrintersDebug(t *testing.T) {
	printers, fake1, fake2 := getNewTestPrinters()

	fake3 := &fakeLogPrinterForPrinterTest{
		debugHasBeenCalled:   false,
		infoHasBeenCalled:    false,
		warningHasBeenCalled: false,
		errorHasBeenCalled:   false,
		printHasBeenCalled:   false,
	}

	if fake1.debugHasBeenCalled || fake2.debugHasBeenCalled ||
		fake1.infoHasBeenCalled || fake2.infoHasBeenCalled ||
		fake1.warningHasBeenCalled || fake2.warningHasBeenCalled ||
		fake1.errorHasBeenCalled || fake2.errorHasBeenCalled ||
		fake1.printHasBeenCalled || fake2.printHasBeenCalled {
		t.Error("Unexpected initial data")

		return
	}

	// Debug
	printers.Debug("Context title", time.Now(), "Test message")

	if !fake1.debugHasBeenCalled || !fake2.debugHasBeenCalled || fake3.debugHasBeenCalled {
		t.Error("Failed asserting Printers.Debug() method is called correctly")

		return
	}

	if fake1.infoHasBeenCalled || fake2.infoHasBeenCalled ||
		fake1.warningHasBeenCalled || fake2.warningHasBeenCalled ||
		fake1.errorHasBeenCalled || fake2.errorHasBeenCalled ||
		fake1.printHasBeenCalled || fake2.printHasBeenCalled {
		t.Error("Printers.Debug() is poisoning other test data")

		return
	}

	// Info
	printers, fake1, fake2 = getNewTestPrinters()

	printers.Info("Context title", time.Now(), "Test message")

	if !fake1.infoHasBeenCalled || !fake2.infoHasBeenCalled || fake3.infoHasBeenCalled {
		t.Error("Failed asserting Printers.Info() method is called correctly")

		return
	}

	if fake1.debugHasBeenCalled || fake2.debugHasBeenCalled ||
		fake1.warningHasBeenCalled || fake2.warningHasBeenCalled ||
		fake1.errorHasBeenCalled || fake2.errorHasBeenCalled ||
		fake1.printHasBeenCalled || fake2.printHasBeenCalled {
		t.Error("Printers.Info() is poisoning other test data")

		return
	}

	// Warning
	printers, fake1, fake2 = getNewTestPrinters()

	printers.Warning("Context title", time.Now(), "Test message")

	if !fake1.warningHasBeenCalled || !fake2.warningHasBeenCalled || fake3.warningHasBeenCalled {
		t.Error("Failed asserting Printers.Warning() method is called correctly")

		return
	}

	if fake1.debugHasBeenCalled || fake2.debugHasBeenCalled ||
		fake1.infoHasBeenCalled || fake2.infoHasBeenCalled ||
		fake1.errorHasBeenCalled || fake2.errorHasBeenCalled ||
		fake1.printHasBeenCalled || fake2.printHasBeenCalled {
		t.Error("Printers.Warning() is poisoning other test data")

		return
	}

	// Error
	printers, fake1, fake2 = getNewTestPrinters()

	printers.Error("Context title", time.Now(), "Test message")

	if !fake1.errorHasBeenCalled || !fake2.errorHasBeenCalled || fake3.errorHasBeenCalled {
		t.Error("Failed asserting Printers.Error() method is called correctly")

		return
	}

	if fake1.debugHasBeenCalled || fake2.debugHasBeenCalled ||
		fake1.infoHasBeenCalled || fake2.infoHasBeenCalled ||
		fake1.warningHasBeenCalled || fake2.warningHasBeenCalled ||
		fake1.printHasBeenCalled || fake2.printHasBeenCalled {
		t.Error("Printers.Error() is poisoning other test data")

		return
	}

	// Print
	printers, fake1, fake2 = getNewTestPrinters()

	printers.Print("Context title", time.Now(), "Test message")

	if !fake1.printHasBeenCalled || !fake2.printHasBeenCalled || fake3.printHasBeenCalled {
		t.Error("Failed asserting Printers.Print() method is called correctly")

		return
	}

	if fake1.debugHasBeenCalled || fake2.debugHasBeenCalled ||
		fake1.infoHasBeenCalled || fake2.infoHasBeenCalled ||
		fake1.warningHasBeenCalled || fake2.warningHasBeenCalled ||
		fake1.errorHasBeenCalled || fake2.errorHasBeenCalled {
		t.Error("Printers.Print() is poisoning other test data")

		return
	}

	// Real time register a new printer and apply a call
	printers.Add(fake3)

	if fake3.debugHasBeenCalled || fake3.infoHasBeenCalled ||
		fake3.warningHasBeenCalled || fake3.errorHasBeenCalled {
		t.Error("Initial data not ready")

		return
	}

	printers.Print("Context title", time.Now(), "Test message")

	if !fake1.printHasBeenCalled || !fake2.printHasBeenCalled || !fake3.printHasBeenCalled {
		t.Error("Failed asserting Printers.Print() method is called correctly")

		return
	}

	if fake1.debugHasBeenCalled || fake2.debugHasBeenCalled || fake3.debugHasBeenCalled ||
		fake1.infoHasBeenCalled || fake2.infoHasBeenCalled || fake3.infoHasBeenCalled ||
		fake1.warningHasBeenCalled || fake2.warningHasBeenCalled || fake3.warningHasBeenCalled ||
		fake1.errorHasBeenCalled || fake2.errorHasBeenCalled || fake3.errorHasBeenCalled {
		t.Error("Printers.Print() is poisoning other test data")

		return
	}
}

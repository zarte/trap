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

package messager

import (
	"github.com/raincious/trap/trap/core/logger"
	"github.com/raincious/trap/trap/core/types"

	"bytes"
	"fmt"
	"sync"
	"testing"
)

type dummyWriteData struct {
	Length int
	Error  error
}

type dummyReadWriter struct {
	ReadChan  chan []byte
	WriteChan chan dummyWriteData
}

func (c *dummyReadWriter) Write(b []byte) (int, error) {
	writeData := dummyWriteData{
		Length: len(b),
		Error:  nil,
	}

	c.WriteChan <- writeData

	return writeData.Length, writeData.Error
}

func (c *dummyReadWriter) Read(b []byte) (int, error) {
	bCap := cap(b)

	chanData := <-c.ReadChan

	cLen := len(chanData)

	if bCap > cLen {
		bCap = cLen
	}

	for i := 0; i < bCap; i++ {
		b[i] = chanData[i]
	}

	return len(chanData), nil
}

func TestMessagerPack(t *testing.T) {
	testDataEsc := []byte{
		SYNC_CONTROLCHAR_ESCAPE,
		SYNC_CONTROLCHAR_TRANSMITTED,
		SYNC_CONTROLCHAR_SEPARATOR,
		SYNC_CONTROLCHAR_RESERVED,
	}
	testDataSegment := [][]byte{
		[]byte(
			"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890",
		),
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc,
		[]byte(
			"11111",
		),
	}
	expecting := string([]byte{
		0, 0, 0, 0, 0, 0, 2, 49, 50, 51, 52, 53, 54, 55,
		56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56,
		57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48,
		49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49,
		50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50,
		51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51,
		52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52,
		53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53,
		54, 55, 56, 57, 48, 2, 0, 0, 0, 1, 0, 2, 0,
		255, 2, 0, 0, 0, 1, 0, 2, 0, 255, 2, 0, 0,
		0, 1, 0, 2, 0, 255, 2, 0, 0, 0, 1, 0, 2, 0,
		255, 2, 0, 0, 0, 1, 0, 2, 0, 255, 2, 0, 0,
		0, 1, 0, 2, 0, 255, 2, 0, 0, 0, 1, 0, 2, 0,
		255, 2, 0, 0, 0, 1, 0, 2, 0, 255, 2, 0, 0,
		0, 1, 0, 2, 0, 255, 2, 0, 0, 0, 1, 0, 2, 0,
		255, 2, 0, 0, 0, 1, 0, 2, 0, 255, 2, 0, 0,
		0, 1, 0, 2, 0, 255, 2, 0, 0, 0, 1, 0, 2, 0,
		255, 2, 0, 0, 0, 1, 0, 2, 0, 255, 2, 0, 0, 0,
		1, 0, 2, 0, 255, 2, 0, 0, 0, 1, 0, 2, 0, 255,
		2, 49, 49, 49, 49, 49, 1,
	})
	messager := Messager{}

	result := messager.pack(0, 0, 0, testDataSegment)

	if string(result) != string(expecting) {
		t.Error("Messager.pack() failed to pack up data for sending")

		return
	}
}

func TestMessagerFillBytes(t *testing.T) {
	messager := Messager{}
	expecting1 := string([]byte{
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65,
	})
	expecting2 := string([]byte{
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
		65, 65, 65, 65, 65, 65, 65, 65, 65, 65,
	})

	if string(messager.fillBytes(10, byte('A'))) != expecting1 {
		t.Error("Messager.fillBytes() failed to generate correctly amount " +
			"of filling bytes")

		return
	}

	if string(messager.fillBytes(16, byte('A'))) != expecting2 {
		t.Error("Messager.fillBytes() failed to generate correctly amount " +
			"of filling bytes")

		return
	}

	if len(messager.fillBytes(251, byte('A'))) != 5 {
		t.Error("Messager.fillBytes() failed to generate correctly amount " +
			"of filling bytes")

		return
	}

	if len(messager.fillBytes(511, byte('A'))) != 1 {
		t.Error("Messager.fillBytes() failed to generate correctly amount " +
			"of filling bytes")

		return
	}

	if len(messager.fillBytes(4096, byte('A'))) != 0 {
		t.Error("Messager.fillBytes() failed to generate correctly amount " +
			"of filling bytes")

		return
	}
}

func TestMessagerFindBytePositions(t *testing.T) {
	seg := []byte{}
	result := [][]byte{}
	messager := Messager{}
	segHead := 0
	str := []byte{
		'T', 'e', 's',
		't', ' ', 'S',
		'\\', 't',
	}

	ctlCharTable := byteCheckTable{}

	ctlCharTable[byte('t')] = true

	clPos, _, _ := messager.findBytePositions(
		0,
		ctlCharTable,
		str,
		byte('\\'),
	)

	for _, ctrlPos := range clPos {
		seg = append(seg, str[segHead:ctrlPos]...)

		result = append(result, seg)

		seg = nil

		segHead = ctrlPos
	}

	result = append(result, str[segHead:len(str)])

	if len(result) != 3 {
		t.Errorf("Unexpected result length, expecting '3', got '%d'", len(result))

		return
	}

	if string(result[0]) != "Tes" ||
		string(result[1]) != "t S" ||
		string(result[2]) != "\\t" {
		t.Errorf("Messager.findBytePositions() resulting an invalid byte positions")

		return
	}
}

func benchmarkMessagerFindBytePositions(messager *Messager,
	searchFor byteCheckTable, data []byte) int {
	clPos, _, _ := messager.findBytePositions(
		0,
		searchFor,
		data,
		byte('\\'),
	)

	return len(clPos)
}

func BenchmarkMessagerFindBytePositions(b *testing.B) {
	messager := Messager{}
	ctlCharTable := byteCheckTable{}

	testDataEsc := []byte{
		SYNC_CONTROLCHAR_ESCAPE,
		SYNC_CONTROLCHAR_TRANSMITTED,
		SYNC_CONTROLCHAR_SEPARATOR,
		SYNC_CONTROLCHAR_RESERVED,
	}
	testDataSegment := [][]byte{
		[]byte(
			"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890",
		),
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
	}

	data := messager.pack(0, 0, 0, testDataSegment)

	for _, cChar := range testDataEsc {
		ctlCharTable[cChar] = true
	}

	for i := 0; i < b.N; i++ {
		benchmarkMessagerFindBytePositions(&messager, ctlCharTable, data)
	}
}

func TestMessagerBypassBytePosition(t *testing.T) {
	str1 := []byte("........... a \\.test. string")
	str2 := []byte("This is.1.... a \\.test. string")
	messager := Messager{}
	checkTable := byteCheckTable{}
	endIdx := 0
	stopLoop := false

	checkTable[byte('.')] = true
	checkTable[byte(' ')] = true

	clPoses, _, _ := messager.findBytePositions(
		0,
		checkTable,
		str1,
		byte('\\'),
	)

	for clIdx, clPos := range clPoses {
		switch str1[clPos] {
		case byte(' '):

		case byte('.'): // End control char
			endIdx = clIdx
			stopLoop = true
		}

		if stopLoop {
			break
		}
	}

	if string(
		str1[messager.bypassBytePosition(clPoses, endIdx, str1, byte('.'))+1:],
	) != " a \\.test. string" {
		t.Errorf("Messager.bypassBytePosition() failed to correctly " +
			"bypass positions")

		return
	}

	clPoses, _, _ = messager.findBytePositions(
		0,
		checkTable,
		str2,
		byte('\\'),
	)

	stopLoop = false

	for clIdx, clPos := range clPoses {
		switch str2[clPos] {
		case byte(' '):

		case byte('.'): // End control char
			endIdx = clIdx
			stopLoop = true
		}

		if stopLoop {
			break
		}
	}

	if string(
		str2[messager.bypassBytePosition(clPoses, endIdx, str2, byte('.'))+1:],
	) != "1.... a \\.test. string" {
		t.Errorf("Messager.bypassBytePosition() failed to correctly " +
			"bypass positions")

		return
	}
}

func TestMessagerEscapeUnescape(t *testing.T) {
	messager := Messager{}
	data := []byte{
		SYNC_CONTROLCHAR_TRANSMITTED,
		SYNC_CONTROLCHAR_TRANSMITTED,
		'A',
		SYNC_CONTROLCHAR_SEPARATOR,
		'C',
		SYNC_CONTROLCHAR_SEPARATOR,
		SYNC_CONTROLCHAR_SEPARATOR,
	}

	if !bytes.Equal(messager.escape(data), []byte{
		SYNC_CONTROLCHAR_ESCAPE, SYNC_CONTROLCHAR_TRANSMITTED,
		SYNC_CONTROLCHAR_ESCAPE, SYNC_CONTROLCHAR_TRANSMITTED,
		'A',
		SYNC_CONTROLCHAR_ESCAPE, SYNC_CONTROLCHAR_SEPARATOR,
		'C',
		SYNC_CONTROLCHAR_ESCAPE, SYNC_CONTROLCHAR_SEPARATOR,
		SYNC_CONTROLCHAR_ESCAPE, SYNC_CONTROLCHAR_SEPARATOR,
	}) {
		t.Errorf("Messager.escape() failed to convert control charactor " +
			"to expected result")

		return
	}

	if !bytes.Equal(messager.unescape(messager.escape(data)), data) {
		t.Errorf("Messager.unescape() failed convert "+
			"bytes back to it's original form. Expecting: '%d', got '%d'",
			data, messager.unescape(messager.escape(data)))

		return
	}

	if !bytes.Equal(messager.unescape(data), data) {
		t.Errorf("Messager.unescape() safe check failed. "+
			"Expecting: '%d', got '%d'", data, messager.unescape(data))

		return
	}
}

func TestMessagerCombine(t *testing.T) {
	messager := Messager{}
	data := [][]byte{
		[]byte("This is some test data"),
		[]byte{SYNC_CONTROLCHAR_TRANSMITTED, 'A', SYNC_CONTROLCHAR_SEPARATOR},
	}

	combined := messager.combine(data)

	if !bytes.Equal(combined, []byte{84, 104, 105, 115, 32, 105, 115, 32, 115,
		111, 109, 101, 32, 116, 101, 115, 116, 32, 100, 97, 116, 97, 2, 0, 1,
		65, 0, 2}) {
		t.Errorf("messager.combine() failed to product expected result")

		return
	}
}

func TestMessagerParse(t *testing.T) {
	wait := sync.WaitGroup{}
	dummyIO := &dummyReadWriter{
		ReadChan:  make(chan []byte),
		WriteChan: make(chan dummyWriteData),
	}
	responds := Callbacks{}
	dummyResponder := func(req Request) *types.Throw {
		return nil
	}

	responds.Register(SYNC_SIGNAL_UNDEFINED, dummyResponder)
	responds.Register(SYNC_SIGNAL_HEATBEAT, dummyResponder)
	responds.Register(SYNC_SIGNAL_HELLO, dummyResponder)
	responds.Register(SYNC_SIGNAL_HELLO_ACCEPT, dummyResponder)
	responds.Register(SYNC_SIGNAL_HELLO_DENIED, dummyResponder)
	responds.Register(SYNC_SIGNAL_HELLO_CONFLICT, dummyResponder)

	messager := NewMessager(responds, 10240, logger.NewLogger())

	testDataEsc := []byte{
		SYNC_CONTROLCHAR_ESCAPE,
		SYNC_CONTROLCHAR_TRANSMITTED,
		SYNC_CONTROLCHAR_SEPARATOR,
		SYNC_CONTROLCHAR_RESERVED,
	}
	testDataOne := [][]byte{
		[]byte("!SO you know, this is just a simple test!"),
		testDataEsc,
		[]byte("!SO you know, this is just a simple test!"),
	}
	testDataTwo := [][]byte{
		[]byte(
			"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes" +
				"ANd this is another testing data, and this one must longer than 256 bytes",
		),
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
	}
	fragmentized := func(b []byte) [][]byte {
		result := [][]byte{}
		temp := []byte{}
		curTmpLen := 0

		for _, data := range b {
			temp = append(temp, data)
			curTmpLen += 1

			if curTmpLen < SYNC_BUFFER_LENGTH {
				continue
			}

			result = append(result, temp)

			curTmpLen = 0
			temp = nil
		}

		if len(temp) > 0 {
			result = append(result, temp)
		}

		return result
	}
	curTestIndex := 0
	errored := ""

	wait.Add(1)

	go func() {
		defer wait.Done()

		for _, val := range fragmentized(messager.pack(0, 0,
			SYNC_SIGNAL_UNDEFINED, testDataOne)) {
			dummyIO.ReadChan <- val
		}

		for _, val := range fragmentized(messager.pack(0, 0,
			SYNC_SIGNAL_UNDEFINED, testDataTwo)) {
			dummyIO.ReadChan <- val
		}
	}()

	messager.parse(dummyIO, &[]byte{}, func(b []byte) *types.Throw {
		dataIdx := curTestIndex - 1

		if dataIdx >= 0 {
			if !bytes.Equal(b, testDataOne[dataIdx]) {
				errored = fmt.Sprintf("Messager.parse() failed to parse stream data "+
					"correctly. Expecting '%d', got '%d'",
					testDataOne[dataIdx], b)
			}
		}

		curTestIndex++

		return nil
	}, func(int, error) *types.Throw {
		return nil
	})

	if errored != "" {
		t.Error(errored)

		return
	}

	curTestIndex = 0

	messager.parse(dummyIO, &[]byte{}, func(b []byte) *types.Throw {
		dataIdx := curTestIndex - 1

		if dataIdx >= 0 {
			if !bytes.Equal(b, testDataTwo[dataIdx]) {
				errored = fmt.Sprintf("Messager.parse() failed to parse stream data "+
					"correctly. Expecting '%d', got '%d'",
					testDataTwo[dataIdx], b)
			}
		}

		curTestIndex++

		return nil
	}, func(int, error) *types.Throw {
		return nil
	})

	if errored != "" {
		t.Error(errored)
	}

	wait.Wait()
}

func benchmarkMessagerParse(data [][]byte, wait sync.WaitGroup,
	messager *Messager, dummyIO *dummyReadWriter) int {
	length := 0

	wait.Add(1)

	go func() {
		defer wait.Done()

		for _, val := range data {
			dummyIO.ReadChan <- val
		}
	}()

	messager.parse(dummyIO, &[]byte{}, func(b []byte) *types.Throw {
		length += len(b)

		return nil
	}, func(int, error) *types.Throw {
		return nil
	})

	wait.Wait()

	return length
}

func BenchmarkMessagerParse(b *testing.B) {
	wait := sync.WaitGroup{}
	dummyIO := &dummyReadWriter{
		ReadChan:  make(chan []byte),
		WriteChan: make(chan dummyWriteData),
	}
	responds := Callbacks{}
	dummyResponder := func(req Request) *types.Throw {
		return nil
	}

	responds.Register(SYNC_SIGNAL_UNDEFINED, dummyResponder)
	responds.Register(SYNC_SIGNAL_HEATBEAT, dummyResponder)
	responds.Register(SYNC_SIGNAL_HELLO, dummyResponder)
	responds.Register(SYNC_SIGNAL_HELLO_ACCEPT, dummyResponder)
	responds.Register(SYNC_SIGNAL_HELLO_DENIED, dummyResponder)
	responds.Register(SYNC_SIGNAL_HELLO_CONFLICT, dummyResponder)

	messager := NewMessager(responds, 10240, logger.NewLogger())

	testDataEsc := []byte{
		SYNC_CONTROLCHAR_ESCAPE,
		SYNC_CONTROLCHAR_TRANSMITTED,
		SYNC_CONTROLCHAR_SEPARATOR,
		SYNC_CONTROLCHAR_RESERVED,
	}
	testDataSegment := [][]byte{
		[]byte(
			"12345678901234567890123456789012345678901234567890" +
				"12345678901234567890123456789012345678901234567890",
		),
		testDataEsc, testDataEsc, testDataEsc, testDataEsc, testDataEsc,
	}

	testRawData := [][]byte{}

	fragmentized := func(b []byte) [][]byte {
		result := [][]byte{}
		temp := []byte{}
		curTmpLen := 0

		for _, data := range b {
			temp = append(temp, data)
			curTmpLen += 1

			if curTmpLen < SYNC_BUFFER_LENGTH {
				continue
			}

			result = append(result, temp)

			curTmpLen = 0
			temp = nil
		}

		if len(temp) > 0 {
			result = append(result, temp)
		}

		return result
	}

	for i := 0; i < 1000; i++ {
		testRawData = append(testRawData, testDataSegment...)
	}

	testData := fragmentized(messager.pack(0, 0,
		SYNC_SIGNAL_UNDEFINED, testRawData))

	for i := 0; i < b.N; i++ {
		benchmarkMessagerParse(testData, wait, messager, dummyIO)
	}
}

func TestMessagerParsePreservedBuff(t *testing.T) {
	wait := sync.WaitGroup{}
	dummyIO := &dummyReadWriter{
		ReadChan:  make(chan []byte),
		WriteChan: make(chan dummyWriteData),
	}
	responds := Callbacks{}
	dummyResponder := func(req Request) *types.Throw {
		return nil
	}
	pBuffer := []byte{}

	responds.Register(SYNC_SIGNAL_UNDEFINED, dummyResponder)
	responds.Register(SYNC_SIGNAL_HEATBEAT, dummyResponder)
	responds.Register(SYNC_SIGNAL_HELLO, dummyResponder)
	responds.Register(SYNC_SIGNAL_HELLO_ACCEPT, dummyResponder)
	responds.Register(SYNC_SIGNAL_HELLO_DENIED, dummyResponder)
	responds.Register(SYNC_SIGNAL_HELLO_CONFLICT, dummyResponder)

	messager := NewMessager(responds, 10240, logger.NewLogger())

	testDataEsc := []byte{
		'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A',
		SYNC_CONTROLCHAR_TRANSMITTED,
		'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A',
	}

	wait.Add(1)

	go func() {
		defer wait.Done()

		dummyIO.ReadChan <- testDataEsc
	}()

	messager.parse(dummyIO, &pBuffer,
		func(b []byte) *types.Throw { return nil },
		func(int, error) *types.Throw {
			return nil
		})

	wait.Wait()

	if string(pBuffer) != "AAAAAAAAAA" {
		t.Errorf("Messager.parse() failed to output preserved buffer")

		return
	}
}

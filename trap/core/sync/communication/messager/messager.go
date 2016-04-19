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
	"github.com/raincious/trap/trap/core/sync/communication/conn"
	"github.com/raincious/trap/trap/core/types"

	"io"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrMessageExpired *types.Error = types.NewError(
		"`Sync` Message '%d' is expired")

	ErrMessageCleared *types.Error = types.NewError(
		"`Sync` Message '%d' has been cleared")

	ErrMessageEOFReached *types.Error = types.NewError(
		"EOF Reached")

	ErrMessageLengthExceed *types.Error = types.NewError(
		"Recevied '%d' bytes of message, " +
			"but it exceed the max length limit of '%d' bytes")

	ErrMessageUnexpectedResp *types.Error = types.NewError(
		"Unexpected respond '%d' for message '%d'")

	ErrMessageRespondNotReady *types.Error = types.NewError(
		"`Sync` Messager does not ready to " +
			"handle the respond of message '%d'")

	ErrMessageDropped *types.Error = types.NewError(
		"`Sync` Message '%d' has been dropped")

	ErrMessageUnwritable *types.Error = types.NewError(
		"`Sync` Messager can't write any message for now")
)

type Messager struct {
	messages slots

	readerReady messageSignalChan
	writerReady messageSignalChan

	writeable     bool
	writeableLock types.Mutex

	writerChan        chan *message
	exitChan          messageSignalChan
	defaultResponders Callbacks

	transmited int64
	received   int64

	syncCtlCharTable byteCheckTable
}

func NewMessager(defaultResponders Callbacks) *Messager {
	controlCharTable := byteCheckTable{}

	for _, ctrlChar := range []byte{
		SYNC_CONTROLCHAR_TRANSMITTED,
		SYNC_CONTROLCHAR_SEPARATOR,
		SYNC_CONTROLCHAR_RESERVED,
	} {
		controlCharTable[ctrlChar] = true
	}

	return &Messager{
		messages: slots{
			initLock: &types.Mutex{},
		},

		readerReady: make(messageSignalChan),
		writerReady: make(messageSignalChan),

		writeable:     false,
		writeableLock: types.Mutex{},

		writerChan:        make(chan *message, MAX_MESSAGES_HOLDING_SIZE),
		exitChan:          make(messageSignalChan),
		defaultResponders: defaultResponders,

		transmited: 0,
		received:   0,

		syncCtlCharTable: controlCharTable,
	}
}

func (m *Messager) bypassBytePosition(positions []int, lastPosIndex int,
	buf []byte, char byte) int {
	posLen := len(positions)
	lastPos := positions[lastPosIndex]
	tmpLastPos := int(0)

	for i := lastPosIndex; i < posLen; i++ {
		if buf[positions[i]] != char {
			break
		}

		lastPos = positions[i]
	}

	if lastPos == positions[lastPosIndex] {
		return lastPos
	}

	tmpLastPos = lastPos

	for i := positions[lastPosIndex]; i <= tmpLastPos; i++ {
		if buf[i] != char {
			break
		}

		lastPos = i
	}

	return lastPos
}

func (m *Messager) findBytePositions(start int, finds byteCheckTable, in []byte,
	escaper byte) ([]int, int, bool) {
	byteLen := len(in)
	found := false
	indexes := make([]int, byteLen-start)
	founds := int(0)
	lastPos := int(-1)

	for i := start; i < byteLen; i++ {
		if in[i] == escaper {
			indexes[founds] = i
			founds++

			found = true
			lastPos = i

			i++

			continue
		}

		if !finds[in[i]] {
			continue
		}

		indexes[founds] = i
		founds++

		found = true
		lastPos = i
	}

	if lastPos == -1 {
		return []int{}, -1, false
	}

	return indexes[:founds], lastPos, found
}

func (m *Messager) escape(b []byte) []byte {
	result := []byte{}
	bLen := len(b)

	for i := 0; i < bLen; i++ {
		switch b[i] {
		case SYNC_CONTROLCHAR_ESCAPE:
			fallthrough
		case SYNC_CONTROLCHAR_TRANSMITTED:
			fallthrough
		case SYNC_CONTROLCHAR_SEPARATOR:
			fallthrough
		case SYNC_CONTROLCHAR_RESERVED:
			result = append(result,
				[]byte{SYNC_CONTROLCHAR_ESCAPE,
					b[i]}...)
		default:
			result = append(result, b[i])
		}
	}

	return result
}

func (m *Messager) unescape(b []byte) []byte {
	bytesLen := len(b)
	nextIdx := int(0)
	rst := make([]byte, bytesLen)
	rstIdx := int(0)

	for i := 0; i < bytesLen; i++ {
		if b[i] != SYNC_CONTROLCHAR_ESCAPE {
			rst[rstIdx] = b[i]
			rstIdx++

			continue
		}

		nextIdx = i + 1

		switch b[nextIdx] {
		case SYNC_CONTROLCHAR_ESCAPE:
			fallthrough
		case SYNC_CONTROLCHAR_TRANSMITTED:
			fallthrough
		case SYNC_CONTROLCHAR_SEPARATOR:
			fallthrough
		case SYNC_CONTROLCHAR_RESERVED:
			rst[rstIdx] = b[nextIdx]

			i++ // Skip next charactor

		default:
			rst[rstIdx] = b[i]
		}

		rstIdx++
	}

	return rst[:rstIdx]
}

func (m *Messager) combine(b [][]byte) []byte {
	result := []byte{}
	resultLen := 0

	for _, bArr := range b {
		result = append(result, m.escape(bArr)...)
		result = append(result, SYNC_CONTROLCHAR_SEPARATOR)
	}

	resultLen = len(result)

	if resultLen < 1 {
		return []byte{}
	}

	return result[:len(result)-1]
}

func (m *Messager) fillBytes(sentLen int, withWith byte) []byte {
	remain := SYNC_BUFFER_LENGTH -
		(sentLen % SYNC_BUFFER_LENGTH)

	fillBytes := remain % SYNC_BUFFER_LENGTH

	if fillBytes == 0 {
		return []byte{}
	}

	filling := make([]byte, fillBytes)

	for idx, _ := range filling {
		filling[idx] = withWith
	}

	return filling
}

func (m *Messager) pack(id, code byte, data [][]byte) []byte {
	packed := []byte{}

	packed = append(packed, m.escape([]byte{id})...)
	packed = append(packed, SYNC_CONTROLCHAR_SEPARATOR)
	packed = append(packed, m.escape([]byte{code})...)
	packed = append(packed, SYNC_CONTROLCHAR_SEPARATOR)
	packed = append(packed, m.combine(data)...)
	packed = append(packed, SYNC_CONTROLCHAR_TRANSMITTED)

	return packed
}

func (m *Messager) parse(reader io.Reader, perseveredBuf *[]byte,
	result func(b []byte),
	afterRead func(int, error) *types.Throw) *types.Throw {
	var err error = nil
	var throw *types.Throw = nil

	var rLen int = 0
	var bufReserved []byte = []byte{}
	var bufResCur int = 0
	var bufResHed int = 0
	var psvBuf []byte = *perseveredBuf
	var psvBufLen int = 0
	var psvBufCutLen int = 0
	var buffer []byte = make([]byte, SYNC_BUFFER_LENGTH)

	var clPos []int = []int{}
	var cllPos int = 0
	var cHad bool = false

	for {
		bufResHed = 0
		buffer = make([]byte, SYNC_BUFFER_LENGTH)

		psvBufLen = len(psvBuf)

		if psvBufLen > 0 {
			if psvBufLen > SYNC_BUFFER_LENGTH {
				psvBufCutLen = SYNC_BUFFER_LENGTH
			} else {
				psvBufCutLen = psvBufLen
			}

			buffer = psvBuf[:psvBufCutLen]
			rLen = psvBufCutLen

			*perseveredBuf = psvBuf[psvBufCutLen:]
			psvBuf = *perseveredBuf

			err = nil
		} else {
			rLen, err = reader.Read(buffer)

			if err != nil {
				return types.ConvertError(err)
			}

			throw = afterRead(rLen, err)

			if throw != nil {
				return types.ConvertError(throw)
			}
		}

		switch err {
		case nil:
			if rLen < 1 {
				continue
			}

			bufReserved = append(bufReserved, buffer[:rLen]...)

			clPos, cllPos, cHad = m.findBytePositions(
				bufResCur,
				m.syncCtlCharTable,
				bufReserved,
				SYNC_CONTROLCHAR_ESCAPE,
			)

			if !cHad {
				bufResCur += rLen

				continue
			}

			if bufReserved[cllPos] == SYNC_CONTROLCHAR_ESCAPE {
				bufResCur = cllPos
				clPos = clPos[:len(clPos)-1]
			}

			for ctrlIdx, ctrlPos := range clPos {
				switch bufReserved[ctrlPos] {
				case SYNC_CONTROLCHAR_ESCAPE:
					// Do nothing

				case SYNC_CONTROLCHAR_SEPARATOR:
					result(m.unescape(bufReserved[bufResHed:ctrlPos]))

					bufResHed = ctrlPos + 1

				case SYNC_CONTROLCHAR_TRANSMITTED:
					result(m.unescape(bufReserved[bufResHed:ctrlPos]))

					ctrlPos = m.bypassBytePosition(clPos,
						ctrlIdx, bufReserved,
						SYNC_CONTROLCHAR_TRANSMITTED)

					bufResHed = ctrlPos + 1

					*perseveredBuf = bufReserved[bufResHed:]

					return nil

				case SYNC_CONTROLCHAR_RESERVED:
					// Yeah, do nothing again
				}
			}

			bufReserved = bufReserved[bufResHed:]
			bufResCur = 0

		case io.EOF:
			return ErrMessageEOFReached.Throw()

		default:
			return types.ConvertError(err)
		}
	}
}

func (m *Messager) reader(rConn *conn.Conn) *types.Throw {
	var message *message = nil
	var readIndex uint = 0
	var totalReadLength uint = 0
	var data [][]byte = [][]byte{}
	var respCode byte = SYNC_SIGNAL_UNDEFINED
	var responable bool = true
	var prsvBuf []byte = []byte{}
	var replyID byte = 0
	var parseErr *types.Throw = nil

	defer func() {
		m.exitChan <- true
	}()

	m.readerReady <- true

	for {
		message = nil
		readIndex = 0
		totalReadLength = 0
		data = [][]byte{}
		respCode = SYNC_SIGNAL_UNDEFINED
		responable = true
		replyID = 0

		parseErr = m.parse(rConn, &prsvBuf, func(b []byte) {
			readIndex++

			switch readIndex {
			case 0:
				// Drop it

			case 1:
				if len(b) > 0 {
					replyID = b[0]

					msg, msgErr := m.messages.Take(replyID)

					if msgErr != nil {
						return
					}

					msg.StatusLock.Exec(func() {
						if msg.Ready {
							message = msg

							return
						}

						msg.ResultChan <- ErrMessageRespondNotReady.Throw(
							msg.ID)

						responable = false
					})
				}

			case 2:
				if len(b) > 0 {
					respCode = b[0]
				}

			default:
				data = append(data, b)
			}
		}, func(readLen int, readErr error) *types.Throw {
			atomic.AddInt64(&m.received, int64(readLen))

			if readErr != nil {
				return nil
			}

			totalReadLength += uint(readLen)

			if message == nil {
				return nil
			}

			if message.MaxRespondLen == 0 {
				return nil
			}

			if totalReadLength < message.MaxRespondLen {
				return nil
			}

			return ErrMessageLengthExceed.Throw(totalReadLength,
				message.MaxRespondLen)
		})

		if parseErr != nil {
			if parseErr.Is(ErrMessageLengthExceed) {
				continue
			}

			return parseErr
		}

		if !responable {
			continue
		}

		if message != nil {
			if !message.Responds.Has(respCode) {
				message.ResultChan <- ErrMessageUnexpectedResp.Throw(
					respCode, message.ID)

				continue
			}

			message.ResultChan <- message.Responds.Call(
				respCode,
				Request{
					conn:     rConn,
					messager: m,

					data:    data,
					dataLen: totalReadLength,
					code:    respCode,
					id:      0,

					isReplyable: false, // Not request, thus unreplyable
				},
			)
		} else {
			if !m.defaultResponders.Has(respCode) {
				continue
			}

			m.defaultResponders.Call(respCode, Request{
				conn:     rConn,
				messager: m,

				data:    data,
				dataLen: totalReadLength,
				code:    respCode,
				id:      replyID,

				isReplyable: true, // This is a request, so replyable
			})
		}
	}
}

func (m *Messager) writer(wConn *conn.Conn) {
	var wLen int = 0
	var wErr error = nil
	var segWrtLen int = 0

	down := false
	preDown := false
	ticker := time.Tick(1 * time.Second)

	defer func() {
		if preDown {
			return
		}

		m.messages.Deinit()
	}()

	// Init the message container
	m.messages.Init()

	m.writerReady <- true

	m.writeableLock.Exec(func() {
		m.writeable = true
	})

	for {
		select {
		case message := <-m.writerChan:
			if preDown {
				// Take the saved message out the container
				if !message.Held {
					message.ResultChan <- ErrMessageDropped.Throw(
						message.ID)
				} else {
					m.messages.Drop(message.ID, ErrMessageDropped.Throw(
						message.ID))
				}

				continue
			}

			data := m.pack(message.ID, message.Code,
				message.Message)

			if len(m.writerChan) <= 0 {
				data = append(data, m.fillBytes(
					len(data)+segWrtLen,
					SYNC_CONTROLCHAR_TRANSMITTED)...)

				segWrtLen = 0
			}

			message.StatusLock.Exec(func() {
				wLen, wErr = wConn.Write(data)

				if wErr != nil {
					return
				}

				message.Ready = true
			})

			if wErr != nil {
				if !message.Held {
					message.ResultChan <- types.ConvertError(wErr)
				} else {
					m.messages.Drop(message.ID,
						types.ConvertError(wErr))
				}

				continue
			}

			segWrtLen += wLen

			atomic.AddInt64(&m.transmited, int64(wLen))

			if !message.Held {
				message.ResultChan <- nil

				continue
			}

		case <-m.exitChan:
			if preDown {
				continue
			}

			preDown = true

			m.writeableLock.Exec(func() {
				m.writeable = false
			})

			go func() {
				m.messages.Deinit()

				down = true
			}()

		case <-ticker:
			if !down {
				continue
			}

			return
		}
	}
}

func (m *Messager) Listen(lConn *conn.Conn, ready chan<- bool) *types.Throw {
	var err *types.Throw = nil

	listenWait := sync.WaitGroup{}

	listenWait.Add(2)

	go func() {
		defer listenWait.Done()

		m.writer(lConn)
	}()

	go func() {
		defer listenWait.Done()

		err = m.reader(lConn)
	}()

	<-m.writerReady
	<-m.readerReady

	ready <- true

	listenWait.Wait()

	return err
}

func (m *Messager) Stats() Stats {
	return Stats{
		TX: atomic.LoadInt64(&m.transmited),
		RX: atomic.LoadInt64(&m.received),
	}
}

func (m *Messager) Query(code byte, data Data, responds Callbacks,
	maxRespondLen uint, waitTime time.Duration) *types.Throw {
	var err *types.Throw = nil
	var msgByte [][]byte = [][]byte{}

	msgByte, err = data.Build()

	if err != nil {
		return err
	}

	msg := &message{
		ID:            MESSAGES_RESEVERED_ID,
		Code:          code,
		Held:          false,
		Message:       msgByte,
		Responds:      responds,
		ResultChan:    make(chan *types.Throw),
		MaxRespondLen: maxRespondLen,
		Ready:         false,
		StatusLock:    types.Mutex{},
	}

	holdErr := m.messages.Hold(
		msg,
		waitTime,
		func(reason MessageDeleteReason, error *types.Throw) {
			var outputRst bool = false

			msg.StatusLock.Exec(func() {
				if error != nil {
					msg.ResultChan <- error

					msg.Ready = false

					return
				}

				switch reason {
				case MESSAGE_DELETE_REASON_CLEAR:
					err = ErrMessageCleared.Throw(msg.ID)
					outputRst = true

				case MESSAGE_DELETE_REASON_EXPIRE:
					err = ErrMessageExpired.Throw(msg.ID)
					outputRst = true
				}

				if !outputRst {
					return
				}

				msg.ResultChan <- err

				msg.Ready = false
			})
		})

	if holdErr != nil {
		return holdErr
	}

	m.writeableLock.Exec(func() {
		if m.writeable {
			return
		}

		err = ErrMessageUnwritable.Throw()
	})

	if err != nil {
		m.messages.Drop(msg.ID, nil)

		return err
	}

	m.writerChan <- msg

	return <-msg.ResultChan
}

func (m *Messager) Reply(msgID byte, code byte, data Data) *types.Throw {
	var err *types.Throw = nil
	var msgByte [][]byte = [][]byte{}

	msgByte, err = data.Build()

	if err != nil {
		return err
	}

	msg := &message{
		ID:            msgID,
		Code:          code,
		Held:          false,
		Message:       msgByte,
		Responds:      m.defaultResponders,
		ResultChan:    make(chan *types.Throw),
		MaxRespondLen: 0,
		Ready:         false,
		StatusLock:    types.Mutex{},
	}

	m.writeableLock.Exec(func() {
		if m.writeable {
			return
		}

		err = ErrMessageUnwritable.Throw()
	})

	if err != nil {
		return err
	}

	m.writerChan <- msg

	return <-msg.ResultChan
}

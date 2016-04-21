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
	"github.com/raincious/trap/trap/core/types"

	"sync"
	"time"
)

const (
	MAX_MESSAGES_HOLDING_SIZE = uint16(^byte(0)) + 1
	MESSAGES_RESEVERED_ID     = ^byte(0)
)

const (
	MESSAGE_DELETE_REASON_CLEAR MessageDeleteReason = iota
	MESSAGE_DELETE_REASON_EXPIRE
	MESSAGE_DELETE_REASON_TAKEN
	MESSAGE_DELETE_REASON_DROP
)

var (
	ErrSlotsChannelFailed *types.Error = types.NewError(
		"Can't get channel data from slot %d")

	ErrSlotsDisabled *types.Error = types.NewError(
		"`Slots` has been disabled")

	ErrSlotsNotFound *types.Error = types.NewError(
		"Can't found message %d")
)

type slots struct {
	inited             bool
	initLock           *types.Mutex
	initWait           sync.WaitGroup
	messages           [MAX_MESSAGES_HOLDING_SIZE]*messageSlot
	messageEnabled     bool
	messageEnabledLock *types.Mutex
	messageIndex       byte
	messageIndexLock   *types.Mutex

	monitoringExit messageSignalChan
}

func (m *slots) Init() *types.Throw {
	var err *types.Throw = nil

	m.initLock.Exec(func() {
		err = m.init()
	})

	return err
}

func (m *slots) Deinit() *types.Throw {
	var err *types.Throw = nil

	m.initLock.Exec(func() {
		err = m.deinit()
	})

	return err
}

func (m *slots) init() *types.Throw {
	if m.inited {
		return nil
	}

	m.inited = true

	m.messages = [MAX_MESSAGES_HOLDING_SIZE]*messageSlot{}
	m.messageIndex = byte(0)
	m.messageIndexLock = &types.Mutex{}
	m.messageEnabledLock = &types.Mutex{}

	m.initWait = sync.WaitGroup{}
	m.monitoringExit = make(messageSignalChan)

	for id, _ := range m.messages {
		m.messages[id] = &messageSlot{
			msg: nil,

			enabled:    false,
			busyChan:   make(messageSignalChan),
			expire:     time.Time{},
			lock:       types.Mutex{},
			insertLock: types.Mutex{},
			deleter:    nil,
		}
	}

	m.enable()

	m.initWait.Add(1)

	go m.monitoring()

	m.messageEnabledLock.Exec(func() {
		m.messageEnabled = true
	})

	return nil
}

func (m *slots) deinit() *types.Throw {
	if !m.inited {
		return nil
	}

	m.messageEnabledLock.Exec(func() {
		m.messageEnabled = false
	})

	m.disable()

	m.monitoringExit <- true

	m.initWait.Wait()

	for {
		if m.remains() <= 0 {
			break
		}

		m.clear()
	}

	m.inited = false

	return nil
}

func (m *slots) monitoring() {
	defer m.initWait.Done()

	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-tick:
			m.expire()

		case <-m.monitoringExit:
			return
		}
	}
}

func (m *slots) nextIndex() byte {
	var result byte = 0

	m.messageIndexLock.Exec(func() {
		newIndex := m.messageIndex + 1

		result = byte(uint16(newIndex) %
			MAX_MESSAGES_HOLDING_SIZE)

		m.messageIndex = result
	})

	return result
}

func (m *slots) has(id byte) bool {
	if m.messages[id].msg == nil {
		return false
	}

	return true
}

func (m *slots) delete(id byte, rs MessageDeleteReason, err *types.Throw) {
	if !m.has(id) {
		return
	}

	deleter := m.messages[id].deleter

	m.messages[id].msg = nil
	m.messages[id].expire = time.Time{}
	m.messages[id].deleter = nil

	if deleter != nil {
		deleter(rs, err)
	}

	select {
	case m.messages[id].busyChan <- false:
	default:
	}
}

func (m *slots) scan(callback func(id byte, msg *messageSlot)) {
	for i := int(MAX_MESSAGES_HOLDING_SIZE) - 1; i >= 0; i-- {
		callback(byte(i), m.messages[byte(i)])
	}
}

func (m *slots) lockedScan(callback func(id byte, msg *messageSlot)) {
	m.scan(func(id byte, msg *messageSlot) {
		msg.lock.Exec(func() {
			callback(id, msg)
		})
	})
}

func (m *slots) insertLockedScan(callback func(id byte, msg *messageSlot)) {
	m.scan(func(id byte, msg *messageSlot) {
		msg.insertLock.Exec(func() {
			callback(id, msg)
		})
	})
}

func (m *slots) clear() {
	m.lockedScan(func(id byte, msg *messageSlot) {
		m.delete(id, MESSAGE_DELETE_REASON_CLEAR, nil)
	})
}

func (m *slots) remains() uint {
	var totalEnabled uint = 0

	m.insertLockedScan(func(id byte, msg *messageSlot) {
		if msg.msg == nil {
			return
		}

		totalEnabled++
	})

	return totalEnabled
}

func (m *slots) enable() {
	m.insertLockedScan(func(id byte, msg *messageSlot) {
		msg.enabled = true
	})
}

func (m *slots) disable() {
	m.insertLockedScan(func(id byte, msg *messageSlot) {
		msg.enabled = false
	})
}

func (m *slots) expire() {
	now := time.Now()

	m.lockedScan(func(id byte, msg *messageSlot) {
		if msg.msg == nil {
			return
		}

		if !now.After(msg.expire) {
			return
		}

		m.delete(id, MESSAGE_DELETE_REASON_EXPIRE, nil)
	})
}

func (m *slots) Drop(id byte, error *types.Throw) *types.Throw {
	var err *types.Throw = nil

	m.messages[id].lock.Exec(func() {
		if !m.has(id) {
			err = ErrSlotsNotFound.Throw(id)

			return
		}

		m.delete(id, MESSAGE_DELETE_REASON_DROP, error)
	})

	return err
}

func (m *slots) Take(id byte) (*message, *types.Throw) {
	var msg *message = nil
	var err *types.Throw = nil

	m.messages[id].lock.Exec(func() {
		if !m.has(id) {
			err = ErrSlotsNotFound.Throw(id)

			return
		}

		msg = m.messages[id].msg

		m.delete(id, MESSAGE_DELETE_REASON_TAKEN, nil)
	})

	return msg, err
}

func (m *slots) Hold(msg *message, expire time.Duration,
	deleter func(MessageDeleteReason, *types.Throw)) *types.Throw {
	var waitChan messageSignalChan
	var err *types.Throw = nil
	var waitUp bool = false
	var waitTime time.Duration = time.Duration(0)
	var idx byte = m.nextIndex()

	m.initLock.Exec(func() {
		m.messageEnabledLock.Exec(func() {
			if m.messageEnabled {
				return
			}

			err = ErrSlotsDisabled.Throw()
		})
	})

	if err != nil {
		return err
	}

	m.messages[idx].insertLock.Exec(func() {
		// Check again
		m.messageEnabledLock.Exec(func() {
			if m.messageEnabled {
				return
			}

			err = ErrSlotsDisabled.Throw()
		})

		if err != nil {
			return
		}

		if !m.messages[idx].enabled {
			err = ErrSlotsDisabled.Throw()

			return
		}

		m.messages[idx].lock.Exec(func() {
			if m.messages[idx].msg == nil {
				return
			}

			waitUp = true
			waitChan = m.messages[idx].busyChan
			waitTime = m.messages[idx].expire.Sub(time.Now())
		})

		if waitUp {
			select {
			case _, ok := <-waitChan:
				if !ok {
					err = ErrSlotsChannelFailed.Throw(idx)

					return
				}

			case <-time.After(waitTime):
				m.messages[idx].lock.Exec(func() {
					m.delete(idx, MESSAGE_DELETE_REASON_EXPIRE, nil)
				})
			}
		}

		m.messages[idx].lock.Exec(func() {
			msg.ID = idx
			msg.Held = true

			m.messages[idx].msg = msg
			m.messages[idx].expire = time.Now().Add(expire)
			m.messages[idx].deleter = deleter
		})
	})

	return err
}

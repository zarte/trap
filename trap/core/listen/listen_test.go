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

package listen

import (
	"github.com/raincious/trap/trap/core/logger"
	"github.com/raincious/trap/trap/core/types"

	"net"
	"testing"
	"time"
)

var (
	ErrProtocolFakeErr  *types.Error = types.NewError("Protocol fake error")
	ErrProtocolFakeErr2 *types.Error = types.NewError("Protocol fake error 2")

	ErrListenerFakeErr  *types.Error = types.NewError("Listener fake error")
	ErrListenerFakeErr2 *types.Error = types.NewError("Listener fake error 2")
)

// Fake protocol
type fakeProtocol struct {
	returnError bool

	initCalled         bool
	initCalledNSucceed bool
	spawnCalled        bool

	config *ProtocolConfig

	settings types.String
}

func (f *fakeProtocol) Init(cfg *ProtocolConfig) *types.Throw {
	f.initCalled = true

	if f.returnError {
		return ErrProtocolFakeErr.Throw()
	}

	f.config = cfg

	f.initCalledNSucceed = true

	return nil
}

func (f *fakeProtocol) Spawn(setting types.String) (Listener, *types.Throw) {
	f.spawnCalled = true

	if f.returnError {
		return nil, ErrProtocolFakeErr2.Throw()
	}

	ipStr, portStr := setting.SpiltWith("@")

	ip := net.ParseIP(ipStr.String())

	if ip == nil {
		ip = net.ParseIP("0.0.0.0")
	}

	return &fakeListener{
		ip:   ip,
		port: portStr.UInt16(),
	}, nil
}

var (
	fakeListenerReturnError bool = false
)

// Fake Listener
type fakeListener struct {
	ip   net.IP
	port types.UInt16

	upCalled   bool
	downCalled bool
}

func (f *fakeListener) Up() (*ListeningInfo, *types.Throw) {
	f.upCalled = true

	if fakeListenerReturnError {
		return nil, ErrListenerFakeErr.Throw()
	}

	return &ListeningInfo{
		Port:     8080,
		IP:       net.IP{},
		Protocol: "udp",
	}, nil
}

func (f *fakeListener) Down() (*ListeningInfo, *types.Throw) {
	f.downCalled = true

	if fakeListenerReturnError {
		return nil, ErrListenerFakeErr2.Throw()
	}

	return &ListeningInfo{
		Port:     8080,
		IP:       net.IP{},
		Protocol: "udp",
	}, nil
}

func TestListenRegister(t *testing.T) {
	listen := Listen{}
	logger := logger.NewLogger()
	fakePl1 := &fakeProtocol{}
	fakePl2 := &fakeProtocol{}
	onErrCalled := false
	onPikCalled := false

	listen.Init(&Config{
		OnError: func(cInfo ConnectionInfo, err *types.Throw) {
			onErrCalled = true
		},
		OnPick: func(cInfo ConnectionInfo, rInfo RespondedResult) {
			onPikCalled = true
		},
		OnListened: func(lInfo *ListeningInfo) {

		},
		OnUnListened: func(lInfo *ListeningInfo) {

		},
		MaxBytes:   512,
		Logger:     logger,
		Concurrent: 100,
		Timeout:    1 * time.Second,
	})

	// Succeed register
	regErr := listen.Register("test", fakePl1)

	if regErr != nil {
		t.Errorf("Operation failed due to error: %s", regErr)

		return
	}

	if len(listen.protocols) != 1 {
		t.Error("Listen.Register() seems didn't add the new protocol into `protocols` map")

		return
	}

	// Listen.Register() should call Protocol.Init() for succeed registation
	if !fakePl1.initCalled {
		t.Error("Listen.Register() didn't call Protocol.Init()")

		return
	}

	if !fakePl1.initCalledNSucceed {
		t.Error("Listen.Register() didn't finish Protocol.Init()")

		return
	}

	if fakePl1.config.OnError == nil || onErrCalled != false ||
		fakePl1.config.OnPick == nil || onPikCalled != false ||
		fakePl1.config.MaxBytes != listen.maxBytes ||
		fakePl1.config.ReadTimeout != listen.timeout ||
		fakePl1.config.WriteTimeout != listen.timeout ||
		fakePl1.config.TotalTimeout != listen.timeout ||
		fakePl1.config.Logger != listen.logger ||
		fakePl1.config.Concurrent != listen.concurrent {
		t.Error("Listen.Register() failed to pass the setting to protocol")

		return
	}

	fakePl1.config.OnError(ConnectionInfo{}, nil)
	fakePl1.config.OnPick(ConnectionInfo{}, RespondedResult{})

	if onErrCalled != true || onPikCalled != true {
		t.Error("Listen.Register() failed to pass the setting to protocol")

		return
	}

	onErrCalled = false
	onPikCalled = false // reset

	// Register should be failed due to duplicate
	regErr = listen.Register("test", fakePl2)

	if regErr == nil || !regErr.Is(ErrProtocolAlreadyRegistered) {
		t.Error("Expected error doesn't happen")

		return
	}

	if len(listen.protocols) != 1 {
		t.Error("Listen.Register() seems added a wrong protocol into the map")

		return
	}

	if fakePl2.initCalled {
		t.Error("Listen.Register() called Protocol.Init() on a invalid protocol")

		return
	}
}

func TestListenRegisterErrored(t *testing.T) {
	listen := Listen{}
	logger := logger.NewLogger()
	fakePl1 := &fakeProtocol{
		returnError: true,
	}
	onErrCalled := false
	onPikCalled := false

	listen.Init(&Config{
		OnError: func(cInfo ConnectionInfo, err *types.Throw) {
			onErrCalled = true
		},
		OnPick: func(cInfo ConnectionInfo, rInfo RespondedResult) {
			onPikCalled = true
		},
		OnListened: func(lInfo *ListeningInfo) {

		},
		OnUnListened: func(lInfo *ListeningInfo) {

		},
		MaxBytes:   512,
		Logger:     logger,
		Concurrent: 100,
		Timeout:    1 * time.Second,
	})

	regErr := listen.Register("test", fakePl1)

	if regErr == nil || !regErr.Is(ErrProtocolFakeErr) {
		t.Error("Listen.Register() didn't return expected error")

		return
	}

	if !fakePl1.initCalled {
		t.Error("Listen.Register() didn't call the Protocol.Init()")

		return
	}

	if fakePl1.initCalledNSucceed {
		t.Error("Listen.Register() unexpectly finished Protocol.Init()")

		return
	}

	if len(listen.protocols) != 0 {
		t.Error("Listen.Register() seems added a invalid protocol into the map")

		return
	}
}

func TestListenAdd(t *testing.T) {
	listen := Listen{}
	logger := logger.NewLogger()
	fakePl := &fakeProtocol{}

	listen.Init(&Config{
		OnError: func(cInfo ConnectionInfo, err *types.Throw) {

		},
		OnPick: func(cInfo ConnectionInfo, rInfo RespondedResult) {

		},
		OnListened: func(lInfo *ListeningInfo) {

		},
		OnUnListened: func(lInfo *ListeningInfo) {

		},
		MaxBytes:   512,
		Logger:     logger,
		Concurrent: 100,
		Timeout:    1 * time.Second,
	})

	regErr := listen.Register("test", fakePl)

	if regErr != nil {
		t.Error("Unexpected error when trying to register fake protocol")

		return
	}

	ptcErr := listen.Add("test", "8080@0.0.0.0")

	if ptcErr != nil {
		t.Error("Unexpected error when trying to register fake listener")

		return
	}

	// Try add another one with non-registered protocol type,
	// should return an error for me
	ptcErr = listen.Add("test2", "8080@0.0.0.0")

	if ptcErr == nil || !ptcErr.Is(ErrProtocolNotSupported) {
		t.Error("Expected error didn't happen")

		return
	}

	fakePl.returnError = true

	ptcErr = listen.Add("test", "8080@0.0.0.0")

	if ptcErr == nil || !ptcErr.Is(ErrProtocolFakeErr2) {
		t.Error("Protocol.Spawn() didn't return the expected error")

		return
	}
}

func TestListenServ(t *testing.T) {
	listen := Listen{}
	logger := logger.NewLogger()
	fakePl := &fakeProtocol{}
	onErrCalled := false
	onPikCalled := false
	onLisCalled := false
	onUnlCalled := false

	listen.Init(&Config{
		OnError: func(cInfo ConnectionInfo, err *types.Throw) {
			onErrCalled = true
		},
		OnPick: func(cInfo ConnectionInfo, rInfo RespondedResult) {
			onPikCalled = true
		},
		OnListened: func(lInfo *ListeningInfo) {
			onLisCalled = true
		},
		OnUnListened: func(lInfo *ListeningInfo) {
			onUnlCalled = true
		},
		MaxBytes:   512,
		Logger:     logger,
		Concurrent: 100,
		Timeout:    1 * time.Second,
	})

	regErr := listen.Register("test", fakePl)

	if regErr != nil {
		t.Errorf("listen.Register() failed due to error: %s", regErr)

		return
	}

	ptcErr := listen.Add("test", "8080@0.0.0.0")

	if ptcErr != nil {
		t.Errorf("listen.Add() failed due to error: %s", regErr)

		return
	}

	// Let the fake listener don't make any mistake
	fakeListenerReturnError = false

	lisErr := listen.Serv()

	if lisErr != nil {
		t.Errorf("listen.Serv() failed to being up Listener due to error: %s", lisErr)

		return
	}

	if !onLisCalled {
		t.Error("listen.Serv() didn't call `OnListened` callback")

		return
	}

	// Try again, this time with an error Listener
	fakeListenerReturnError = true

	lisErr = listen.Serv()

	if lisErr == nil || !lisErr.Is(ErrListenerFakeErr) {
		t.Error("listen.Serv() failed pickup expected error")

		return
	}
}

func TestListenDown(t *testing.T) {
	listen := Listen{}
	logger := logger.NewLogger()
	fakePl := &fakeProtocol{}
	onErrCalled := false
	onPikCalled := false
	onLisCalled := false
	onUnlCalled := false

	listen.Init(&Config{
		OnError: func(cInfo ConnectionInfo, err *types.Throw) {
			onErrCalled = true
		},
		OnPick: func(cInfo ConnectionInfo, rInfo RespondedResult) {
			onPikCalled = true
		},
		OnListened: func(lInfo *ListeningInfo) {
			onLisCalled = true
		},
		OnUnListened: func(lInfo *ListeningInfo) {
			onUnlCalled = true
		},
		MaxBytes:   512,
		Logger:     logger,
		Concurrent: 100,
		Timeout:    1 * time.Second,
	})

	regErr := listen.Register("test", fakePl)

	if regErr != nil {
		t.Errorf("listen.Register() failed due to error: %s", regErr)

		return
	}

	ptcErr := listen.Add("test", "8080@0.0.0.0")

	if ptcErr != nil {
		t.Errorf("listen.Add() failed due to error: %s", regErr)

		return
	}

	// Let the fake listener don't make any mistake
	fakeListenerReturnError = false

	lisErr := listen.Down()

	if lisErr != nil {
		t.Errorf("listen.Down() failed to being up Listener due to error: %s",
			lisErr)

		return
	}

	if !onUnlCalled {
		t.Error("listen.Down() didn't call `OnUnListened` callback")

		return
	}

	// Try again, this time with an error Listener
	fakeListenerReturnError = true

	lisErr = listen.Down()

	if lisErr == nil || !lisErr.Is(ErrListenerFakeErr2) {
		t.Error("listen.Down() failed pickup expected error")

		return
	}
}

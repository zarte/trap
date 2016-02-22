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
    "net"
    "errors"
)

var (
    ErrHostInvalidIP *Error =
        NewError("'%s' is not an IP address")


    emptyIP                 =   IP{}
)

type IP [16]byte

func (ip IP) IP() (net.IP) {
    if ip.IsEmpty() {
        return net.IP{}
    }

    return net.IP(ip[:])
}

func (ip *IP) String() (string) {
    return ip.IP().String()
}

func (ip *IP) IsEqual(anotherIP *IP) (bool) {
    return *ip == *anotherIP
}

func (ip *IP) IsEmpty() (bool) {
    if !ip.IsEqual(&emptyIP) {
        return false
    }

    return true
}

func (ip *IP) MarshalText() ([]byte, error) {
    return []byte(ip.String()), nil
}

func (ip *IP) UnmarshalText(text []byte) (error) {
    newIP               :=  IP{}

    var err *Throw      =   nil

    newIP, err          =   ConvertIPFromString(String(text[:]))

    if err != nil {
        return errors.New(err.Error())
    }

    *ip                 =   newIP

    return nil
}

type HostAddress struct {
    Host        String
    Port        UInt16
}

func (a *HostAddress) IsEmpty() (bool) {
    if a.Host == "" {
        return true
    }

    if a.Port == 0 {
        return true
    }

    return false
}

type IPAddress struct {
    IP          IP
    Port        UInt16
}

func (a *IPAddress) IsEmpty() (bool) {
    if a.IP.IsEmpty() {
        return true
    }

    if a.Port == 0 {
        return true
    }

    return false
}

func ConvertAddress(addr net.Addr) (HostAddress, *Throw) {
    aHost, aPort, aErr  :=  net.SplitHostPort(addr.String())

    if aErr != nil {
        return HostAddress{ Host: "", Port: 0 }, ConvertError(aErr)
    }

    return HostAddress{
        Host:           String(aHost),
        Port:           String(aPort).UInt16(),
    }, nil
}

func ConvertIP(ip net.IP) (IP) {
    newIP   :=  IP{}

    bip     :=  ip.To16()

    copy(newIP[:], bip[:])

    return newIP
}

func ConvertIPFromString(addr String) (IP, *Throw) {
    ipAddr              :=  net.ParseIP(addr.String())

    if ipAddr == nil {
        return IP{}, ErrHostInvalidIP.Throw(addr)
    }

    return ConvertIP(ipAddr), nil
}

func ConvertIPAddress(addr net.Addr) (IPAddress, *Throw) {
    aHost, aPort, aErr  :=  net.SplitHostPort(addr.String())

    if aErr != nil {
        return IPAddress{ IP: IP{}, Port: 0 }, ConvertError(aErr)
    }

    ip, ipErr           :=  ConvertIPFromString(String(aHost))

    if ipErr != nil {
        return IPAddress{ IP: IP{}, Port: 0 }, ipErr
    }

    return IPAddress{
        IP:             ip,
        Port:           String(aPort).UInt16(),
    }, nil
}
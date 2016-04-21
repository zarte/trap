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
	"encoding/binary"
	"errors"
	"math"
	"net"
)

const (
	IP_ADDR_SLICE_LEN = 16
	IP_PORT_LEN       = 2
)

var (
	ErrHostInvalidIP *Error = NewError("'%s' is not an IP address")

	emptyIP = IP{}
)

type IP [IP_ADDR_SLICE_LEN]byte

func (ip IP) IP() net.IP {
	if ip.IsEmpty() {
		return net.IP{}
	}

	return net.IP(ip[:])
}

func (ip *IP) String() string {
	return ip.IP().String()
}

func (ip *IP) IsEqual(anotherIP *IP) bool {
	return *ip == *anotherIP
}

func (ip *IP) IsEmpty() bool {
	if !ip.IsEqual(&emptyIP) {
		return false
	}

	return true
}

func (ip *IP) MarshalText() ([]byte, error) {
	return []byte(ip.String()), nil
}

func (ip *IP) UnmarshalText(text []byte) error {
	newIP := IP{}

	var err *Throw = nil

	newIP, err = ConvertIPFromString(String(text[:]))

	if err != nil {
		return errors.New(err.Error())
	}

	*ip = newIP

	return nil
}

type HostAddress struct {
	Host String
	Port UInt16
}

func (a *HostAddress) String() string {
	return net.JoinHostPort(
		a.Host.String(), a.Port.String().String())
}

func (a *HostAddress) IsEmpty() bool {
	if a.Host == "" {
		return true
	}

	if a.Port == 0 {
		return true
	}

	return false
}

type IPAddress struct {
	IP   IP
	Port UInt16
}

func (a *IPAddress) Host() *HostAddress {
	return &HostAddress{
		Host: String(a.IP.String()),
		Port: a.Port,
	}
}

func (a *IPAddress) IsEqual(addr *IPAddress) bool {
	if !a.IP.IsEqual(&addr.IP) {
		return false
	}

	if a.Port != addr.Port {
		return false
	}

	return true
}

func (a *IPAddress) IsEmpty() bool {
	if a.IP.IsEmpty() {
		return true
	}

	if a.Port == 0 {
		return true
	}

	return false
}

func (a IPAddress) String() String {
	return String(a.IP.String()).Join(":", a.Port.String())
}

func (a *IPAddress) Serialize() ([]byte, *Throw) {
	result := [IP_ADDR_SLICE_LEN + IP_PORT_LEN]byte{}
	portByte := make([]byte, IP_PORT_LEN)

	for ipIdx, _ := range a.IP {
		result[ipIdx] = a.IP[ipIdx]
	}

	binary.LittleEndian.PutUint16(portByte, uint16(a.Port))

	for portIdx, _ := range portByte {
		result[portIdx+IP_ADDR_SLICE_LEN] = portByte[portIdx]
	}

	return result[:], nil
}

func (a *IPAddress) Unserialize(text []byte) *Throw {
	newIP := IP{}
	tmpPt := [IP_PORT_LEN]byte{}

	if len(text) != IP_ADDR_SLICE_LEN+IP_PORT_LEN {
		return ErrTypesUnserializeInvalidDataLength.Throw(
			IP_ADDR_SLICE_LEN + IP_PORT_LEN)
	}

	for ipIdx := 0; ipIdx < IP_ADDR_SLICE_LEN; ipIdx++ {
		newIP[ipIdx] = text[ipIdx]
	}

	if newIP.IsEmpty() {
		return ErrHostInvalidIP.Throw(newIP.String())
	}

	for prtIdx, _ := range tmpPt {
		tmpPt[prtIdx] = text[prtIdx+IP_ADDR_SLICE_LEN]
	}

	a.IP = newIP
	a.Port = UInt16(binary.LittleEndian.Uint16(tmpPt[:]))

	return nil
}

type IPAddresses []IPAddress

func (ipAddrs *IPAddresses) Serialize() ([]byte, *Throw) {
	var err *Throw = nil
	const IP_FULL_LEN = IP_ADDR_SLICE_LEN + IP_PORT_LEN

	segment := []byte{}
	result := make([]byte, IP_FULL_LEN*len(*ipAddrs))
	curSegStart := 0

	for _, ipAddr := range *ipAddrs {
		segment, err = ipAddr.Serialize()

		if err != nil {
			return nil, err
		}

		for byteIdx, sgBytes := range segment {
			result[byteIdx+curSegStart] = sgBytes
		}

		curSegStart += IP_FULL_LEN
	}

	return result[:], nil
}

func (ipAddrs *IPAddresses) Unserialize(data []byte) *Throw {
	totalSegLen := int64(IP_ADDR_SLICE_LEN + IP_PORT_LEN)
	totalTxtLen := int64(len(data))
	segStart := 0
	segEnds := IP_ADDR_SLICE_LEN + IP_PORT_LEN

	if totalTxtLen%totalSegLen != 0 {
		return ErrTypesUnserializeInvalidDataLength.Throw(totalSegLen)
	}

	nIPAddrs := make(IPAddresses,
		int64(
			math.Ceil(
				float64(totalTxtLen/totalSegLen))))

	for ipIdx, _ := range nIPAddrs {
		segEnds = segStart + IP_ADDR_SLICE_LEN + IP_PORT_LEN

		err := nIPAddrs[ipIdx].Unserialize(data[segStart:segEnds])

		if err != nil {
			return err
		}

		segStart = segEnds
	}

	*ipAddrs = nIPAddrs

	return nil
}

func (ipAddrs *IPAddresses) Contains(companions *IPAddresses) Int64 {
	result := Int64(0)

	for _, ipAddr := range *ipAddrs {
		for _, companion := range *companions {
			if !companion.IsEqual(&ipAddr) {
				continue
			}

			result += 1
		}
	}

	return result
}

func (ipAddrs *IPAddresses) Intersection(companions *IPAddresses) IPAddresses {
	intersection := IPAddresses{}

	for _, ipAddr := range *ipAddrs {
		for _, companion := range *companions {
			if !companion.IsEqual(&ipAddr) {
				continue
			}

			intersection = append(intersection, companion)
		}
	}

	return intersection
}

func ConvertAddress(addr net.Addr) (HostAddress, *Throw) {
	aHost, aPort, aErr := net.SplitHostPort(addr.String())

	if aErr != nil {
		return HostAddress{Host: "", Port: 0}, ConvertError(aErr)
	}

	return HostAddress{
		Host: String(aHost),
		Port: String(aPort).UInt16(),
	}, nil
}

func ConvertIP(ip net.IP) IP {
	newIP := IP{}

	bip := ip.To16()

	copy(newIP[:], bip[:])

	return newIP
}

func ConvertIPFromString(addr String) (IP, *Throw) {
	ipAddr := net.ParseIP(addr.String())

	if ipAddr == nil {
		return IP{}, ErrHostInvalidIP.Throw(addr)
	}

	return ConvertIP(ipAddr), nil
}

func ConvertIPAddress(addr net.Addr) (IPAddress, *Throw) {
	aHost, aPort, aErr := net.SplitHostPort(addr.String())

	if aErr != nil {
		return IPAddress{IP: IP{}, Port: 0}, ConvertError(aErr)
	}

	ip, ipErr := ConvertIPFromString(String(aHost))

	if ipErr != nil {
		return IPAddress{IP: IP{}, Port: 0}, ipErr
	}

	return IPAddress{
		IP:   ip,
		Port: String(aPort).UInt16(),
	}, nil
}

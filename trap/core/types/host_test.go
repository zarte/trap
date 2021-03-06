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
	"testing"
)

type fakeNetAddress struct {
	IPPort string
	Net    string
}

func (f fakeNetAddress) String() string {
	return f.IPPort
}

func (f fakeNetAddress) Network() string {
	return f.Net
}

// Test IP struct
func TestIPStruct(t *testing.T) {
	ip := IP{}
	ip2 := IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ip3 := IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 127, 0, 0, 1}
	ip4 := IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}

	if ip.String() != "<nil>" {
		t.Errorf("Failed asserting string format of an empty IP is '<nil>'")

		return
	}

	if ip.IP().String() != ip.String() {
		t.Errorf("When convert the IP to net.IP, the data is broken")

		return
	}

	if !ip.IP().Equal(ip2.IP()) { // Empty IP == Another type of empty IP
		t.Errorf("When convert the IP to net.IP, the data is broken")

		return
	}

	if !ip.IsEmpty() {
		t.Errorf("Failed asserting an empty IP is empty")

		return
	}

	if !ip.IsEqual(&emptyIP) {
		t.Errorf("Failed asserting an empty IP is equal to another empty IP")

		return
	}

	marshaledIP, marshalErr := ip.MarshalText()

	if marshalErr != nil {
		t.Errorf("Error happened when trying to marshal the IP: %s", marshalErr)

		return
	}

	if String(marshaledIP) != "<nil>" {
		t.Errorf("Marshaled empty IP is not excepted '<nil>', got '%s'", marshaledIP)

		return
	}

	unmarshalErr := ip.UnmarshalText([]byte("0.0.0.0"))

	if unmarshalErr != nil {
		t.Errorf("Can't unmarshal text due to error: %s", unmarshalErr)

		return
	}

	if ip.String() != "0.0.0.0" {
		t.Errorf("Failed asserting the IP '%s' is '0.0.0.0'", ip.String())

		return
	}

	if ip.IP().String() != ip.String() {
		t.Errorf("When convert the IP to net.IP, the data is broken")

		return
	}

	if ip.IP().Equal(ip2.IP()) {
		// the `ip` is not empty now so it can't be equal with `ip2`
		t.Errorf("When convert the IP to net.IP, the data is broken")

		return
	}

	if ip.IsEmpty() {
		t.Errorf("Failed asserting an filled IP is filled")

		return
	}

	if !ip.IsEqual(&ip) {
		t.Errorf("Failed asserting the IP is equal to itself")

		return
	}

	if ip.IsEqual(&ip2) {
		t.Errorf("Failed asserting the IP is equals to another IP which contains same value")

		return
	}

	if ip.IsEqual(&ip3) {
		t.Errorf("Failed asserting the IP is not equals to another IP which contains different value")

		return
	}

	if !ip2.IsZero() {
		t.Errorf("Failed asserting the IP '%s' is zero", ip2.String())

		return
	}

	if ip3.String() != "127.0.0.1" {
		t.Errorf("Failed asserting the IP is excepted '127.0.0.1', got '%s'", ip3.String())

		return
	}

	if ip4.String() != "2001:db8:85a3::8a2e:370:7334" {
		t.Errorf("Failed asserting the IP is excepted '2001:db8:85a3::8a2e:370:7334', got '%s'", ip4.String())

		return
	}

	if ip3.IsEqual(&ip4) {
		t.Errorf("Failed asserting two different IP is not equal")

		return
	}

	if ip4.IsEqual(&ip3) {
		t.Errorf("Failed asserting two different IP is not equal")

		return
	}

	if ip4.IsEmpty() {
		t.Errorf("Failed asserting a filled IPv6 address is not empty")

		return
	}
}

// Test HostAddress struct
func TestHostAddressStruct(t *testing.T) {
	emptyAddr1 := HostAddress{
		Host: "",
		Port: 0,
	}

	emptyAddr2 := HostAddress{
		Host: "localhost",
		Port: 0,
	}

	emptyAddr3 := HostAddress{
		Host: "",
		Port: 123,
	}

	if !emptyAddr1.IsEmpty() || !emptyAddr2.IsEmpty() || !emptyAddr3.IsEmpty() {
		t.Errorf("Failed asserting an empty Host Address is not empty")

		return
	}

	filledAddr := HostAddress{
		Host: "localhost",
		Port: 123,
	}

	if filledAddr.IsEmpty() {
		t.Errorf("Failed asserting an filled Host Address is empty")

		return
	}

	if filledAddr.String() != "localhost:123" {
		t.Errorf("Failed to convert HostAddress to string")

		return
	}

	if filledAddr.Host != "localhost" || filledAddr.Port != 123 {
		t.Errorf("Struct data somehow mutated unexpectly")

		return
	}
}

// Test IPAddress struct
func TestIPAddressStruct(t *testing.T) {
	testIP := IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}

	emptyAddr1 := IPAddress{
		IP:   IP{},
		Port: 0,
	}

	emptyAddr2 := IPAddress{
		IP:   IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52},
		Port: 0,
	}

	emptyAddr3 := IPAddress{
		IP:   IP{},
		Port: 123,
	}

	if !emptyAddr1.IsEmpty() || !emptyAddr2.IsEmpty() || !emptyAddr3.IsEmpty() {
		t.Errorf("Failed asserting an empty IP Address is not empty")

		return
	}

	filledAddr := IPAddress{
		IP:   testIP,
		Port: 123,
	}

	if filledAddr.Host().String() != "[2001:db8:85a3::8a2e:370:7334]:123" {
		t.Errorf("Failed to convert HostAddress to string")

		return
	}

	if filledAddr.IsEmpty() {
		t.Errorf("Failed asserting an filled IP Address is empty")

		return
	}

	if !filledAddr.IP.IsEqual(&testIP) || filledAddr.Port != 123 {
		t.Errorf("Struct data somehow mutated unexpectly")

		return
	}
}

// Test ConvertIPAddress function
func testConvertIPAddress(t *testing.T, ip string, expectedIP string, port UInt16) {
	testIPAddress, testIPAddrErr := ConvertIPAddress(fakeNetAddress{
		IPPort: ip + ":" + port.String().String(),
	})

	if testIPAddrErr != nil {
		t.Errorf("Unexpected error returned when trying to convert "+
			"IP Address: %s", testIPAddrErr)

		return
	}

	if testIPAddress.IP.String() != expectedIP {
		t.Errorf("ConvertIPAddress failed to convert IP address. "+
			"Expecting address to be '%s', got '%s'",
			expectedIP, testIPAddress.IP.String())

		return
	}

	if testIPAddress.Port != port {
		t.Errorf("ConvertIPAddress failed to convert IP address. "+
			"Expecting port to be '%d', got '%d'", port, testIPAddress.Port)

		return
	}

	if testIPAddress.String() !=
		String(testIPAddress.IP.String()).Join(":", port.String()) {
		t.Errorf("ConvertIPAddress failed to convert IP address. "+
			"Expecting String() to be '%s', got '%s'",
			ip+":"+port.String().String(), testIPAddress.String())

		return
	}
}

func TestConvertIPAddressV4(t *testing.T) {
	testConvertIPAddress(t, "127.0.0.1", "127.0.0.1", 0)
}

func TestConvertIPAddressV6(t *testing.T) {
	testConvertIPAddress(t, "[fe80::222:68ff:fea8:56bd]",
		"fe80::222:68ff:fea8:56bd", 0)
}

func TestConvertIPAddressInvalid(t *testing.T) {
	_, testIPAddrErr := ConvertIPAddress(fakeNetAddress{
		IPPort: "this.is.not.an.IP:0",
	})

	if testIPAddrErr == nil {
		t.Errorf("Expecting some error, but nothing happened")

		return
	}

	if !testIPAddrErr.Is(ErrHostInvalidIP) {
		t.Errorf("Expecting error '%s', got: %s",
			ErrHostInvalidIP.Throw("[An IP Here]"), testIPAddrErr)

		return
	}
}

// Test ConvertIPFromString function
func testConvertIPFromString(t *testing.T, ipStr String) {
	ip, ipErr := ConvertIPFromString(ipStr)
	ip2, _ := ConvertIPFromString(ipStr)

	if ipErr != nil {
		t.Errorf("Unexpected error: %s", ipErr)

		return
	}

	if ip.String() != ipStr.String() {
		t.Errorf("Expecting '%s', got: %s",
			ipStr, ip.String())

		return
	}

	if !ip.IsEqual(&ip2) {
		t.Errorf("Failed asserting '%s' is equals to '%s'",
			ip.String(), ip2.String())

		return
	}

	if ip.IsEmpty() {
		t.Errorf("Failed asserting '%s' is not empty", ip.String())

		return
	}
}

func TestConvertIPFromStringV4(t *testing.T) {
	testConvertIPFromString(t, "127.0.0.1")
	testConvertIPFromString(t, "2001:db8:85a3::8a2e:370:7334")
}

func TestConvertIPFromStringInvalid(t *testing.T) {
	_, testIPAddrErr := ConvertIPAddress(fakeNetAddress{
		IPPort: "this.is.not.an.IP:0",
	})

	if testIPAddrErr == nil {
		t.Errorf("Expecting some error, but nothing happened")

		return
	}

	if !testIPAddrErr.Is(ErrHostInvalidIP) {
		t.Errorf("Expecting error '%s', got: %s",
			ErrHostInvalidIP.Throw("[An IP Here]"), testIPAddrErr)

		return
	}
}

// Test ConvertIP function
func testConvertIP(t *testing.T, ipStr string) {
	rawIP := net.ParseIP(ipStr)

	ip := ConvertIP(rawIP)

	if ip.String() != ipStr {
		t.Errorf("Failed to covert IP, excepting '%s', got '%s'",
			ipStr, ip.String())
	}
}

func TestConvertIP(t *testing.T) {
	testConvertIP(t, "127.0.0.1")
	testConvertIP(t, "2001:db8:85a3::8a2e:370:7334")
}

func TestConvertIPInvalid(t *testing.T) {
	rawIP := net.ParseIP("Not an IP")

	ip := ConvertIP(rawIP)

	if !ip.IsEmpty() {
		t.Errorf("Excepting the IP address to be empty, but it's not")

		return
	}

	if ip.String() != "<nil>" {
		t.Errorf("Excepting the IP address to be '<nil>', got '%s'",
			ip.String())

		return
	}
}

// Test ConvertAddress function
func testConvertAddress(t *testing.T, host String, exceptingHost String, port UInt16) {
	addr, err := ConvertAddress(fakeNetAddress{
		IPPort: String(host + ":" + port.String()).String(),
		Net:    "",
	})

	if err != nil {
		t.Errorf("Unexpected error: %s", err)

		return
	}

	if addr.Host != exceptingHost {
		t.Errorf("Failed to convert Host Address. Excepting host to be '%s', got '%s'",
			exceptingHost, addr.Host.String())

		return
	}

	if addr.Port != port {
		t.Errorf("Failed to convert Host Address. Excepting port to be '%d', got '%d'",
			port, addr.Port)

		return
	}
}

func TestConvertAddress(t *testing.T) {
	testConvertAddress(t, "www.google.com", "www.google.com", 443)
	testConvertAddress(t, "hostname", "hostname", 443)
	testConvertAddress(t, "127.0.0.1", "127.0.0.1", 443)
	testConvertAddress(t, "[2001:db8:85a3::8a2e:370:7334]", "2001:db8:85a3::8a2e:370:7334", 443)
	testConvertAddress(t, "[::]", "::", 443)
}

func TestIPAddressInvalid(t *testing.T) {
	addr, err := ConvertAddress(fakeNetAddress{
		IPPort: "::0",
		Net:    "",
	})

	if err == nil {
		t.Errorf("Excepting error, but got nothing")
	}

	if !addr.IsEmpty() {
		t.Errorf("Excepting the address is empty, got something else")
	}
}

func TestIPAddressesSerializeUnarshalText(t *testing.T) {
	ips := IPAddresses{}
	ips2 := IPAddresses{}

	testIPAddress, _ := ConvertIPAddress(fakeNetAddress{
		IPPort: "127.0.0.1:8080",
	})

	ips = append(ips, testIPAddress)
	ips = append(ips, testIPAddress)

	out, _ := ips.Serialize()

	ips2.Unserialize(out)

	if !ips2[0].IsEqual(&ips[0]) || !ips2[1].IsEqual(&ips[1]) {
		t.Error("IPAddress.Serialize() or IPAddress.Unserialize() " +
			"is failed")

		return
	}

	if ips2[0].IP.String() != "127.0.0.1" {
		t.Error("IPAddress.Serialize() or IPAddress.Unserialize() " +
			"is failed")

		return
	}

	if ips2[1].Port != 8080 {
		t.Error("IPAddress.Serialize() or IPAddress.Unserialize() " +
			"is failed")

		return
	}
}

func TestIPAddressesContains(t *testing.T) {
	ips := IPAddresses{}
	ips2 := IPAddresses{}

	testIPAddress, _ := ConvertIPAddress(fakeNetAddress{
		IPPort: "127.0.0.1:8080",
	})

	testIPAddress2, _ := ConvertIPAddress(fakeNetAddress{
		IPPort: "127.0.0.2:8080",
	})

	ips = append(ips, testIPAddress2)
	ips2 = append(ips2, testIPAddress)

	if ips.Contains(&ips2) {
		t.Errorf("IPAddresses.Contains() failed to count the right " +
			"amount of companion objects")

		return
	}

	ips = append(ips, testIPAddress)
	ips = append(ips, testIPAddress)

	if !ips.Contains(&ips2) {
		t.Errorf("IPAddresses.Contains() failed to count the right " +
			"amount of companion objects")

		return
	}
}

func TestIPAddressesIntersection(t *testing.T) {
	ips := IPAddresses{}
	ips2 := IPAddresses{}

	testIPAddress, _ := ConvertIPAddress(fakeNetAddress{
		IPPort: "127.0.0.1:8080",
	})

	testIPAddress2, _ := ConvertIPAddress(fakeNetAddress{
		IPPort: "127.0.0.2:8080",
	})

	ips = append(ips, testIPAddress)
	ips2 = append(ips2, testIPAddress2)

	if len(ips.Intersection(&ips2)) != 0 {
		t.Errorf("IPAddresses.Intersection() failed to pickup the right " +
			"amount of companion objects")

		return
	}

	ips2 = append(ips2, testIPAddress)

	intersection := ips.Intersection(&ips2)

	if len(intersection) != 1 {
		t.Errorf("IPAddresses.Intersection() failed to pickup the right " +
			"amount of companion objects")

		return
	}

	if !intersection[0].IsEqual(&testIPAddress) {
		t.Errorf("IPAddresses.Intersection() failed to pickup the right result")

		return
	}
}

func TestIPAddressWrapped(t *testing.T) {
	testIPAddress, _ := ConvertIPAddress(fakeNetAddress{
		IPPort: "127.0.0.1:8080",
	})

	wrapped := testIPAddress.Wrapped()

	wrappedIPAddr := wrapped.IPAddress()

	if !wrappedIPAddr.IsEqual(&testIPAddress) {
		t.Errorf("IPAddressWrapped.Import() failed import IPAddress")

		return
	}

	if wrapped.String() != IPAddressString(testIPAddress.String()) {
		t.Errorf("IPAddressWrapped.Import() failed import IPAddress")

		return
	}
}

func TestSearchableIPAddresses(t *testing.T) {
	testIPAddress, _ := ConvertIPAddress(fakeNetAddress{
		IPPort: "127.0.0.1:8080",
	})

	testIPAddress2, _ := ConvertIPAddress(fakeNetAddress{
		IPPort: "127.0.0.1:8081",
	})

	testIPAddress3, _ := ConvertIPAddress(fakeNetAddress{
		IPPort: "127.0.0.2:8081",
	})

	testIPAddress4, _ := ConvertIPAddress(fakeNetAddress{
		IPPort: "127.0.0.3:8081",
	})

	wrapped := testIPAddress.Wrapped()
	wrapped2 := testIPAddress4.Wrapped()

	ips := IPAddresses{testIPAddress, testIPAddress2, testIPAddress3}
	ips2 := IPAddresses{testIPAddress, testIPAddress2, testIPAddress4}

	searchable := ips.Searchable()
	searchable2 := ips2.Searchable()

	if !searchable.Has(&wrapped) {
		t.Error("SearchableIPAddresses.Import() failed import IPAddress")

		return
	}

	if searchable.Has(&wrapped2) {
		t.Error("SearchableIPAddresses.Has() failed to return a false on fail")

		return
	}

	inter := searchable.Intersection(&searchable2)

	if inter.Len() != 2 ||
		len(inter.Export()) != 2 ||
		!inter.Contains(&searchable) {
		t.Error("SearchableIPAddresses.Intersection() failed to return " +
			"a right amount of result")

		return
	}

	errorString := ""

	inter.Through(func(key IPAddressString, val IPAddressWrapped) *Throw {
		switch key {
		case "127.0.0.1:8080":
			if String(val.String()) != testIPAddress.String() {
				errorString = "SearchableIPAddresses.Intersection() failed " +
					"to return the right result"

				return nil
			}

		case "127.0.0.1:8081":
			if String(val.String()) != testIPAddress2.String() {
				errorString = "SearchableIPAddresses.Intersection() failed " +
					"to return the right result"

				return nil
			}

		default:
			errorString = "SearchableIPAddresses.Intersection() returning an " +
				"unexpected result"

			return nil
		}

		inter.Delete(&val)

		return nil
	})

	if errorString != "" {
		t.Error(errorString)

		return
	}

	if inter.Len() != 0 {
		t.Error("SearchableIPAddresses.Delete() failed to delete " +
			"items")

		return
	}

	if inter.Contains(&searchable) {
		t.Error("SearchableIPAddresses.Contains() failed to return false when" +
			" it's not contains specified item")

		return
	}
}

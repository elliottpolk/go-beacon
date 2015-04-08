// Copyright 2015 elliott@tkwcafe.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package beacon

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/paypal/gatt"
)

func NewServer(name string, manufacturer []byte, uuid string, major, minor int) (*gatt.Server, error) {
	id, err := parseUUID(uuid)
	if err != nil {
		return nil, err
	}

	//  if name is left blank, go ahead and default to iBeacon
	if name == "" {
		name = "iBeacon"
	}

	//  if a manufacturer isn't provided default to the Bluetooth SIG Apple
	//  identifier, flagging it as an iBeacon. For a list of company identifiers:
	//  https://www.bluetooth.org/en-us/specification/assigned-numbers/company-identifiers
	if manufacturer == nil || len(manufacturer) < 2 {
		manufacturer = []byte{0x4C, 0x00}
	}

	//  payload format is based on the info from:
	//  http://stackoverflow.com/questions/18906988/what-is-the-ibeacon-bluetooth-profile
	payload := []byte{}
	payload = append(payload, 0x02) // Number of bytes that follow in first advertising structure
	payload = append(payload, 0x01) // Number of flags
	payload = append(payload, 0x1A) // Flag -> 0x1A = 000011010
	payload = append(payload, 0x1A) // Number of bytes that follow in second advertising structure
	payload = append(payload, 0xFF) // Manufacturer specific data advertising type
	payload = append(payload, manufacturer...)
	payload = append(payload, []byte{0x02, 0x15}...) // iBeacon identifier
	payload = append(payload, id...)
	payload = append(payload, []byte{uint8(major >> 8), uint8(major & 0xff)}...)
	payload = append(payload, []byte{uint8(minor >> 8), uint8(minor & 0xff)}...)

	//  TODO: this will vary based on the beacon hardware
	//  Note, the current value (0xC5) is pulled via the ibeacon profile link
	//  above and will NOT be accurate in most cases. PLEASE CHANGE
	payload = append(payload, 0xC5)

	return &gatt.Server{Name: name, AdvertisingPacket: payload}, nil
}

//  modified via github.com/paypal/gatt/uuid.go
//
//  The original (gatt.ParseUUID) returns a UUID item which "hides" the []byte
//  keeping us from building the advertising payload required for beacon-like
//  functionality
func parseUUID(s string) ([]byte, error) {
	s = strings.Replace(s, "-", "", -1)
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(b) != 16 {
		return nil, errors.New("UUID length must be 16 bytes")
	}

	return b, nil
}

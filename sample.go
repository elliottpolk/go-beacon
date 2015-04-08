// Copyright 2015 elliott@tkwcafe.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"log"

	"beacon"

	"github.com/paypal/gatt"
)

func main() {
	srv, err := beacon.NewServer("MyBeacon", nil, "D9B9EC1F-3925-43D0-80A9-1E39D4CEA95C", 1, 2)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(srv.AdvertiseAndServe())
}

// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package speedtest

import (
	"encoding/json"
	"errors"
	"net"
	"time"
)

// RunClient dials the given address and starts a speedtest.
// It returns any errors that come up in the tests.
// If there are no errors in the test, it returns a slice of results.
func RunClient(direction Direction, duration time.Duration, host string) ([]Result, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	conf := config{TestDuration: duration, Version: version, Direction: direction}

	defer conn.Close()
	encoder := json.NewEncoder(conn)

	if err = encoder.Encode(conf); err != nil {
		return nil, err
	}

	var response configResponse
	decoder := json.NewDecoder(conn)
	if err = decoder.Decode(&response); err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return doTest(conn, conf)
}

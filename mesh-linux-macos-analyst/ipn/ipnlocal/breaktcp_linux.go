// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package ipnlocal

import (
	"log"

	"golang.org/x/sys/unix"
)

func init() {
	breakTCPConns = breakTCPConnsLinux
}

func breakTCPConnsLinux() error {
	var matched int
	for fd := 0; fd < 1000; fd++ {
		_, err := unix.GetsockoptTCPInfo(fd, unix.IPPROTO_TCP, unix.TCP_INFO)
		if err == nil {
			matched++
			err = unix.Close(fd)
			log.Printf("debug: closed TCP fd %v: %v", fd, err)
		}
	}
	if matched == 0 {
		log.Printf("debug: no TCP connections found")
	}
	return nil
}

// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build darwin && !ios

package magicsock

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func getDontFragOpt(network string) int {
	if network == "udp4" {
		return unix.IP_DONTFRAG
	}
	return unix.IPV6_DONTFRAG
}

func (c *Conn) setDontFragment(network string, enable bool) error {
	optArg := 1
	if enable == false {
		optArg = 0
	}
	var err error
	rcErr := c.connControl(network, func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), getIPProto(network), getDontFragOpt(network), optArg)
	})

	if rcErr != nil {
		return rcErr
	}
	return err
}

func (c *Conn) getDontFragment(network string) (bool, error) {
	var v int
	var err error
	rcErr := c.connControl(network, func(fd uintptr) {
		v, err = syscall.GetsockoptInt(int(fd), getIPProto(network), getDontFragOpt(network))
	})

	if rcErr != nil {
		return false, rcErr
	}
	if v == 1 {
		return true, err
	}
	return false, err
}

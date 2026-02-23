// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build linux && !android

package magicsock

import (
	"syscall"
)

func getDontFragOpt(network string) int {
	if network == "udp4" {
		return syscall.IP_MTU_DISCOVER
	}
	return syscall.IPV6_MTU_DISCOVER
}

func (c *Conn) setDontFragment(network string, enable bool) error {
	optArg := syscall.IP_PMTUDISC_DO
	if enable == false {
		optArg = syscall.IP_PMTUDISC_DONT
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
	if v == syscall.IP_PMTUDISC_DO {
		return true, err
	}
	return false, err
}

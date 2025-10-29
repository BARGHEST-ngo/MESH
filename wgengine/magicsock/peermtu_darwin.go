// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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

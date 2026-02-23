// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build cgo && linux

// Package linuxfwtest contains tests for the linuxfw package. Go does not
// support cgo in tests, and we don't want the main package to have a cgo
// dependency, so we put all the tests here and call them from the main package
// in tests intead.
package linuxfwtest

import (
	"testing"
	"unsafe"
)

/*
#include <sys/socket.h>  // socket()
*/
import "C"

type SizeInfo struct {
	SizeofSocklen uintptr
}

func TestSizes(t *testing.T, si *SizeInfo) {
	want := unsafe.Sizeof(C.socklen_t(0))
	if want != si.SizeofSocklen {
		t.Errorf("sockLen has wrong size; want=%d got=%d", want, si.SizeofSocklen)
	}
}

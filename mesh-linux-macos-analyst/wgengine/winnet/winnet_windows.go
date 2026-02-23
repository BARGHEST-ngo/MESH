// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package winnet

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

func (v *INetworkConnection) GetAdapterId() (string, error) {
	buf := ole.GUID{}
	hr, _, _ := syscall.Syscall(
		v.VTable().GetAdapterId,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&buf)),
		0)
	if hr != 0 {
		return "", fmt.Errorf("GetAdapterId failed: %08x", hr)
	}
	return buf.String(), nil
}

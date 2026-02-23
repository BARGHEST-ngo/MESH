// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package pidowner

import (
	"fmt"
	"os"
	"strings"

	"tailscale.com/util/lineiter"
)

func ownerOfPID(pid int) (userID string, err error) {
	file := fmt.Sprintf("/proc/%d/status", pid)
	for lr := range lineiter.File(file) {
		line, err := lr.Value()
		if err != nil {
			if os.IsNotExist(err) {
				return "", ErrProcessNotFound
			}
			return "", err
		}
		if len(line) < 4 || string(line[:4]) != "Uid:" {
			continue
		}
		f := strings.Fields(string(line))
		if len(f) >= 2 {
			userID = f[1] // real userid
		}
	}
	if userID == "" {
		return "", fmt.Errorf("missing Uid line in %s", file)
	}
	return userID, nil
}

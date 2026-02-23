// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.
//go:build gofuzz

package disco

func Fuzz(data []byte) int {
	m, _ := Parse(data)

	newBytes := m.AppendMarshal(data)
	parsedMarshall, _ := Parse(newBytes)

	if m != parsedMarshall {
		panic("Parsing error")
	}
	return 1
}

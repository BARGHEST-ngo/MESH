// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.
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

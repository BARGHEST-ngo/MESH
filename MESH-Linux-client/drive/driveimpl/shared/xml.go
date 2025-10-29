// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package shared

import (
	"bytes"
	"encoding/xml"
)

// EscapeForXML escapes the given string for use in XML text.
func EscapeForXML(s string) string {
	result := bytes.NewBuffer(nil)
	xml.Escape(result, []byte(s))
	return result.String()
}

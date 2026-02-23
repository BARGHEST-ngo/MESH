// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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

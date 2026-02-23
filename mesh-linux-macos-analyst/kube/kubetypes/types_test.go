// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package kubetypes

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalAPIServerProxyMode(t *testing.T) {
	tests := []struct {
		data     string
		expected APIServerProxyMode
	}{
		{data: `{"mode":"auth"}`, expected: APIServerProxyModeAuth},
		{data: `{"mode":"noauth"}`, expected: APIServerProxyModeNoAuth},
		{data: `{"mode":""}`, expected: ""},
		{data: `{"mode":"Auth"}`, expected: ""},
		{data: `{"mode":"unknown"}`, expected: ""},
	}

	for _, tc := range tests {
		var s struct {
			Mode *APIServerProxyMode `json:",omitempty"`
		}
		err := json.Unmarshal([]byte(tc.data), &s)
		if tc.expected == "" {
			if err == nil {
				t.Errorf("expected error for %q, got none", tc.data)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error for %q: %v", tc.data, err)
			continue
		}
		if *s.Mode != tc.expected {
			t.Errorf("for %q expected %q, got %q", tc.data, tc.expected, *s.Mode)
		}
	}
}

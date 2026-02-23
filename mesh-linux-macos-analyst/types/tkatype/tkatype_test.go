// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package tkatype

import (
	"encoding/json"
	"testing"

	"golang.org/x/crypto/blake2s"
)

func TestSigHashSize(t *testing.T) {
	var sigHash AUMSigHash
	if len(sigHash) != blake2s.Size {
		t.Errorf("AUMSigHash is wrong size: got %d, want %d", len(sigHash), blake2s.Size)
	}

	var nksHash NKSSigHash
	if len(nksHash) != blake2s.Size {
		t.Errorf("NKSSigHash is wrong size: got %d, want %d", len(nksHash), blake2s.Size)
	}
}

func TestMarshaledSignatureJSON(t *testing.T) {
	sig := MarshaledSignature("abcdef")
	j, err := json.Marshal(sig)
	if err != nil {
		t.Fatal(err)
	}
	const encoded = `"YWJjZGVm"`
	if string(j) != encoded {
		t.Errorf("got JSON %q; want %q", j, encoded)
	}

	var back MarshaledSignature
	if err := json.Unmarshal([]byte(encoded), &back); err != nil {
		t.Fatal(err)
	}
	if string(back) != string(sig) {
		t.Errorf("decoded JSON back to %q; want %q", back, sig)
	}
}

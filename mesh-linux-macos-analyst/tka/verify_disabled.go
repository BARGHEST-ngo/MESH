// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build ts_omit_tailnetlock

package tka

import (
	"errors"

	"tailscale.com/types/tkatype"
)

// signatureVerify returns a nil error if the signature is valid over the
// provided AUM BLAKE2s digest, using the given key.
func signatureVerify(s *tkatype.Signature, aumDigest tkatype.AUMSigHash, key Key) error {
	return errors.New("tailnetlock disabled in build")
}

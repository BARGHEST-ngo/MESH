// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !js && !ts_omit_acme

package tailscale

import (
	"context"
	"crypto/tls"

	"tailscale.com/client/local"
)

// GetCertificate is an alias for [tailscale.com/client/local.GetCertificate].
//
// Deprecated: import [tailscale.com/client/local] instead and use [local.Client.GetCertificate].
func GetCertificate(hi *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return local.GetCertificate(hi)
}

// CertPair is an alias for [tailscale.com/client/local.CertPair].
//
// Deprecated: import [tailscale.com/client/local] instead and use [local.Client.CertPair].
func CertPair(ctx context.Context, domain string) (certPEM, keyPEM []byte, err error) {
	return local.CertPair(ctx, domain)
}

// ExpandSNIName is an alias for [tailscale.com/client/local.ExpandSNIName].
//
// Deprecated: import [tailscale.com/client/local] instead and use [local.Client.ExpandSNIName].
func ExpandSNIName(ctx context.Context, name string) (fqdn string, ok bool) {
	return local.ExpandSNIName(ctx, name)
}

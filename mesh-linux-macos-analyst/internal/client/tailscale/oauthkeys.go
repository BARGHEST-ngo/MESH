// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package tailscale

import (
	"context"

	"tailscale.com/feature"
)

// HookResolveAuthKey resolves to [oauthkey.ResolveAuthKey] when the
// corresponding feature tag is enabled in the build process.
//
// authKey is a standard device auth key or an OAuth client secret to
// resolve into an auth key.
// tags is the list of tags being advertised by the client (required to be
// provided for the OAuth secret case, and required to be the same as the
// list of tags for which the OAuth secret is allowed to issue auth keys).
var HookResolveAuthKey feature.Hook[func(ctx context.Context, authKey string, tags []string) (string, error)]

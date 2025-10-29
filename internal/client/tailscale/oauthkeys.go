// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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

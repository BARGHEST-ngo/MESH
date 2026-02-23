// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package cli

import (
	"context"
	"errors"

	"github.com/peterbourgon/ff/v3/ffcli"
	"tailscale.com/envknob"
)

var idTokenCmd = &ffcli.Command{
	Name:       "id-token",
	ShortUsage: "tailscale id-token <aud>",
	ShortHelp:  "Fetch an OIDC id-token for the Tailscale machine",
	LongHelp:   hidden,
	Exec:       runIDToken,
}

func runIDToken(ctx context.Context, args []string) error {
	if !envknob.UseWIPCode() {
		return errors.New("tailscale id-token: works-in-progress require TAILSCALE_USE_WIP_CODE=1 envvar")
	}
	if len(args) != 1 {
		return errors.New("usage: tailscale id-token <aud>")
	}

	tr, err := localClient.IDToken(ctx, args[0])
	if err != nil {
		return err
	}

	outln(tr.IDToken)
	return nil
}

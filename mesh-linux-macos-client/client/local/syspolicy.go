// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !ts_omit_syspolicy

package local

import (
	"context"
	"net/http"

	"tailscale.com/util/syspolicy/setting"
)

// GetEffectivePolicy returns the effective policy for the specified scope.
func (lc *Client) GetEffectivePolicy(ctx context.Context, scope setting.PolicyScope) (*setting.Snapshot, error) {
	scopeID, err := scope.MarshalText()
	if err != nil {
		return nil, err
	}
	body, err := lc.get200(ctx, "/localapi/v0/policy/"+string(scopeID))
	if err != nil {
		return nil, err
	}
	return decodeJSON[*setting.Snapshot](body)
}

// ReloadEffectivePolicy reloads the effective policy for the specified scope
// by reading and merging policy settings from all applicable policy sources.
func (lc *Client) ReloadEffectivePolicy(ctx context.Context, scope setting.PolicyScope) (*setting.Snapshot, error) {
	scopeID, err := scope.MarshalText()
	if err != nil {
		return nil, err
	}
	body, err := lc.send(ctx, "POST", "/localapi/v0/policy/"+string(scopeID), 200, http.NoBody)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*setting.Snapshot](body)
}

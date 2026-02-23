// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package dns

type noopManager struct{}

func (m noopManager) SetDNS(OSConfig) error  { return nil }
func (m noopManager) SupportsSplitDNS() bool { return false }
func (m noopManager) Close() error           { return nil }
func (m noopManager) GetBaseConfig() (OSConfig, error) {
	return OSConfig{}, ErrGetBaseConfigNotSupported
}

func NewNoopManager() (noopManager, error) {
	return noopManager{}, nil
}

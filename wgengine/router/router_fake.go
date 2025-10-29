// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package router

import (
	"tailscale.com/types/logger"
)

// NewFake returns a Router that does nothing when called and always
// returns nil errors.
func NewFake(logf logger.Logf) Router {
	return fakeRouter{logf: logf}
}

type fakeRouter struct {
	logf logger.Logf
}

func (r fakeRouter) Up() error {
	r.logf("[v1] warning: fakeRouter.Up: not implemented.")
	return nil
}

func (r fakeRouter) Set(cfg *Config) error {
	r.logf("[v1] warning: fakeRouter.Set: not implemented.")
	return nil
}

func (r fakeRouter) Close() error {
	r.logf("[v1] warning: fakeRouter.Close: not implemented.")
	return nil
}

// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build ts_omit_usermetrics

package usermetric

type Registry struct {
	m Metrics
}

func (*Registry) NewGauge(name, help string) *Gauge { return nil }

type MultiLabelMap[T comparable] = noopMap[T]

type noopMap[T comparable] struct{}

type Gauge struct{}

func (*Gauge) Set(float64) {}

func NewMultiLabelMapWithRegistry[T comparable](m *Registry, name string, promType, helpText string) *MultiLabelMap[T] {
	return nil
}

func (*noopMap[T]) Add(T, int64) {}
func (*noopMap[T]) Set(T, any)   {}

func (r *Registry) Handler(any, any) {} // no-op HTTP handler

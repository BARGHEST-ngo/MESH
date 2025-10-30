// androidqf - Android Quick Forensics
// Copyright (c) 2021-2022 Claudio Guarnieri.
// Use of this software is governed by the MVT License 1.1 that can be found at
//   https://license.mvt.re/1.1/

package assets

import (
	"errors"
)

// CollectorFS is a stub for the collector filesystem
// We don't bundle collector binaries - users should rely on system adb commands
type CollectorFS struct{}

// ReadFile returns an error since we don't bundle collector binaries
func (c CollectorFS) ReadFile(name string) ([]byte, error) {
	return nil, errors.New("collector binaries not bundled - please use system adb commands")
}

// Collector is a stub instance since we don't bundle collector binaries
var Collector = CollectorFS{}

type Asset struct {
	Name string
	Data []byte
}

// DeployAssets is a no-op since we rely on system adb instead of embedded binaries
func DeployAssets() error {
	// No-op: we use system adb instead of embedded binaries
	return nil
}

// CleanAssets is a no-op since we don't deploy any assets
func CleanAssets() error {
	// No-op: we don't deploy any assets
	return nil
}

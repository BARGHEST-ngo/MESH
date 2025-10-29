// androidqf - Android Quick Forensics
// Copyright (c) 2021-2022 Claudio Guarnieri.
// Use of this software is governed by the MVT License 1.1 that can be found at
//   https://license.mvt.re/1.1/

package adb

import (
	"os/exec"
	"path/filepath"

	saveRuntime "github.com/botherder/go-savetime/runtime"
)

func (a *ADB) findExe() error {
	// Skip asset deployment and use system adb
	// err := assets.DeployAssets()
	// if err != nil {
	// 	return err
	// }

	adbPath, err := exec.LookPath("adb")
	if err == nil {
		a.ExePath = adbPath
	} else {
		a.ExePath = filepath.Join(saveRuntime.GetExecutableDirectory(), "adb")
	}
	return nil
}

// Copyright (c) 2021-2023 Claudio Guarnieri.
// Use of this source code is governed by the MVT License 1.1
// which can be found in the LICENSE file.

package modules

import (
	"fmt"

	"github.com/mvt-project/androidqf_ward/acquisition"
	"github.com/mvt-project/androidqf_ward/adb"
	"github.com/mvt-project/androidqf_ward/log"
)

type Environment struct {
	StoragePath string
}

func NewEnvironment() *Environment {
	return &Environment{}
}

func (e *Environment) Name() string {
	return "environment"
}

func (e *Environment) InitStorage(storagePath string) error {
	e.StoragePath = storagePath
	return nil
}

func (e *Environment) Run(acq *acquisition.Acquisition, fast bool) error {
	log.Info("Collecting environment...")

	out, err := adb.Client.Shell("env")
	if err != nil {
		return fmt.Errorf("failed to run `adb shell env`: %v", err)
	}

	return saveStringToAcquisition(acq, "env.txt", out)
}

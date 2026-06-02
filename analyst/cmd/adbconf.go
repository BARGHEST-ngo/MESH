// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type adbConfig struct {
	Host     string `json:"host"`
	Hostport string `json:"hostport"`
	Oem      string `json:"oem"`
}

var (
	adbConf     adbConfig
	adbConfOnce sync.Once
)

const adbConfFilename = "adbconf.json"

func ensureAdbConf() {
	adbConfOnce.Do(func() {
		if _, err := os.Stat(adbConfFilename); os.IsNotExist(err) {
			f, err := os.Create(adbConfFilename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "adbconf create: %v\n", err)
				return
			}
			defer f.Close()
			b, _ := json.Marshal(adbConfig{})
			_, _ = f.Write(b)
		}

		data, err := os.ReadFile(adbConfFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "adbconf read: %v\n", err)
			return
		}

		var cfg adbConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			fmt.Fprintf(os.Stderr, "error unmarshalling adbConfig: %v\n", err)
			return
		}
		adbConf = cfg
	})
}

func saveHost(host string) error {
	ensureAdbConf()
	adbConf.Host = host
	return writeAdbConf()
}

func saveHostport(port string) error {
	ensureAdbConf()
	adbConf.Hostport = port
	return writeAdbConf()
}

func saveOem(oem string) error {
	ensureAdbConf()
	adbConf.Oem = oem
	return writeAdbConf()
}

func writeAdbConf() error {
	data, err := json.MarshalIndent(adbConf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(adbConfFilename, data, 0644)
}

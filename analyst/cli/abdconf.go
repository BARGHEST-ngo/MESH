// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cli

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host     string `json:"host"`
	Hostport string `json:"hostport"`
	Oem      string `json:"oem"`
}

var config Config

const filename = "adbconf.json"

func main() {
	if err := checkFile(filename); err != nil {
		errf("checkFile: %v\n", err)
		return
	}
	loadConfig()
}

func checkFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		_, _ = f.Write([]byte("[{}]"))
	}
	return nil
}

func loadConfig() {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		errf("read: %v\n", err)
		return
	}
	var cfgs []Config
	if err := json.Unmarshal(bytes, &cfgs); err == nil && len(cfgs) > 0 {
		config = cfgs[0]
	}
}

func saveHost(host string) error {
	config.Host = host
	return writeConfig()
}

func saveHostport(port string) error {
	config.Hostport = port
	return writeConfig()
}

func saveOem(oem string) error {
	config.Oem = oem
	return writeConfig()
}

func writeConfig() error {
	bytes, err := json.MarshalIndent([]Config{config}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, 0644)
}

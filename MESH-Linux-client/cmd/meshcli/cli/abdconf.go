// SPDX-License-Identifier: CC0-1.0
// Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.

package cli

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("checkFile:", err)
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
		fmt.Println("read:", err)
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

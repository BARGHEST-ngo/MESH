// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package wgcfg

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"tailscale.com/types/logger"
)

type AmneziaConfig struct {
	Jc   uint8
	Jmin uint16
	Jmax uint16
	S1   uint16
	S2   uint16
	H1   uint32
	H2   uint32
	H3   uint32
	H4   uint32
}

func DefaultAmneziaConfig() AmneziaConfig {
	return AmneziaConfig{
		Jc:   0,
		Jmin: 0,
		Jmax: 0,
		S1:   0,
		S2:   0,
		H1:   1,
		H2:   2,
		H3:   3,
		H4:   4,
	}
}

var getConfigPath = func() string {
	switch runtime.GOOS {
	case "linux":
		return "/etc/mesh/amneziawg.conf"
	case "darwin":
		return "/usr/local/etc/mesh/amneziawg.conf"
	case "windows":
		return "C:\\ProgramData\\MESH\\amneziawg.conf"
	case "android":
		return "/data/data/com.barghest.mesh/files/amneziawg.conf"
	default:
		return "/etc/mesh/amneziawg.conf"
	}
}

func LoadAmneziaConfig(logf logger.Logf) AmneziaConfig {
	cfg := DefaultAmneziaConfig()

	path := getConfigPath()
	file, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			logf("amneziawg: failed to open config file %s: %v", path, err)
		} else {
			logf("amneziawg: config file not found at %s, using defaults (backward compatible mode)", path)
		}
		return cfg
	}
	defer file.Close()

	logf("amneziawg: loading config from %s", path)

	scanner := bufio.NewScanner(file)
	inInterface := false
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section := strings.ToLower(strings.Trim(line, "[]"))
			inInterface = (section == "interface")
			continue
		}

		if !inInterface {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			logf("amneziawg: invalid line %d in %s: %s", lineNum, path, line)
			continue
		}

		key := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])

		switch key {
		case "jc":
			if v, err := strconv.ParseUint(value, 10, 8); err == nil {
				cfg.Jc = uint8(v)
			} else {
				logf("amneziawg: invalid Jc value at line %d: %v", lineNum, err)
			}
		case "jmin":
			if v, err := strconv.ParseUint(value, 10, 16); err == nil {
				cfg.Jmin = uint16(v)
			} else {
				logf("amneziawg: invalid Jmin value at line %d: %v", lineNum, err)
			}
		case "jmax":
			if v, err := strconv.ParseUint(value, 10, 16); err == nil {
				cfg.Jmax = uint16(v)
			} else {
				logf("amneziawg: invalid Jmax value at line %d: %v", lineNum, err)
			}
		case "s1":
			if v, err := strconv.ParseUint(value, 10, 16); err == nil {
				cfg.S1 = uint16(v)
			} else {
				logf("amneziawg: invalid S1 value at line %d: %v", lineNum, err)
			}
		case "s2":
			if v, err := strconv.ParseUint(value, 10, 16); err == nil {
				cfg.S2 = uint16(v)
			} else {
				logf("amneziawg: invalid S2 value at line %d: %v", lineNum, err)
			}
		case "h1":
			if v, err := strconv.ParseUint(value, 10, 32); err == nil {
				cfg.H1 = uint32(v)
			} else {
				logf("amneziawg: invalid H1 value at line %d: %v", lineNum, err)
			}
		case "h2":
			if v, err := strconv.ParseUint(value, 10, 32); err == nil {
				cfg.H2 = uint32(v)
			} else {
				logf("amneziawg: invalid H2 value at line %d: %v", lineNum, err)
			}
		case "h3":
			if v, err := strconv.ParseUint(value, 10, 32); err == nil {
				cfg.H3 = uint32(v)
			} else {
				logf("amneziawg: invalid H3 value at line %d: %v", lineNum, err)
			}
		case "h4":
			if v, err := strconv.ParseUint(value, 10, 32); err == nil {
				cfg.H4 = uint32(v)
			} else {
				logf("amneziawg: invalid H4 value at line %d: %v", lineNum, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logf("amneziawg: error reading config file %s: %v", path, err)
	}

	if err := cfg.Validate(); err != nil {
		logf("amneziawg: invalid configuration: %v, using defaults", err)
		return DefaultAmneziaConfig()
	}

	// Log the loaded configuration
	if cfg.IsObfuscationEnabled() {
		logf("amneziawg: obfuscation ENABLED - Jc=%d Jmin=%d Jmax=%d S1=%d S2=%d H1=%d H2=%d H3=%d H4=%d",
			cfg.Jc, cfg.Jmin, cfg.Jmax, cfg.S1, cfg.S2, cfg.H1, cfg.H2, cfg.H3, cfg.H4)
	} else {
		logf("amneziawg: obfuscation DISABLED (backward compatible mode)")
	}

	return cfg
}

func (c *AmneziaConfig) Validate() error {
	if c.Jc > 128 {
		return fmt.Errorf("Jc must be between 0 and 128, got %d", c.Jc)
	}

	if c.Jmax > 1280 {
		return fmt.Errorf("Jmax must be <= 1280, got %d", c.Jmax)
	}

	if c.Jmin > c.Jmax && c.Jmax != 0 {
		return fmt.Errorf("Jmin (%d) must be <= Jmax (%d)", c.Jmin, c.Jmax)
	}

	if c.S1 > 0 && c.S1 < 15 {
		return fmt.Errorf("S1 must be 0 or >= 15, got %d", c.S1)
	}

	if c.S1 > 150 {
		return fmt.Errorf("S1 must be <= 150, got %d", c.S1)
	}

	if c.S2 > 0 && c.S2 < 15 {
		return fmt.Errorf("S2 must be 0 or >= 15, got %d", c.S2)
	}

	if c.S2 > 150 {
		return fmt.Errorf("S2 must be <= 150, got %d", c.S2)
	}

	if c.S1 > 0 && c.S2 > 0 && c.S1+56 == c.S2 {
		return fmt.Errorf("S1+56 must not equal S2 (S1=%d, S2=%d)", c.S1, c.S2)
	}

	if c.H1 < 1 {
		return fmt.Errorf("H1 must be >= 1, got %d", c.H1)
	}

	if c.H2 < 1 {
		return fmt.Errorf("H2 must be >= 1, got %d", c.H2)
	}

	if c.H3 < 1 {
		return fmt.Errorf("H3 must be >= 1, got %d", c.H3)
	}

	if c.H4 < 1 {
		return fmt.Errorf("H4 must be >= 1, got %d", c.H4)
	}

	hValues := map[uint32]bool{c.H1: true, c.H2: true, c.H3: true, c.H4: true}
	if len(hValues) != 4 {
		return fmt.Errorf("H1, H2, H3, H4 must all be different values")
	}

	return nil
}

// IsObfuscationEnabled returns true if any obfuscation parameters are non-default
func (c *AmneziaConfig) IsObfuscationEnabled() bool {
	defaults := DefaultAmneziaConfig()
	return c.Jc != defaults.Jc ||
		c.Jmin != defaults.Jmin ||
		c.Jmax != defaults.Jmax ||
		c.S1 != defaults.S1 ||
		c.S2 != defaults.S2 ||
		c.H1 != defaults.H1 ||
		c.H2 != defaults.H2 ||
		c.H3 != defaults.H3 ||
		c.H4 != defaults.H4
}

func (c *AmneziaConfig) ApplyTo(cfg *Config) {
	cfg.Jc = c.Jc
	cfg.Jmin = c.Jmin
	cfg.Jmax = c.Jmax
	cfg.S1 = c.S1
	cfg.S2 = c.S2
	cfg.H1 = c.H1
	cfg.H2 = c.H2
	cfg.H3 = c.H3
	cfg.H4 = c.H4
}

// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package wgcfg

import (
	"strings"
	"testing"

	"tailscale.com/types/key"
)

// TestAWGParametersInUAPI verifies that AWG parameters are correctly written to UAPI
func TestAWGParametersInUAPI(t *testing.T) {
	// Create a config with AWG parameters
	cfg := &Config{
		Name:       "test",
		PrivateKey: key.NewNode(),
		Jc:         5,
		Jmin:       50,
		Jmax:       1000,
		S1:         30,
		S2:         40,
		H1:         100,
		H2:         200,
		H3:         300,
		H4:         400,
	}

	// Previous config with defaults (simulating first config)
	prev := &Config{
		Name:       "test",
		PrivateKey: cfg.PrivateKey,
		Jc:         0,
		Jmin:       0,
		Jmax:       0,
		S1:         0,
		S2:         0,
		H1:         1,
		H2:         2,
		H3:         3,
		H4:         4,
	}

	// Write to UAPI
	buf := new(strings.Builder)
	err := cfg.ToUAPI(t.Logf, buf, prev)
	if err != nil {
		t.Fatalf("ToUAPI failed: %v", err)
	}

	uapi := buf.String()
	t.Logf("UAPI output:\n%s", uapi)

	// Verify all AWG parameters are present
	requiredParams := map[string]string{
		"jc":   "5",
		"jmin": "50",
		"jmax": "1000",
		"s1":   "30",
		"s2":   "40",
		"h1":   "100",
		"h2":   "200",
		"h3":   "300",
		"h4":   "400",
	}

	for param, expectedValue := range requiredParams {
		expected := param + "=" + expectedValue
		if !strings.Contains(uapi, expected) {
			t.Errorf("UAPI output missing %q, got:\n%s", expected, uapi)
		}
	}
}

// TestAWGParametersNotWrittenWhenUnchanged verifies that AWG parameters are NOT written when unchanged
func TestAWGParametersNotWrittenWhenUnchanged(t *testing.T) {
	// Create a config with AWG parameters
	cfg := &Config{
		Name:       "test",
		PrivateKey: key.NewNode(),
		Jc:         5,
		Jmin:       50,
		Jmax:       1000,
		S1:         30,
		S2:         40,
		H1:         100,
		H2:         200,
		H3:         300,
		H4:         400,
	}

	// Previous config with SAME AWG parameters
	prev := &Config{
		Name:       "test",
		PrivateKey: cfg.PrivateKey,
		Jc:         5,
		Jmin:       50,
		Jmax:       1000,
		S1:         30,
		S2:         40,
		H1:         100,
		H2:         200,
		H3:         300,
		H4:         400,
	}

	// Write to UAPI
	buf := new(strings.Builder)
	err := cfg.ToUAPI(t.Logf, buf, prev)
	if err != nil {
		t.Fatalf("ToUAPI failed: %v", err)
	}

	uapi := buf.String()
	t.Logf("UAPI output:\n%s", uapi)

	// Verify AWG parameters are NOT present (since they haven't changed)
	awgParams := []string{"jc=", "jmin=", "jmax=", "s1=", "s2=", "h1=", "h2=", "h3=", "h4="}
	for _, param := range awgParams {
		if strings.Contains(uapi, param) {
			t.Errorf("UAPI output should not contain %q when unchanged, but got:\n%s", param, uapi)
		}
	}
}


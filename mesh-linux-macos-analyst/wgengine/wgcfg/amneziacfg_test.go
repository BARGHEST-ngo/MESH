// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package wgcfg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultAmneziaConfig(t *testing.T) {
	cfg := DefaultAmneziaConfig()
	
	if cfg.Jc != 0 {
		t.Errorf("Expected Jc=0, got %d", cfg.Jc)
	}
	if cfg.Jmin != 0 {
		t.Errorf("Expected Jmin=0, got %d", cfg.Jmin)
	}
	if cfg.Jmax != 0 {
		t.Errorf("Expected Jmax=0, got %d", cfg.Jmax)
	}
	if cfg.S1 != 0 {
		t.Errorf("Expected S1=0, got %d", cfg.S1)
	}
	if cfg.S2 != 0 {
		t.Errorf("Expected S2=0, got %d", cfg.S2)
	}
	if cfg.H1 != 1 {
		t.Errorf("Expected H1=1, got %d", cfg.H1)
	}
	if cfg.H2 != 2 {
		t.Errorf("Expected H2=2, got %d", cfg.H2)
	}
	if cfg.H3 != 3 {
		t.Errorf("Expected H3=3, got %d", cfg.H3)
	}
	if cfg.H4 != 4 {
		t.Errorf("Expected H4=4, got %d", cfg.H4)
	}
}

func TestAmneziaConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     AmneziaConfig
		wantErr bool
	}{
		{
			name:    "default config is valid",
			cfg:     DefaultAmneziaConfig(),
			wantErr: false,
		},
		{
			name: "valid obfuscation config",
			cfg: AmneziaConfig{
				Jc: 5, Jmin: 50, Jmax: 1000,
				S1: 30, S2: 40,
				H1: 100, H2: 200, H3: 300, H4: 400,
			},
			wantErr: false,
		},
		{
			name: "Jc too high",
			cfg: AmneziaConfig{
				Jc: 129, Jmin: 50, Jmax: 1000,
				S1: 30, S2: 40,
				H1: 100, H2: 200, H3: 300, H4: 400,
			},
			wantErr: true,
		},
		{
			name: "Jmax too high",
			cfg: AmneziaConfig{
				Jc: 5, Jmin: 50, Jmax: 1281,
				S1: 30, S2: 40,
				H1: 100, H2: 200, H3: 300, H4: 400,
			},
			wantErr: true,
		},
		{
			name: "Jmin > Jmax",
			cfg: AmneziaConfig{
				Jc: 5, Jmin: 1000, Jmax: 50,
				S1: 30, S2: 40,
				H1: 100, H2: 200, H3: 300, H4: 400,
			},
			wantErr: true,
		},
		{
			name: "S1 too low",
			cfg: AmneziaConfig{
				Jc: 5, Jmin: 50, Jmax: 1000,
				S1: 10, S2: 40,
				H1: 100, H2: 200, H3: 300, H4: 400,
			},
			wantErr: true,
		},
		{
			name: "S1 too high",
			cfg: AmneziaConfig{
				Jc: 5, Jmin: 50, Jmax: 1000,
				S1: 151, S2: 40,
				H1: 100, H2: 200, H3: 300, H4: 400,
			},
			wantErr: true,
		},
		{
			name: "S1+56 == S2",
			cfg: AmneziaConfig{
				Jc: 5, Jmin: 50, Jmax: 1000,
				S1: 30, S2: 86,
				H1: 100, H2: 200, H3: 300, H4: 400,
			},
			wantErr: true,
		},
		{
			name: "H1 is zero",
			cfg: AmneziaConfig{
				Jc: 5, Jmin: 50, Jmax: 1000,
				S1: 30, S2: 40,
				H1: 0, H2: 200, H3: 300, H4: 400,
			},
			wantErr: true,
		},
		{
			name: "duplicate H values",
			cfg: AmneziaConfig{
				Jc: 5, Jmin: 50, Jmax: 1000,
				S1: 30, S2: 40,
				H1: 100, H2: 100, H3: 300, H4: 400,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadAmneziaConfigFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "amneziawg.conf")

	configContent := `# Test config
[Interface]
Jc = 5
Jmin = 50
Jmax = 1000
S1 = 30
S2 = 40
H1 = 100
H2 = 200
H3 = 300
H4 = 400
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	oldGetConfigPath := getConfigPath
	defer func() { getConfigPath = oldGetConfigPath }()
	
	getConfigPath = func() string { return configPath }

	logf := t.Logf
	cfg := LoadAmneziaConfig(logf)

	if cfg.Jc != 5 {
		t.Errorf("Expected Jc=5, got %d", cfg.Jc)
	}
	if cfg.Jmin != 50 {
		t.Errorf("Expected Jmin=50, got %d", cfg.Jmin)
	}
	if cfg.Jmax != 1000 {
		t.Errorf("Expected Jmax=1000, got %d", cfg.Jmax)
	}
	if cfg.S1 != 30 {
		t.Errorf("Expected S1=30, got %d", cfg.S1)
	}
	if cfg.S2 != 40 {
		t.Errorf("Expected S2=40, got %d", cfg.S2)
	}
	if cfg.H1 != 100 {
		t.Errorf("Expected H1=100, got %d", cfg.H1)
	}
	if cfg.H2 != 200 {
		t.Errorf("Expected H2=200, got %d", cfg.H2)
	}
	if cfg.H3 != 300 {
		t.Errorf("Expected H3=300, got %d", cfg.H3)
	}
	if cfg.H4 != 400 {
		t.Errorf("Expected H4=400, got %d", cfg.H4)
	}
}

func TestApplyTo(t *testing.T) {
	amneziaCfg := AmneziaConfig{
		Jc: 5, Jmin: 50, Jmax: 1000,
		S1: 30, S2: 40,
		H1: 100, H2: 200, H3: 300, H4: 400,
	}

	cfg := &Config{}
	amneziaCfg.ApplyTo(cfg)

	if cfg.Jc != 5 {
		t.Errorf("Expected Jc=5, got %d", cfg.Jc)
	}
	if cfg.Jmin != 50 {
		t.Errorf("Expected Jmin=50, got %d", cfg.Jmin)
	}
	if cfg.Jmax != 1000 {
		t.Errorf("Expected Jmax=1000, got %d", cfg.Jmax)
	}
	if cfg.S1 != 30 {
		t.Errorf("Expected S1=30, got %d", cfg.S1)
	}
	if cfg.S2 != 40 {
		t.Errorf("Expected S2=40, got %d", cfg.S2)
	}
	if cfg.H1 != 100 {
		t.Errorf("Expected H1=100, got %d", cfg.H1)
	}
	if cfg.H2 != 200 {
		t.Errorf("Expected H2=200, got %d", cfg.H2)
	}
	if cfg.H3 != 300 {
		t.Errorf("Expected H3=300, got %d", cfg.H3)
	}
	if cfg.H4 != 400 {
		t.Errorf("Expected H4=400, got %d", cfg.H4)
	}
}


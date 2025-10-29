//
// Copyright 2024 Stacklok, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package config provides the frizbee configuration.
package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type contextConfigKey struct{}

// ContextConfigKey is the context key for the configuration.
// nolint:gochecknoglobals // this is a context key
var ContextConfigKey = contextConfigKey{}

var (
	// ErrNoConfigInContext is returned when no configuration is found in the context.
	ErrNoConfigInContext = errors.New("no configuration found in context")
)

// FromCommand returns the configuration from the cobra command.
func FromCommand(cmd *cobra.Command) (*Config, error) {
	ctx := cmd.Context()
	cfg, ok := ctx.Value(ContextConfigKey).(*Config)
	if !ok {
		return nil, ErrNoConfigInContext
	}

	// If the platform flag is set, override the platform in the configuration.
	if cmd.Flags().Lookup("platform") != nil {
		cfg.Platform = cmd.Flag("platform").Value.String()
	}
	return cfg, nil
}

// Config is the frizbee configuration.
type Config struct {
	Platform  string    `yaml:"platform" mapstructure:"platform"`
	GHActions GHActions `yaml:"ghactions" mapstructure:"ghactions"`
	Images    Images    `yaml:"images" mapstructure:"images"`
}

// GHActions is the GitHub Actions configuration.
type GHActions struct {
	Filter `yaml:",inline" mapstructure:",inline"`
}

// Filter is a common configuration for filtering out patterns.
type Filter struct {
	// Exclude is a list of patterns to exclude.
	Exclude         []string `yaml:"exclude" mapstructure:"exclude"`
	ExcludeBranches []string `yaml:"exclude_branches" mapstructure:"exclude_branches"`
}

// Images is the image configuration.
type Images struct {
	ImageFilter `yaml:",inline" mapstructure:",inline"`
}

// ImageFilter is the image filter configuration.
type ImageFilter struct {
	// ExcludeImages is a regex that must match in order for an image to be excluded and not pinned
	ExcludeImages []string `yaml:"exclude_images" mapstructure:"exclude_images"`
	ExcludeTags   []string `yaml:"exclude_tags" mapstructure:"exclude_tags"`
}

// ParseConfigFile parses a configuration file.
func ParseConfigFile(configfile string) (*Config, error) {
	bfs := osfs.New(".")
	return ParseConfigFileFromFS(bfs, configfile)
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		GHActions: GHActions{
			Filter: Filter{
				ExcludeBranches: []string{"main", "master"},
			},
		},
		Images: Images{
			ImageFilter: ImageFilter{
				ExcludeImages: []string{"scratch"},
				ExcludeTags:   []string{"latest"},
			},
		},
	}
}

// MergeUserConfig merges the user configuration with the default configuration.
// mostly making sure that we don't try to pin the scratch image
func MergeUserConfig(userConfig *Config) *Config {
	if userConfig == nil {
		return DefaultConfig()
	}

	if userConfig.Images.ExcludeImages == nil {
		userConfig.Images.ExcludeImages = []string{"scratch"}
	}

	if !slices.Contains(userConfig.Images.ExcludeImages, "scratch") {
		userConfig.Images.ExcludeImages = append(userConfig.Images.ExcludeImages, "scratch")
	}

	return userConfig
}

// ParseConfigFileFromFS parses a configuration file from a filesystem.
func ParseConfigFileFromFS(fs billy.Filesystem, configfile string) (*Config, error) {
	cfg := DefaultConfig()
	cleancfgfile := filepath.Clean(configfile)
	cfgF, err := fs.Open(cleancfgfile)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}

		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer cfgF.Close() // nolint:errcheck

	dec := yaml.NewDecoder(cfgF)
	if err := dec.Decode(cfg); err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("failed to decode config file: %w", err)
		}
	}

	return cfg, nil
}

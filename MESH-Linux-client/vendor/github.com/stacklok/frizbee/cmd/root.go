//
// Copyright 2023 Stacklok, Inc.
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

// Package cmd provides the frizbee command line interface.
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/stacklok/frizbee/cmd/actions"
	"github.com/stacklok/frizbee/cmd/image"
	"github.com/stacklok/frizbee/cmd/version"
	"github.com/stacklok/frizbee/pkg/utils/config"
)

// Execute runs the root command.
func Execute() {
	var rootCmd = &cobra.Command{
		Use:               "frizbee",
		Short:             "frizbee is a tool you may throw a tag at and it comes back with a checksum",
		PersistentPreRunE: prerun,
	}

	rootCmd.PersistentFlags().StringP("config", "c", ".frizbee.yml", "config file (default is .frizbee.yml)")

	rootCmd.AddCommand(actions.CmdGHActions())
	rootCmd.AddCommand(image.CmdContainerImage())
	rootCmd.AddCommand(version.CmdVersion())

	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}

func prerun(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	cfg, err := readConfig(cmd)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	ctx = context.WithValue(ctx, config.ContextConfigKey, cfg)

	cmd.SetContext(ctx)

	return nil
}

func readConfig(cmd *cobra.Command) (*config.Config, error) {
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, fmt.Errorf("failed to get config file: %w", err)
	}

	return config.ParseConfigFile(configFile)
}

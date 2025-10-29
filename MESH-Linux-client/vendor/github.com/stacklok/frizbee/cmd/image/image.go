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

// Package image provides command-line utilities to work with container images.
package image

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/stacklok/frizbee/internal/cli"
	"github.com/stacklok/frizbee/pkg/interfaces"
	"github.com/stacklok/frizbee/pkg/replacer"
	"github.com/stacklok/frizbee/pkg/utils/config"
)

// CmdContainerImage represents the containers command
func CmdContainerImage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Replace container image references with checksums",
		Long: `This utility replaces tag or branch references in yaml/yml files
with the latest commit hash of the referenced tag or branch.
	
Example:

	$ frizbee image <path-to-yaml-files> or <ghcr.io/stacklok/minder/server:latest>

This will replace all tag or branch references in all yaml files for the given directory.
`,
		RunE:         replaceCmd,
		SilenceUsage: true,
		Aliases:      []string{"containerimage", "dockercompose", "compose"}, // backwards compatibility
		Args:         cobra.ExactArgs(1),
	}

	// flags
	cli.DeclareFrizbeeFlags(cmd, false)

	// sub-commands
	cmd.AddCommand(CmdList())

	return cmd
}

func replaceCmd(cmd *cobra.Command, args []string) error {
	// Extract the CLI flags from the cobra command
	cliFlags, err := cli.NewHelper(cmd)
	if err != nil {
		return err
	}

	// Set up the config
	cfg, err := config.FromCommand(cmd)
	if err != nil {
		return err
	}

	// Create a new replacer
	r := replacer.NewContainerImagesReplacer(cfg).
		WithUserRegex(cliFlags.Regex)

	if cli.IsPath(args[0]) {
		dir := filepath.Clean(args[0])
		// Replace the tags in the directory
		res, err := r.ParsePath(cmd.Context(), dir)
		if err != nil {
			return err
		}
		// Process the output files
		return cliFlags.ProcessOutput(dir, res.Processed, res.Modified)
	}
	// Replace the passed reference
	res, err := r.ParseString(cmd.Context(), args[0])
	if err != nil {
		if errors.Is(err, interfaces.ErrReferenceSkipped) {
			fmt.Fprintln(cmd.OutOrStdout(), args[0]) // nolint:errcheck
			return nil
		}
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%s@%s\n", res.Name, res.Ref) // nolint:errcheck
	return nil

}

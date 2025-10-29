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

// Package version adds a version command.
package version

import (
	"github.com/spf13/cobra"

	"github.com/stacklok/frizbee/internal/cli"
)

// CmdVersion is the Cobra command for the version command.
func CmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print frizbee CLI version",
		Long:  "The frizbee version command prints the version of the frizbee CLI.",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Println(cli.VerboseCLIVersion)
		},
	}
}

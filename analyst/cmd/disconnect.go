// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"

	"github.com/BARGHEST-ngo/MESH/analyst/session"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func DisconnectCmd(s *session.Session, in *Input) *ffcli.Command {
	return &ffcli.Command{
		Name:       "disconnect",
		ShortUsage: "disconnect",
		ShortHelp:  "End a connected device's session",
		Exec: func(ctx context.Context, args []string) error {
			d, err := chooseDevice(ctx, s, in)
			if err != nil {
				return err
			}
			return s.Disconnect(d.ID)
		},
	}
}

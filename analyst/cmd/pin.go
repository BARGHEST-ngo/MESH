// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/BARGHEST-ngo/MESH/analyst/session"
	meshnet "github.com/BARGHEST-ngo/MESHnet"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func PinCmd(s *session.Session) *ffcli.Command {
	return &ffcli.Command{
		Name: "pin",
		Exec: func(ctx context.Context, args []string) error {
			pin, err := s.NewPairing()
			switch {
			case errors.Is(err, meshnet.ErrServerFull):
				return errors.New("server at capacity, cannot connect new clients")
			case errors.Is(err, meshnet.ErrPairingExhausted):
				return errors.New("too many failed pairing attempts, restart MESH")
			case err != nil:
				return err
			}

			fmt.Printf("PIN: %s\n", pin)
			return nil
		},
	}
}

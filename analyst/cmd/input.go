// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

type Input struct {
	t  *term.Terminal
	fd int
}

func NewInput(r io.Reader) *Input {
	rw := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	return &Input{
		t:  term.NewTerminal(rw, ""),
		fd: int(os.Stdin.Fd()),
	}
}

func (in *Input) Output() io.Writer {
	return in.t
}

// ReadLine prints msg and returns the next line
// It returns ctx's error on Ctrl-C and io.EOF when input ends
func (in *Input) ReadLine(ctx context.Context, msg string) (string, error) {
	in.t.SetPrompt(msg)
	if old, err := term.MakeRaw(in.fd); err == nil {
		defer term.Restore(in.fd, old)
	}
	line, err := in.t.ReadLine()
	return strings.TrimSpace(line), err
}

// ReadValid asks until validate accepts the answer
func (in *Input) ReadValid(ctx context.Context, msg string, validate func(string) error) (string, error) {
	for {
		answer, err := in.ReadLine(ctx, msg)
		if err != nil {
			return "", err
		}
		if err := validate(answer); err != nil {
			fmt.Printf("Invalid input: %v\n", err)
			continue
		}
		return answer, nil
	}
}

// ReadPort asks until the answer is a TCP port number
func (in *Input) ReadPort(ctx context.Context, msg string) (int, error) {
	answer, err := in.ReadValid(ctx, msg, func(s string) error {
		p, err := strconv.Atoi(s)
		if err != nil || p < 1 || p > 65535 {
			return fmt.Errorf("%q is not a port number (1-65535)", s)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	p, _ := strconv.Atoi(answer)
	return p, nil
}

// Select lists choices and asks until one is picked, returning its index
func (in *Input) Select(ctx context.Context, title string, choices []string) (int, error) {
	fmt.Println(title)
	for i, c := range choices {
		fmt.Printf(" %d) %s\n", i+1, c)
	}
	answer, err := in.ReadValid(ctx, "Choice: ", func(s string) error {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > len(choices) {
			return fmt.Errorf("choose 1-%d", len(choices))
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	n, _ := strconv.Atoi(answer)
	return n - 1, nil
}

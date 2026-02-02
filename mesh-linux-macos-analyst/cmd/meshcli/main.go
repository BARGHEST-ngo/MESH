package main

import (
	"fmt"
	"os"

	"github.com/BARGHEST-ngo/MESH/analyst/cli"
)

func main() {
	args := os.Args[1:]
	if err := cli.Run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

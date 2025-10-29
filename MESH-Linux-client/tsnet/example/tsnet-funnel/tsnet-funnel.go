// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

// The tsnet-funnel server demonstrates how to use tsnet with Funnel.
//
// To use it, generate an auth key from the Tailscale admin panel and
// run the demo with the key:
//
//	TS_AUTHKEY=<yourkey> go run tsnet-funnel.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"tailscale.com/tsnet"
)

func main() {
	flag.Parse()
	s := &tsnet.Server{
		Dir:      "./funnel-demo-config",
		Hostname: "fun",
	}
	defer s.Close()

	ln, err := s.ListenFunnel("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Printf("Listening on https://%v\n", s.CertDomains()[0])

	err = http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body><h1>Hello, internet!</h1>")
	}))
	log.Fatal(err)
}

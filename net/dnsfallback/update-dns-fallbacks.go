// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"tailscale.com/tailcfg"
)

func main() {
	res, err := http.Get("https://login.tailscale.com/derpmap/default")
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		res.Write(os.Stderr)
		os.Exit(1)
	}
	dm := new(tailcfg.DERPMap)
	if err := json.NewDecoder(res.Body).Decode(dm); err != nil {
		log.Fatal(err)
	}
	for rid, r := range dm.Regions {
		// Names misleading to check into git, as this is a
		// static snapshot and doesn't reflect the live DERP
		// map.
		r.RegionCode = fmt.Sprintf("r%d", rid)
		r.RegionName = r.RegionCode
	}
	out, err := json.MarshalIndent(dm, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("dns-fallback-servers.json", out, 0644); err != nil {
		log.Fatal(err)
	}
}

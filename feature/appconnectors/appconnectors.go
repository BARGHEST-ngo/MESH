// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

// Package appconnectors registers support for Tailscale App Connectors.
package appconnectors

import (
	"encoding/json"
	"net/http"

	"tailscale.com/ipn/ipnlocal"
	"tailscale.com/tailcfg"
)

func init() {
	ipnlocal.RegisterC2N("GET /appconnector/routes", handleC2NAppConnectorDomainRoutesGet)
}

// handleC2NAppConnectorDomainRoutesGet handles returning the domains
// that the app connector is responsible for, as well as the resolved
// IP addresses for each domain. If the node is not configured as
// an app connector, an empty map is returned.
func handleC2NAppConnectorDomainRoutesGet(b *ipnlocal.LocalBackend, w http.ResponseWriter, r *http.Request) {
	logf := b.Logger()
	logf("c2n: GET /appconnector/routes received")

	var res tailcfg.C2NAppConnectorDomainRoutesResponse
	appConnector := b.AppConnector()
	if appConnector == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Domains = appConnector.DomainRoutes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

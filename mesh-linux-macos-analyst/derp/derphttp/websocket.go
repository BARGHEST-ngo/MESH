// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build js || ((linux || darwin) && ts_debug_websockets)

package derphttp

import (
	"context"
	"log"
	"net"

	"github.com/coder/websocket"
	"tailscale.com/net/wsconn"
)

const canWebsockets = true

func init() {
	dialWebsocketFunc = dialWebsocket
}

func dialWebsocket(ctx context.Context, urlStr string) (net.Conn, error) {
	c, res, err := websocket.Dial(ctx, urlStr, &websocket.DialOptions{
		Subprotocols: []string{"derp"},
	})
	if err != nil {
		log.Printf("websocket Dial: %v, %+v", err, res)
		return nil, err
	}
	log.Printf("websocket: connected to %v", urlStr)
	netConn := wsconn.NetConn(context.Background(), c, websocket.MessageBinary, urlStr)
	return netConn, nil
}

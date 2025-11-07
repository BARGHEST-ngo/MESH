// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package libtailscale

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"tailscale.com/ipn"
)

func (app *App) WatchNotifications(mask int, cb NotificationCallback) NotificationManager {
	log.Printf("WatchNotifications: begin app=%p", app)
	app.ready.Wait()
	be := app.backend
	log.Printf("WatchNotifications: after ready.Wait app=%p backend=%p", app, be)
	// Defensive: if backend is unexpectedly nil, wait briefly to avoid a crash
	if be == nil {
		deadline := time.Now().Add(5 * time.Second)
		for be == nil && time.Now().Before(deadline) {
			time.Sleep(50 * time.Millisecond)
			be = app.backend
		}
		if be == nil {
			log.Printf("WatchNotifications: backend still nil after wait; returning no-op manager to avoid crash")
			ctx, cancel := context.WithCancel(context.Background())
			// Hold a goroutine so the manager is well-formed; it exits on Stop().
			go func() { <-ctx.Done() }()
			return &notificationManager{cancel}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	go be.WatchNotifications(ctx, ipn.NotifyWatchOpt(mask), func() {}, func(notify *ipn.Notify) bool {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("panic in WatchNotifications %s: %s", p, debug.Stack())
				panic(p)
			}
		}()

		b, err := json.Marshal(notify)
		if err != nil {
			log.Printf("error: WatchNotifications: marshal notify: %s", err)
			return true
		}
		if err = cb.OnNotify(b); err != nil {
			log.Printf("error: WatchNotifications: OnNotify: %s", err)
			return true
		}
		return true
	})
	_ = fmt.Sprintf // ensure fmt is referenced even if build strips debug logs
	return &notificationManager{cancel}
}

type notificationManager struct {
	cancel func()
}

func (nm *notificationManager) Stop() {
	nm.cancel()
}

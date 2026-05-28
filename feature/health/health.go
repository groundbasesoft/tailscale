// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

// Package health registers the localapi /health endpoint that exposes the
// daemon's structured health.State (with WarnableCode-keyed warnings) to
// localapi clients.
//
// It must be imported by binaries that want the endpoint; the
// feature/condregister/maybe_health.go file does this conditionally so it
// can be linked out with the ts_omit_health build tag.
package health

import (
	"encoding/json"
	"net/http"

	"tailscale.com/ipn/localapi"
)

func init() {
	localapi.Register("health", serveHealth)
}

// serveHealth returns a JSON snapshot of the daemon's current health state,
// with each unhealthy Warnable keyed by its WarnableCode. This is the
// structured equivalent of the rendered strings in ipnstate.Status.Health.
func serveHealth(h *localapi.Handler, w http.ResponseWriter, r *http.Request) {
	if !h.PermitRead {
		http.Error(w, "health access denied", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.LocalBackend().HealthTracker().CurrentState())
}

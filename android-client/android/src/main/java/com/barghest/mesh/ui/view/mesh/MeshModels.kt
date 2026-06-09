// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package com.barghest.mesh.ui.view.mesh

/** Identity of the analyst on the other end of a session; live name/IP come from the netmap peer. */
data class Analyst(
    val name: String,
    val org: String,
    val ip: String,
    val verified: Boolean,
)

/** Empty-state defaults shown until a real peer identity is available. */
object MeshDefaults {
    // No display copy here — a blank name renders the localized placeholder at the call site.
    val analyst = Analyst(
        name = "",
        org = "",
        ip = "—",
        verified = false,
    )
}

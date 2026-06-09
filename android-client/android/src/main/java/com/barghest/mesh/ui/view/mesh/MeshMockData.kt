// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package com.barghest.mesh.ui.view.mesh

/** Sample analyst + guided copy; live name/IP come from the netmap peer in production. */

data class Analyst(
    val name: String,
    val org: String,
    val ip: String,
    val verified: Boolean,
)

object MeshSample {
    val analyst = Analyst(
        name = "Your analyst",
        org = "",
        ip = "—",
        verified = false,
    )
}

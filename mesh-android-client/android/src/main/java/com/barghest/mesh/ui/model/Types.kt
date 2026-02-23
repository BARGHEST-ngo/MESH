// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package com.barghest.mesh.ui.model

import kotlinx.serialization.Serializable

typealias Addr = String

typealias Prefix = String

typealias NodeID = Long

typealias KeyNodePublic = String

typealias MachineKey = String

typealias UserID = Long

typealias Time = String

typealias StableNodeID = String

typealias BugReportID = String

val GoZeroTimeString = "0001-01-01T00:00:00Z"

// Represents and empty message with a single 'property' field.
class Empty {
  @Serializable data class Message(val property: String = "")
}

// Parsable errors returned by localApiService
class Errors {
  @Serializable data class GenericError(val error: String)
}

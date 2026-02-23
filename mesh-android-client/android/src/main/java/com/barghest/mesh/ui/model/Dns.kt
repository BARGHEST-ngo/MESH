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

class Dns {
  @Serializable data class HostEntry(val addr: Addr?, val hosts: List<String>?)

  @Serializable
  data class OSConfig(
      val hosts: List<HostEntry>? = null,
      val nameservers: List<Addr>? = null,
      val searchDomains: List<String>? = null,
      val matchDomains: List<String>? = null,
  ) {
    val isEmpty: Boolean
      get() =
          (hosts.isNullOrEmpty()) &&
              (nameservers.isNullOrEmpty()) &&
              (searchDomains.isNullOrEmpty()) &&
              (matchDomains.isNullOrEmpty())
  }
}

class DnsType {
  @Serializable
  data class Resolver(var Addr: String? = null, var BootstrapResolution: List<Addr>? = null)
}

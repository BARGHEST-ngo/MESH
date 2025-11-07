// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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

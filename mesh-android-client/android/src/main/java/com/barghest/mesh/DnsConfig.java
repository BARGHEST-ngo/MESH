// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh;

// Tailscale DNS Config retrieval
//
// Tailscale's DNS support can either override the local DNS servers with a set of servers
// configured in the admin panel, or supplement the local DNS servers with additional
// servers for specific domains like example.com.beta.barghest.net. In the non-override mode,
// we need to retrieve the current set of DNS servers from the platform. These will typically
// be the DNS servers received from DHCP.
//
// Importantly, after the Tailscale VPN comes up it will set a DNS server of 100.100.100.100
// but we still want to retrieve the underlying DNS servers received from DHCP. If we roam
// from Wi-Fi to LTE, we want the DNS servers received from LTE.

public class DnsConfig {
    private String dnsConfigs;

    // getDnsConfigAsString returns the current DNS configuration as a multiline string:
    // line[0] DNS server addresses separated by spaces
    // line[1] search domains separated by spaces
    //
    // For example:
    // 8.8.8.8 8.8.4.4
    // example.com
    //
    // an empty string means the current DNS configuration could not be retrieved.
    String getDnsConfigAsString() {
        String dnsConfig = getDnsConfigs();
        if (dnsConfig != null) {
            return getDnsConfigs().trim();
        }
        return "";
    }

    private String getDnsConfigs() {
        synchronized (this) {
            return this.dnsConfigs;
        }
    }

    boolean updateDNSFromNetwork(String dnsConfigs) {
        synchronized (this) {
            if (!dnsConfigs.equals(this.dnsConfigs)) {
                this.dnsConfigs = dnsConfigs;
                return true;
            } else {
                return false;
            }
        }
    }
}

// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package dns

import (
	"fmt"
	"net/netip"
	"reflect"
	"testing"

	"tailscale.com/tstest"
	"tailscale.com/util/dnsname"
)

func TestOSConfigPrintable(t *testing.T) {
	ocfg := OSConfig{
		Hosts: []*HostEntry{
			{
				Addr:  netip.AddrFrom4([4]byte{100, 1, 2, 3}),
				Hosts: []string{"server", "client"},
			},
			{
				Addr:  netip.AddrFrom4([4]byte{100, 1, 2, 4}),
				Hosts: []string{"otherhost"},
			},
		},
		Nameservers: []netip.Addr{
			netip.AddrFrom4([4]byte{8, 8, 8, 8}),
		},
		SearchDomains: []dnsname.FQDN{
			dnsname.FQDN("foo.beta.tailscale.net."),
			dnsname.FQDN("bar.beta.tailscale.net."),
		},
		MatchDomains: []dnsname.FQDN{
			dnsname.FQDN("ts.com."),
		},
	}
	s := fmt.Sprintf("%+v", ocfg)

	const expected = `{Nameservers:[8.8.8.8] SearchDomains:[foo.beta.tailscale.net. bar.beta.tailscale.net.] MatchDomains:[ts.com.] Hosts:[&{Addr:100.1.2.3 Hosts:[server client]} &{Addr:100.1.2.4 Hosts:[otherhost]}]}`
	if s != expected {
		t.Errorf("format mismatch:\n   got: %s\n  want: %s", s, expected)
	}
}

func TestIsZero(t *testing.T) {
	tstest.CheckIsZero[OSConfig](t, map[reflect.Type]any{
		reflect.TypeFor[dnsname.FQDN](): dnsname.FQDN("foo.bar."),
		reflect.TypeFor[*HostEntry](): &HostEntry{
			Addr:  netip.AddrFrom4([4]byte{100, 1, 2, 3}),
			Hosts: []string{"foo", "bar"},
		},
	})
}

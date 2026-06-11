package main

import "testing"

func TestSplitFQDN(t *testing.T) {
	cases := []struct{ in, wantName, wantDomain string }{
		{"_acme-challenge.grafana.kralovi.net.", "_acme-challenge.grafana", "kralovi.net"},
		{"_acme-challenge.kralovi.net.", "_acme-challenge", "kralovi.net"},
		{"_acme-challenge.grafana.kralovi.net", "_acme-challenge.grafana", "kralovi.net"}, // no trailing dot
	}
	for _, c := range cases {
		name, domain := splitFQDN(c.in)
		if name != c.wantName || domain != c.wantDomain {
			t.Errorf("splitFQDN(%q) = (%q, %q); want (%q, %q)", c.in, name, domain, c.wantName, c.wantDomain)
		}
	}
}

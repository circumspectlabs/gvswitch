package cmd

import (
	"testing"
)

func TestIPAddressConvertions(t *testing.T) {
	t.Parallel()
	cases := [][]string{
		{"192.168.127.1/24", "192.168.127.1", "192.168.127.254"},
		{"10.10.0.0/16", "10.10.0.1", "10.10.255.254"},
		{"172.16.16.16/12", "172.16.0.1", "172.31.255.254"},
		{"fc00::fff/64", "fc00::1", "fc00::ffff:ffff:ffff:fffe"},
	}
	for _, v := range cases {
		fuaddr, err := getFirstUsableIPFromSubnet(v[0])
		if err != nil {
			t.Errorf("getFirstUsableIPFromSubnet returns error for \"%s\" -> \"%s\": %s", v[0], fuaddr, err.Error())
		}
		luaddr, err := getLastUsableIPFromSubnet(v[0])
		if err != nil {
			t.Errorf("getLastUsableIPFromSubnet returns error for \"%s\" -> \"%s\": %s", v[0], luaddr, err.Error())
		}
		if fuaddr != v[1] {
			t.Errorf("getFirstUsableIPFromSubnet returns wrong result: expects \"%s\", got \"%s\"", v[1], fuaddr)
		}
		if luaddr != v[2] {
			t.Errorf("getLastUsableIPFromSubnet returns wrong result: expects \"%s\", got \"%s\"", v[2], luaddr)
		}
	}
}

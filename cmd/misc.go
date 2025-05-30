package cmd

import (
	"errors"
	"net"
	"net/netip"
)

func getFirstUsableIPFromSubnet(sa string) (string, error) {
	pa, err := netip.ParsePrefix(sa)
	if err != nil {
		return "", err
	}

	// The network must have at least 5 IP addresses: network, broadcast, gateway, guest, and preferably host
	// v4/30 has only 2 devices, thus prefer at least v4/29 CIDR. This code works also for IPv6, just in case
	if (pa.Bits() + 3) > pa.Addr().BitLen() {
		return "", errors.New("too small network")
	}

	return pa.Masked().Addr().Next().String(), nil
}

func getLastUsableIPFromSubnet(sa string) (string, error) {
	pa, err := netip.ParsePrefix(sa)
	if err != nil {
		return "", err
	}

	// The network must have at least 5 IP addresses: network, broadcast, gateway, guest, and preferably host
	// v4/30 has only 2 devices, thus prefer at least v4/29 CIDR. This code works also for IPv6, just in case
	if (pa.Bits() + 3) > pa.Addr().BitLen() {
		return "", errors.New("too small network")
	}

	var b []byte = pa.Masked().Addr().AsSlice()
	for i, v := range net.CIDRMask(pa.Bits(), pa.Addr().BitLen()) {
		b[i] += ^v
	}
	b[len(b)-1] -= 1

	addr, ok := netip.AddrFromSlice(b)
	if !ok {
		return "", errors.New("bad ip address")
	}

	return addr.String(), nil
}

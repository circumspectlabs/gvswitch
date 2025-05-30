package stack

import (
	"fmt"

	"github.com/containers/gvisor-tap-vsock/pkg/types"
	log "github.com/sirupsen/logrus"
)

func (c *Config) AdapterToOriginalStruct() types.Configuration {
	DNS := make([]types.Zone, 0)
	for _, v := range c.DNS {
		entry := types.Zone{
			Name:    v.Name,
			Records: make([]types.Record, 0),
		}
		for _, vv := range v.Records {
			entry.Records = append(entry.Records, types.Record{
				Name: vv.Name,
				IP:   vv.IP,
			})
		}
		DNS = append(DNS, entry)
	}

	forwards := make(map[string]string)
	for _, v := range c.Forwards {
		if (v.From.Proto != "tcp") || (v.To.Proto != "tcp") {
			log.Warningf("Only TCP to TCP forwarding is supported right now. Skipping \"%s:%d->%s:%d\" rule", v.From.IP, v.From.Port, v.To.IP, v.To.Port)
			continue
		}
		forwards[fmt.Sprintf("%s:%d", v.From.IP.String(), v.From.Port)] = fmt.Sprintf("%s:%d", v.To.IP.String(), v.To.Port)
	}

	gwvips := make([]string, 0)
	gwvips = append(gwvips, c.Hypervisor.IP.String())

	dhcpStaticLeases := make(map[string]string)
	dhcpStaticLeases[c.Gateway.IP.String()] = c.Gateway.MAC.String()
	dhcpStaticLeases[c.Hypervisor.IP.String()] = c.Hypervisor.MAC.String()
	for _, v := range c.DHCPStaticLeases {
		dhcpStaticLeases[v.IP.String()] = v.MAC.String()
	}

	return types.Configuration{
		Debug:                  false,
		CaptureFile:            "",
		MTU:                    c.MTU,
		Subnet:                 c.Subnet.String(),
		GatewayIP:              c.Gateway.IP.String(),
		GatewayMacAddress:      c.Gateway.MAC.String(),
		DNS:                    DNS,
		DNSSearchDomains:       c.DNSSearchDomains,
		Forwards:               forwards,
		NAT:                    make(map[string]string), // empty
		GatewayVirtualIPs:      gwvips,
		DHCPStaticLeases:       dhcpStaticLeases,
		VpnKitUUIDMacAddresses: make(map[string]string),
		Protocol:               types.QemuProtocol,
	}
}

package stack

import (
	"net"
)

type Config struct {
	// MTU
	MTU int `json:"mtu,omitempty" yaml:"mtu,omitempty" mapstructure:"mtu"`

	// Subnet of the network
	Subnet IPNet `json:"subnet,omitempty" yaml:"subnet,omitempty" mapstructure:"subnet"`

	// Gateway for devices in the network
	Gateway IPMac `json:"gateway" yaml:"gateway" mapstructure:"gateway"`

	// Hypervisor's IP and MAC, for forwarding
	Hypervisor IPMac `json:"hypervisor" yaml:"hypervisor" mapstructure:"hypervisor"`

	// Port forwarding between the hypervisor and the virtual network devices
	Forwards []Forward `json:"forwards,omitempty" yaml:"forwards,omitempty" mapstructure:"forwards"`

	// DHCP static leases
	DHCPStaticLeases []IPMac `json:"dhcpStaticLeases,omitempty" yaml:"dhcpStaticLeases,omitempty" mapstructure:"dhcpStaticLeases"`

	// Built-in DNS records that will be served by the DNS server embedded in the gateway
	DNS []Zone `json:"dns,omitempty" yaml:"dns,omitempty" mapstructure:"dns"`

	// List of search domains that will be added in all DHCP replies
	DNSSearchDomains []string `json:"dnsSearchDomains,omitempty" yaml:"dnsSearchDomains,omitempty" mapstructure:"dnsSearchDomains"`
}

type Zone struct {
	Name    string   `yaml:"name,omitempty"`
	Records []Record `yaml:"records,omitempty"`
}

type Record struct {
	Name string `yaml:"name,omitempty"`
	IP   net.IP `yaml:"ip,omitempty"`
}

type IPMac struct {
	IP  net.IP `json:"ip,omitempty" yaml:"ip,omitempty" mapstructure:"ip"`
	MAC MAC    `json:"mac,omitempty" yaml:"mac,omitempty" mapstructure:"mac"`
}

type Forward struct {
	From IPProtoPort `json:"from,omitempty" yaml:"from,omitempty" mapstructure:"from"`
	To   IPProtoPort `json:"to,omitempty" yaml:"to,omitempty" mapstructure:"to"`
}

type IPProtoPort struct {
	Proto string `json:"proto,omitempty" yaml:"proto,omitempty" mapstructure:"proto"`
	IP    net.IP `json:"ip,omitempty" yaml:"ip,omitempty" mapstructure:"ip"`
	Port  int    `json:"port,omitempty" yaml:"port,omitempty" mapstructure:"port"`
}

type MAC net.HardwareAddr
type IPNet net.IPNet

package cmd

import (
	"fmt"
	"os"

	"github.com/circumspectlabs/gvswitch/stack"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Config    string       `json:"-" mapstructure:"config"`
	LogLevel  string       `json:"log-level,omitempty" yaml:"log-level,omitempty" mapstructure:"log-level"`
	LogFormat string       `json:"log-format,omitempty" yaml:"log-format,omitempty" mapstructure:"log-format"`
	Daemon    bool         `json:"daemon,omitempty" yaml:"daemon,omitempty" mapstructure:"daemon"`
	Listen    string       `json:"listen,omitempty" yaml:"listen,omitempty" mapstructure:"listen"`
	Stack     stack.Config `json:"stack" yaml:"stack" mapstructure:"stack"`
	Serve     struct {
		Qemu string `json:"qemu,omitempty" yaml:"qemu,omitempty" mapstructure:"qemu"`
	} `json:"serve,omitempty" yaml:"serve,omitempty" mapstructure:"serve"`
}

var cmdConfig = &cobra.Command{
	Use:   "config",
	Short: "Different config files manipulations",
}

func initConfig() {
	viper.SetDefault("stack.mtu", 1500)
	viper.SetDefault("stack.subnet", "192.168.127.0/24")
	viper.SetDefault("stack.gateway.ip", "<DETECT LATER>")
	viper.SetDefault("stack.gateway.mac", stack.GetRandomMAC())
	viper.SetDefault("stack.hypervisor.ip", "<DETECT LATER>")
	viper.SetDefault("stack.hypervisor.mac", stack.GetRandomMAC())

	cmdConfig.AddCommand(&cobra.Command{
		Use:   "template",
		Short: "Show config template",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(configTemplate + "\n")
		},
	})

	cmdConfig.AddCommand(&cobra.Command{
		Use:     "validate",
		Short:   "Validate provided config file",
		PreRunE: loadConfig,
		Run: func(cmd *cobra.Command, args []string) {
			log.Debugln("config is valid")
		},
	})

	cmdConfig.AddCommand(&cobra.Command{
		Use:     "show",
		Short:   "Show config file (with flag overrides)",
		PreRunE: loadConfig,
		Run: func(cmd *cobra.Command, args []string) {
			viper.WriteConfigTo(os.Stdout)
		},
	})

	rcmd.AddCommand(cmdConfig)
}

func loadConfig(cmd *cobra.Command, args []string) error {
	if err := initLogs(rLogLevel, rLogFormat); err != nil {
		return err
	}

	log.Debugln("initializing config file")

	cnf := rConfig
	if rConfig == "" {
		cnf = viper.GetString("config")
	}
	viper.SetConfigFile(cnf)
	viper.SetConfigType("yaml")

	log.Debugf("identified config file is \"%s\"", cnf)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("configuration file \"%s\" hasn't been found", rConfig)
		} else {
			log.Fatalln(err)
		}
	}

	// Never take config from config file
	viper.Set("config", cnf)

	// Reinit logs if necessary
	if rLogLevel == "" {
		rLogLevel = viper.GetString("log-level")
	}
	if rLogFormat == "" {
		rLogFormat = viper.GetString("log-format")
	}
	if err := initLogs(rLogLevel, rLogFormat); err != nil {
		return err
	}

	log.Debugln("validating config file")

	// Basic types check
	if err := validateConfig(cmd, args); err != nil {
		return err
	}

	// Set delayed defaults
	if viper.GetString("stack.gateway.ip") == "<DETECT LATER>" {
		fuip, err := getFirstUsableIPFromSubnet(viper.GetString("stack.subnet"))
		if err != nil {
			log.Fatalf("Failed to set default Gateway IP address: %v", err)
		}
		viper.Set("stack.gateway.ip", fuip)
	}
	if viper.GetString("stack.hypervisor.ip") == "<DETECT LATER>" {
		luip, err := getLastUsableIPFromSubnet(viper.GetString("stack.subnet"))
		if err != nil {
			log.Fatalf("Failed to set default Hypervisor IP address: %v", err)
		}
		viper.Set("stack.hypervisor.ip", luip)
	}

	log.Debugln("finished initializing config file")

	return nil
}

func validateConfig(cmd *cobra.Command, args []string) error {
	_, err := GetConfig()
	if err != nil {
		return err
	}

	return nil
}

func GetConfig() (*Config, error) {
	var cnf Config

	composed := mapstructure.ComposeDecodeHookFunc(
		stack.IPNetHookFunc(),
		stack.MACHookFunc(),
		stack.IPHookFunc(),
	)

	if err := viper.Unmarshal(&cnf, viper.DecodeHook(composed)); err != nil {
		return nil, err
	}

	return &cnf, nil
}

const configTemplate string = `## gvswitch config
##

## Enables notify call for systemd daemon. Only use for systemd
##
daemon: false

## Log level. Valid options: "fatal", "error", "info", "debug", "trace".
## "trace" shows traffic in the switch (similar to tcpdump for "text"
## format).
##
log-level: info

## Log format. Valid options: "text", "json". Default is "text". It can
## automatically detect TTY and enable color output.
##
log-format: text

## Management socket listener. No management socket by default (means
## empty string). Supports unix and tcp protocols. Examples:
##   unix:///var/run/gvswitch.ctl.sock # absolute path
##   unix://gvswitch.ctl.sock          # relative path
##   tcp://127.0.0.1:8888
##
listen: ""

## Virtual switch configuration
##
stack:
  ## MTU (https://en.wikipedia.org/wiki/Maximum_transmission_unit)
  ## Usually should be 1500, defaults to 1500.
  ##
  mtu: 1500

  ## Subnet IP with CIDR. Default is 192.168.127.0/24.
  ##
  subnet: 192.168.127.0/24

  ## Gateway IP and MAC addresses. Default IP is the first valid
  ## in the subnet. Default MAC address is just random.
  ##
  gateway:
    ip: 192.168.127.1
    mac: 5a:94:ef:e4:0c:dd

  ## The host machine will use this IP and MAC addresses while
  ## forwarding ports (see below). Default IP is the last valid
  ## in the subnet (not broadcast). Default MAC address is just
  ## random.
  ##
  hypervisor:
    # enabled: true
    ip: 192.168.127.254
    mac: 5a:94:ef:e4:0c:ff

  ## Internal DNS configuration. Default is empty.
  ##
  dns:
    - name: molecule.internal.
      records:
        - name: gateway
          ip: 192.168.127.1
        - name: host
          ip: 192.168.127.254

  ## Forwarding rules from local machine into the virtual network.
  ## Useful for exposing SSH as a local port. No forwards by default.
  ##
  forwards:
    - from:
        proto: tcp
        ip: 127.0.0.1
        port: 2222
      to:
        proto: tcp
        ip: 192.168.127.2
        port: 22

  ## You can enable "static" IP address for particular MAC address.
  ## No static leases by default.
  ##
  dhcpStaticLeases:
    - ip: 192.168.127.2
      mac: 5a:94:ef:e4:0c:ee

## Listen for hypervisor connections. Every socket can accept multiple VMs.
## You must specify at least one listen socket. Examples:
##   unix:///var/run/gvswitch.qemu.sock # absolute path
##   unix://gvswitch.qemu.sock          # relative path
##   tcp://127.0.0.1:9999
##
serve:
  qemu: ""
`

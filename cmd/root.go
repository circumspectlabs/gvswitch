package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rConfig string
var rLogLevel string
var rLogFormat string
var rDaemon bool
var rDebug bool
var rcmd = &cobra.Command{
	Use:   "gvswitch -c config.yaml [-l debug] ... ",
	Short: "A completely userspace low-privileged virtual switch for QEMU VMs",
	Long:  "This is a simple utility for running completely userspace low-privileged virtual L2+ switch with additional features like simple routing and DHCP. It is useful for QEMU virtual machines without root permissions (especially on MacOS).",
}

func init() {
	rcmd.PersistentFlags().StringVarP(&rConfig, "config", "c", "", "config file")
	rcmd.PersistentFlags().StringVarP(&rLogLevel, "log-level", "l", "", "log level: \"fatal\", \"error\", \"info\", \"debug\", or \"trace\" (default \"info\")")
	rcmd.PersistentFlags().StringVar(&rLogFormat, "log-format", "", "log format: \"text\" or \"json\" (default \"text\")")
	rcmd.PersistentFlags().BoolVarP(&rDaemon, "daemon", "d", false, "run as systemd service with notify call")
	viper.BindPFlag("config", rcmd.PersistentFlags().Lookup("config"))
	viper.BindEnv("config", "GVSWITCH_CONFIG")
	viper.BindPFlag("log-level", rcmd.PersistentFlags().Lookup("log-level"))
	viper.BindEnv("log-level", "GVSWITCH_LOG_LEVEL")
	viper.SetDefault("log-level", "info")
	viper.BindPFlag("log-format", rcmd.PersistentFlags().Lookup("log-format"))
	viper.BindEnv("log-format", "GVSWITCH_LOG_FORMAT")
	viper.SetDefault("log-format", "text")
	viper.BindPFlag("daemon", rcmd.PersistentFlags().Lookup("daemon"))
	viper.BindEnv("daemon", "GVSWITCH_DAEMON")
	viper.SetDefault("daemon", false)

	initConfig()

	rcmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	})

	rcmd.AddCommand(cmdServe)
}

// func preInitLogsInterface(cmd *cobra.Command, args []string) error {
// 	return initLogs(rLogLevel, rLogFormat)
// }

func initLogs(l string, f string) error {
	text := &log.TextFormatter{
		FullTimestamp: true,
	}
	json := &log.JSONFormatter{}

	if l != "" {
		lvl, err := log.ParseLevel(l)
		if err != nil {
			return err
		}
		log.SetLevel(lvl)
	}

	switch f {
	case "":
		log.SetFormatter(text)
	case "text":
		log.SetFormatter(text)
	case "json":
		log.SetFormatter(json)
	default:
		return fmt.Errorf("Bad log format \"%s\"", f)
	}

	return nil
}

func Entrypoint() {
	if err := rcmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

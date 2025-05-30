package cmd

import (
	"context"
	"sync"

	"github.com/circumspectlabs/gvswitch/control"
	"github.com/circumspectlabs/gvswitch/misc"
	"github.com/circumspectlabs/gvswitch/qemu"

	"github.com/coreos/go-systemd/daemon"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/containers/gvisor-tap-vsock/pkg/virtualnetwork"
)

var cmdServe = &cobra.Command{
	Use:     "serve",
	Short:   "Start and serve virtual switch",
	PreRunE: loadConfig,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugln("starting the virtual switch")

		// Get the config and original stack structure
		cnf, err := GetConfig()
		if err != nil {
			log.Fatalf("failed to get configuration: %v", err)
		}
		stack := cnf.Stack.AdapterToOriginalStruct()

		// Initialize virtual network
		vn, err := virtualnetwork.New(&stack)
		if err != nil {
			log.Fatalf("failed to initialize virtual network: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		go misc.WaitForSignal(cancel)

		wg := sync.WaitGroup{}
		wg.Add(1)

		// Control socket
		wg.Add(1)
		go func() {
			log.Debugln("listening control socket")
			ctx, _ := context.WithCancel(ctx)

			if err := control.Serve(ctx, cnf.Listen); err != nil {
				log.Errorf("control socket listener has finished with error: %v", err)
			}

			cancel()
			wg.Done()
		}()

		// Listen qemu socket
		wg.Add(1)
		go func() {
			log.Debugln("listening qemu socket")
			ctx, _ := context.WithCancel(ctx)

			if err := qemu.Serve(ctx, cnf.Serve.Qemu, vn); err != nil {
				log.Errorf("control socket listener has finished with error: %v", err)
			}

			cancel()
			wg.Done()
		}()

		if cnf.Daemon {
			daemon.SdNotify(false, daemon.SdNotifyReady)
		}

		<-ctx.Done()

		if cnf.Daemon {
			daemon.SdNotify(false, daemon.SdNotifyStopping)
		}

		wg.Done()
		wg.Wait()
	},
}

package misc

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/containers/winquit/pkg/winquit"
	log "github.com/sirupsen/logrus"
)

func WaitForSignal(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, os.Interrupt)
	winquit.SimulateSigTermOnQuit(signals)
	log.Debugln("waiting for signal")
	log.Infof("received signal: %s", <-signals)
	cancel()
	close(signals)
}

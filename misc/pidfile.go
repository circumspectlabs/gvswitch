package misc

import (
	"context"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func WithPidFile(ctx context.Context, pidfile string) error {
	log.Debugf("creating pidfile \"%s\"", pidfile)

	fd, err := os.Create(pidfile)
	if err != nil {
		return fmt.Errorf("failed to create pidfile: %v", err)
	}
	if _, err := fd.WriteString(strconv.Itoa(os.Getpid())); err != nil {
		return fmt.Errorf("failed to write into pidfile: %v", err)
	}

	finisher := make(chan interface{})
	go func() {
		<-ctx.Done()
		os.Remove(pidfile)
		log.Debugf("released pidfile \"%s\"", pidfile)
		finisher <- true
	}()

	<-finisher

	return nil
}

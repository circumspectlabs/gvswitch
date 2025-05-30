package qemu

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"

	"github.com/circumspectlabs/gvswitch/misc"
	log "github.com/sirupsen/logrus"

	"github.com/containers/gvisor-tap-vsock/pkg/virtualnetwork"
)

func Serve(ctx context.Context, addr string, vn *virtualnetwork.VirtualNetwork) error {
	parsed, err := misc.ParseListenAddressAndVerify(addr, "qemu")
	if err != nil {
		return err
	}

	var listener net.Listener

	switch parsed.Scheme {
	case "unix":
		if parsed.Host != "" {
			parsed.Path = fmt.Sprintf("%s%s%s", parsed.Host, string(os.PathSeparator), parsed.Path)
			parsed.Host = ""
		}
		parsed.Path = strings.TrimRight(parsed.Path, string(os.PathSeparator))
		if runtime.GOOS == "windows" {
			parsed.Path = strings.TrimPrefix(parsed.Path, string(os.PathSeparator))
		}
		listener, err = net.Listen(parsed.Scheme, parsed.Path)
	case "tcp", "tcp4", "tcp6":
		listener, err = net.Listen(parsed.Scheme, parsed.Host)
	default:
		return errors.New("qemu listener: unexpected scheme")
	}

	if err != nil {
		return fmt.Errorf("cannot open qemu listener on %s: %v", addr, err)
	}

	finisher := make(chan interface{})
	go func() {
		<-ctx.Done()
		if err := listener.Close(); err != nil {
			log.Errorf("error while closing qemu listener on %s: %v", addr, err)
		}
		log.Debugf("closed qemu listener on \"%s\"", addr)
		if parsed.Scheme == "unix" {
			os.Remove(parsed.Path)
		}
		finisher <- true
	}()

	log.Debugf("listening for qemu connections on \"%s\"", addr)

acceptloop:
	for {
		select {
		case <-ctx.Done():
			break acceptloop
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Errorf("qemu accept error: %v", err)
				continue
			}
			log.Infof("accepted qemu connection from \"%s\"", conn.RemoteAddr())
			go func() {
				if err := vn.AcceptQemu(ctx, conn); err != nil {
					log.Errorf("qemu accept error: %v", err)
				}
				log.Infof("closed qemu connection from \"%s\"", conn.RemoteAddr())
			}()
		}
	}

	<-finisher

	return nil
}

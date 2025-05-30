package misc

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func ParseListenAddressAndVerify(addr string, source string) (*url.URL, error) {
	parsed, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme == "" {
		parsed.Scheme = "tcp"
	}

	switch source {
	case "control":
		if !((parsed.Scheme == "unix") || strings.HasPrefix(parsed.Scheme, "tcp")) {
			return nil, fmt.Errorf("parse listen address error: \"%s\" doesn't support \"%s\" listening", source, parsed.Scheme)
		}
	case "qemu":
		if !((parsed.Scheme == "unix") || strings.HasPrefix(parsed.Scheme, "tcp")) {
			return nil, fmt.Errorf("parse listen address error: \"%s\" doesn't support \"%s\" listening", source, parsed.Scheme)
		}
	default:
		return nil, errors.New("parse listen address error: undefined source")
	}

	return parsed, err
}

package stack

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

// IPNet
var _ json.Marshaler = IPNet{}
var _ yaml.Marshaler = IPNet{}

func (f IPNet) String() string {
	ones, _ := f.Mask.Size()
	return fmt.Sprintf(`%s/%d`, f.IP.String(), ones)
}

func (f IPNet) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, f.String())), nil
}

func (f IPNet) MarshalYAML() (interface{}, error) {
	return []byte(f.String()), nil
}

func IPNetHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t == reflect.TypeOf(IPNet{}) {
			_, ipnet, err := net.ParseCIDR(data.(string))
			return ipnet, err
		}

		return data, nil
	}
}

// MAC
var _ json.Marshaler = MAC{}
var _ yaml.Marshaler = MAC{}

func (a MAC) String() string {
	const hexDigit = "0123456789abcdef"
	if len(a) == 0 {
		return ""
	}
	buf := make([]byte, 0, len(a)*3-1)
	for i, b := range a {
		if i > 0 {
			buf = append(buf, ':')
		}
		buf = append(buf, hexDigit[b>>4])
		buf = append(buf, hexDigit[b&0xF])
	}
	return string(buf)
}

func (f MAC) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, f.String())), nil
}

func (f MAC) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf(`%s`, f.String()), nil
}

func MACHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t == reflect.TypeOf(MAC{}) {
			mac, err := net.ParseMAC(data.(string))
			return mac, err
		}

		return data, nil
	}
}

// IP
func IPHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t == reflect.TypeOf(net.IP{}) {
			mac := net.ParseIP(data.(string))
			return mac, nil
		}

		return data, nil
	}
}

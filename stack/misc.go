package stack

import (
	"encoding/binary"
	"fmt"
	"math/rand"
)

var unclaimedMACVendorIDs []string = []string{
	"20:67:b1", "28:e2:98", "2C:9e:fd",
	"30:ae:f7", "3c:1a:58", "44:5e:ce",
	"48:fe:eb", "64:55:b2", "74:6a:3b",
	"7c:5a:1d", "84:26:2c", "90:5c:45",
	"cc:ea:1d", "d4:01:2a", "e4:af:a2",
}

func GetRandomMAC() string {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, rand.Uint32())
	return fmt.Sprintf("%s:%02x:%02x:%02x", unclaimedMACVendorIDs[rand.Intn(15)], buf[0], buf[1], buf[2])
}

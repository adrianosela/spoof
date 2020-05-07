package app

import (
	"encoding/binary"
	"net"
)

func getBroadcastIPv4(ip net.IP, mask net.IPMask) net.IP {
	i := binary.BigEndian.Uint32(ip.To4())
	m := binary.BigEndian.Uint32(net.IP(mask).To4())

	bc := make(net.IP, 4)
	binary.BigEndian.PutUint32(bc, i|^m)

	return bc
}

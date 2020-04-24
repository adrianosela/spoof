package payloads

import (
	"fmt"
	"net"
)

// Type represents the payload type
type Type string

// Config carries the payload configuration
type Config struct {
	SrcIP  net.IP
	DstIP  net.IP
	SrcMAC net.HardwareAddr
	DstMAC net.HardwareAddr
}

// Build returns a wire-ready payload of the specified type
func Build(t Type, c Config) ([]byte, error) {
	switch t {
	case TypeICMPEcho:
		return ICMPEchoReq(c.SrcIP, c.DstIP, c.SrcMAC, c.DstMAC)
	case TypeARPReply:
		return ARPReply(c.SrcIP, c.DstIP, c.SrcMAC, c.DstMAC)
	default:
		return nil, fmt.Errorf("payload type %s not supported", t)
	}
}

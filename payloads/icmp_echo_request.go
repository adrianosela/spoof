package payloads

import (
	"net"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// TypeICMPEcho is the payload type for an ICMP Echo Request
var TypeICMPEcho = Type("PAYLOAD_ICMP_ECHO_REQUEST")

// ICMPEchoReq returns a wire-ready Ethernet frame
// encapsulating an IPv4 packet carrying an ICMP Echo Request
func ICMPEchoReq(sIP, dIP net.IP, sMAC, dMAC net.HardwareAddr) ([]byte, error) {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	eth := &layers.Ethernet{
		EthernetType: layers.EthernetTypeIPv4,
		SrcMAC:       sMAC,
		DstMAC:       dMAC,
	}
	ip := &layers.IPv4{
		SrcIP:    sIP,
		DstIP:    dIP,
		Protocol: layers.IPProtocolICMPv4,
		Version:  4,
		TTL:      32,
	}
	icmp := &layers.ICMPv4{
		TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
		Id:       uint16(os.Getpid()) & 0xffff,
		Seq:      0x0001,
	}
	if err := gopacket.SerializeLayers(buf, opts, eth, ip, icmp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

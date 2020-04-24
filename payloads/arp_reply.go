package payloads

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// TypeARPReply is the payload type for an ARP Reply
var TypeARPReply = Type("PAYLOAD_ARP_REPLY")

// ARPReply returns a wire-ready Ethernet frame
// encapsulating an (unsolicited) ARP reply
func ARPReply(sIP, dIP net.IP, sMAC, dMAC net.HardwareAddr) ([]byte, error) {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	eth := &layers.Ethernet{
		EthernetType: layers.EthernetTypeARP,
		SrcMAC:       sMAC,
		DstMAC:       dMAC,
	}
	arp := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		SourceHwAddress:   []byte(sMAC),
		SourceProtAddress: []byte(sIP.To4()),
		DstHwAddress:      []byte(dMAC),
		DstProtAddress:    []byte(dIP.To4()),
	}
	if err := gopacket.SerializeLayers(buf, opts, eth, arp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

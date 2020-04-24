package main

import (
	"net"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// returns a wire-ready Ethernet frame encapsulating an IPv4 ICMP Echo Request
func spoofedICMP(sip, dip net.IP, smac, dmac net.HardwareAddr) ([]byte, error) {
	buf := gopacket.NewSerializeBuffer()

	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	eth := &layers.Ethernet{
		SrcMAC:       smac,
		DstMAC:       dmac,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ip := &layers.IPv4{
		SrcIP:    sip,
		DstIP:    dip,
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

package main

import (
	"encoding/hex"
	"log"
	"net"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func spoofedICMP(spoof, dstIP net.IP) ([]byte, error) {
	ip := &layers.IPv4{
		SrcIP:    spoof, // spoofed source IP
		DstIP:    dstIP,
		Protocol: layers.IPProtocolICMPv4,
		Version:  4,
		TTL:      32,
	}

	icmp := &layers.ICMPv4{
		TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
		Id:       uint16(os.Getpid()) & 0xffff,
		Seq:      0x0001,
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}

	if err := gopacket.SerializeLayers(buf, opts, ip, icmp); err != nil {
		return nil, err
	}

	pl := buf.Bytes()
	pl[0] = byte(8) // gopacket marks ICMP ECHO REQUEST as ICMP type 69 which is incorrect

	log.Printf("[payload] \n%v", hex.Dump(buf.Bytes()))
	return buf.Bytes(), nil
}

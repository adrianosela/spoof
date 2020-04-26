package main

import (
	"net"

	"github.com/adrianosela/spoof/payloads"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
)

type spoofer struct {
	wire    *pcap.Handle
	payload []byte
}

func newSpoofer(sIP, dIP net.IP, iface string) (*spoofer, error) {
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
		SrcMAC: nif.HardwareAddr,                               // doesn't matter
		DstMAC: net.HardwareAddr{255, 255, 255, 255, 255, 255}, // broadcast MAC
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not build a spoofed payload")
	}
	return &spoofer{
		wire:    wire,
		payload: payload,
	}, nil
}

func (s *spoofer) inject() error {
	return s.wire.WritePacketData(s.payload)
}

func (s *spoofer) close() {
	s.wire.Close()
}

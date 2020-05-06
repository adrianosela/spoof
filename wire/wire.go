package wire

import (
	"net"

	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
)

// Wire represents a live link
type Wire struct {
	iface *net.Interface
	pcap  *pcap.Handle
}

// NewWire is the Wire constructor
func NewWire(iface string) (*Wire, error) {
	nif, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, errors.Wrap(err, "could not establish network interface")
	}
	link, err := pcap.OpenLive(iface, int32(nif.MTU), false, pcap.BlockForever)
	if err != nil {
		return nil, errors.Wrap(err, "could not set pcap on network interface")
	}
	return &Wire{
		iface: nif,
		pcap:  link,
	}, nil
}

// Inject a frame onto a live link
func (w *Wire) Inject(frame []byte) error {
	return w.pcap.WritePacketData(frame)
}

// Close a live link
func (w *Wire) Close() {
	w.pcap.Close()
}

// MAC returns the hardware address of the live link
func (w *Wire) MAC() net.HardwareAddr {
	return w.iface.HardwareAddr
}

// IP returns the IPv4 address of the network interface
func (w *Wire) IP() (net.IP, net.IPMask, error) {
	addrs, err := w.iface.Addrs()
	if err != nil {
		return nil, nil, errors.Wrap(err, "no ip address found for network interface")
	}

	var ip net.IP
	var mask net.IPMask

	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
			mask = v.Mask
		case *net.IPAddr:
			ip = v.IP
			mask = ip.DefaultMask()
		}
		if ip == nil {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			continue
		}
		return ip, mask, nil
	}

	return nil, nil, errors.New("no ip address found for network interface")
}

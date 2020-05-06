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

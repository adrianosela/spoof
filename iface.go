package main

import (
	"net"

	"github.com/pkg/errors"
)

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

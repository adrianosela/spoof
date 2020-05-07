package app

import (
	"encoding/binary"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func getBroadcastIPv4(ip net.IP, mask net.IPMask) net.IP {
	i := binary.BigEndian.Uint32(ip.To4())
	m := binary.BigEndian.Uint32(net.IP(mask).To4())

	bc := make(net.IP, 4)
	binary.BigEndian.PutUint32(bc, i|^m)

	return bc
}

// runs the function f periodically
func loop(every time.Duration, f func()) {
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(every)
	defer ticker.Stop()

	f() // run *now*

	for {
		select {
		case <-shutdownSignal:
			return
		case <-ticker.C:
			f()
		}
	}
}

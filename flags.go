package main

import (
	"fmt"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

const (
	mandatoryTag = "[mandatory]"
)

var (
	victimFlag = cli.StringFlag{
		Name:  "victim, v",
		Usage: "IP address of victim for smurf attack",
	}
	proxyFlag = cli.StringFlag{
		Name:  "proxy, p",
		Usage: "IP address of network to help in the attack",
	}
	gatewayFlag = cli.StringFlag{
		Name:  "gateway, g",
		Usage: "MAC address of gateway router for LAN",
	}
	ifaceFlag = cli.StringFlag{
		Name:  "iface, i",
		Usage: "name of network interface to use",
	}
)

// name returns the long name of a flag
// note that the split function returns the original string in index 0
// if it does not contain the given delimiter ","
func name(f cli.Flag) string {
	return strings.Split(f.GetName(), ",")[0]
}

func withDefault(f cli.StringFlag, def string) cli.StringFlag {
	f.Value = def
	return f
}

func asMandatory(f cli.StringFlag) cli.StringFlag {
	f.Usage = fmt.Sprintf("%s %s", mandatoryTag, f.Usage)
	return f
}

func assertSet(ctx *cli.Context, flags ...cli.Flag) error {
	for _, f := range flags {
		if !ctx.IsSet(name(f)) {
			return fmt.Errorf("missing %s argument \"%s\"", mandatoryTag, name(f))
		}
	}
	return nil
}

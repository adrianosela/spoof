package app

import (
	"fmt"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

const (
	mandatoryTag = "[mandatory]"
)

var (
	targetFlag = cli.StringFlag{
		Name:  "target, t",
		Usage: "target IP address",
	}
	ifaceFlag = cli.StringFlag{
		Name:  "iface, i",
		Usage: "network interface to use",
	}
	everyFlag = cli.StringFlag{
		Name:  "every, e",
		Usage: "payload send interval e.g. \"10ns\", \"5us\", \"8ms\", \"1s\"",
	}
	srcIPFlag = cli.StringFlag{
		Name:  "srcIP",
		Usage: "source IP",
	}
	dstIPFlag = cli.StringFlag{
		Name:  "dstIP",
		Usage: "destination IP",
	}
	srcMACFlag = cli.StringFlag{
		Name:  "srcMAC",
		Usage: "source IP",
	}
	dstMACFlag = cli.StringFlag{
		Name:  "dstMAC",
		Usage: "destination MAC address",
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

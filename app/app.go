package app

import (
	"fmt"

	cli "gopkg.in/urfave/cli.v1"
)

// Run runs the spoof cli
func Run(version string, args []string) error {
	app := cli.NewApp()
	app.Version = version
	app.Usage = "a utility for injecting spoofed frames into a network"
	app.Commands = commands
	app.EnableBashCompletion = true
	app.CommandNotFound = func(c *cli.Context, cmd string) {
		c.App.Run([]string{"help"})
		fmt.Printf("\ncommand \"%s\" does not exist\n", cmd)
	}
	return app.Run(args)
}

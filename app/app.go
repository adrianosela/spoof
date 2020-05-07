package app

import (
	"fmt"

	cli "gopkg.in/urfave/cli.v1"
)

const (
	banner = `
             github.com/adrianosela/     $$$$$$\
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
	usage = "a utility for injecting spoofed frames into a network"
)

// Run runs the spoof cli
func Run(version string, args []string) error {
	app := cli.NewApp()

	app.Version = version
	app.Usage = usage
	app.Commands = commands
	app.EnableBashCompletion = true
	app.CommandNotFound = func(c *cli.Context, cmd string) {
		c.App.Run([]string{"help"})
		fmt.Printf("\ncommand \"%s\" does not exist\n", cmd)
	}

	fmt.Println(banner)

	return app.Run(args)
}

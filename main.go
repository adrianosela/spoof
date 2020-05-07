package main

import (
	"fmt"
	"os"

	"github.com/adrianosela/spoof/app"
)

const banner = `
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

var version string // injected at build-time

func main() {
	fmt.Println(banner)
	if err := app.Run(version, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

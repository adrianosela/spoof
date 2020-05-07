package main

import (
	"log"
	"os"

	"github.com/adrianosela/spoof/app"
)

var version string // injected at build-time

func main() {
	if err := app.Run(version, os.Args); err != nil {
		log.Fatal(err)
	}
}

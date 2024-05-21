package main

import (
	"fmt"
	"os"

	"github.com/concrete-eth/ark-rts/engine"
)

func main() {
	app := engine.NewGeth()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

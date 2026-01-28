package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rudrankriyam/App-Store-Connect-CLI/cmd"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	versionInfo := fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date)
	root := cmd.RootCommand(versionInfo)
	defer cmd.CleanupTempPrivateKey()

	if err := root.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			cmd.CleanupTempPrivateKey()
			os.Exit(0)
		}
		cmd.CleanupTempPrivateKey()
		log.Fatalf("error parsing flags: %v\n", err)
	}

	if err := root.Run(context.Background()); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			cmd.CleanupTempPrivateKey()
			os.Exit(1)
		}
		cmd.CleanupTempPrivateKey()
		log.Fatalf("error executing command: %v\n", err)
	}
}

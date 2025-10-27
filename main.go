package main

import (
	"fmt"
	"os"

	"github.com/honeybadger/tf-iamgen/cmd"
)

// Version information injected at build time
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

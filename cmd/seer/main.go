package main

import (
	"os"

	"github.com/flarebyte/baldrick-seer/internal/cli"
)

func main() {
	exitCode := cli.Execute(os.Args[1:], os.Stdout, os.Stderr)
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

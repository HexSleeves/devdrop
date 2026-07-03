package main

import (
	"fmt"
	"os"

	"github.com/liatrio-forge/devdrop-capstone/internal/devspace"
)

var version = "dev"

func main() {
	if err := devspace.NewRootCommand(version).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

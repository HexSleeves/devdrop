package main

import (
	"context"
	"os"

	"charm.land/fang/v2"
	"github.com/liatrio-forge/devdrop-capstone/internal/devspace"
)

var version = "dev"

func main() {
	if err := fang.Execute(context.Background(), devspace.NewRootCommand(version), fang.WithVersion(version)); err != nil {
		os.Exit(1)
	}
}

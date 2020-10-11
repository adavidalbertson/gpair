package main

import (
	"flag"
	"os"

	"github.com/adavidalbertson/gpair/internal/subcommands"
)

func main() {
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case subcommands.AddCmd.Name():
		subcommands.Add()

	case subcommands.RemoveCmd.Name():
		subcommands.Remove()

	case subcommands.UnpairCmd.Name():
		subcommands.Unpair()

	default:
		subcommands.Pair()
	}

	os.Exit(0)
}

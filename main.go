package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adavidalbertson/gpair/internal"
	"github.com/adavidalbertson/gpair/subcommands"
)

func main() {
	var alias, name, email string

	flag.Usage = func() {
		fmt.Println()
	}

	flag.BoolVar(&internal.Verbose, "v", false, "Enable verbose output")

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addCmd.StringVar(&alias, "alias", "", "A short name for the pair")
	addCmd.StringVar(&name, "name", "", "The git username for the pair")
	addCmd.StringVar(&email, "email", "", "The email for the pair")
	addCmd.BoolVar(&internal.Verbose, "v", false, "Enable verbose output")

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case addCmd.Name():
		err := addCmd.Parse(os.Args[2:])
		if err != nil {
			addCmd.Usage()
			os.Exit(0)
		}
		subcommands.Add(*addCmd, alias, name, email)

	default:
		if len(os.Args) >= 2 {
			aliases := os.Args[1:]
			pairs, err := internal.GetPairs(aliases...)
			if err != nil {
				fmt.Println(err.Error())
			}

			for _, pair := range pairs {
				fmt.Println(pair)
			}
			os.Exit(0)
		}

		flag.Usage()

	}

	os.Exit(0)
}

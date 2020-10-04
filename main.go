package main

import (
	"github.com/adavidalbertson/gpair/internal/config"
	"flag"
	"fmt"
	"os"

	"github.com/adavidalbertson/gpair/internal"
	"github.com/adavidalbertson/gpair/internal/subcommands"
)

func main() {
	flag.BoolVar(&internal.Verbose, "v", false, "Enable verbose output")

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addCmd.String("alias", "", "A short name for the pair")
	addCmd.String("name", "", "The git username for the pair")
	addCmd.String("email", "", "The email for the pair")
	addCmd.BoolVar(&internal.Verbose, "v", false, "Enable verbose output")

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	configurator := config.NewConfigurator()

	switch os.Args[1] {
	case addCmd.Name():
		subcommands.Add(*addCmd, configurator)

	default:
		if len(os.Args) >= 2 {
			aliases := os.Args[1:]
			pairs, err := configurator.GetPairs(aliases...)
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

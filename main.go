package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/adavidalbertson/gpair/internal/config"

	"github.com/adavidalbertson/gpair/internal"
	"github.com/adavidalbertson/gpair/internal/subcommands"
)

func init() {
	flag.BoolVar(&internal.Help, "help", false, "Display usage information")
	flag.BoolVar(&internal.Help, "\nh", false, "Display usage information (shorthand)")
	flag.BoolVar(&internal.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&internal.Verbose, "\nv", false, "Enable verbose output (shorthand)")
	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Println("gpair is a utility that makes it easier to share credit for collaboration using git.")
		fmt.Println("It stores the contact info of your frequent collaborators and outputs a 'Co-author' clause for your git commit messages.")
		fmt.Println("Run `gpair ALIAS` to retrieve the 'Co-Author' clause for the collaborator saved under 'ALIAS'.")
		fmt.Println("For multiple collaborators, run `gpair ALIAS_1 [ALIAS_2 ...]`")
		fmt.Println("To add a collaborator, use the 'add' subcommand. For information on using 'add', run `gpair add -h`.")
		oldUsage()
	}
}

func main() {
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	configurator := config.NewConfigurator()

	switch os.Args[1] {
	case subcommands.AddCmd.Name():
		if internal.Help {
			subcommands.AddCmd.Usage()
			os.Exit(0)
		}

		alias, name, email, err := subcommands.ParseAddArgs(os.Args[2:])
		if err != nil {
			panic(err)
		}

		err = subcommands.Add(alias, name, email, configurator)
		if err != nil {
			panic(err)
		}

	case subcommands.RemoveCmd.Name():
		if internal.Help {
			subcommands.RemoveCmd.Usage()
			os.Exit(0)
		}

		aliases, err := subcommands.ParseRemoveArgs(os.Args[2:])
		if err != nil {
			panic(err)
		}

		deleted, err := configurator.DeletePairs(aliases...)
		if err != nil {
			fmt.Println(err.Error())
		}

		if len(deleted) > 0 {
			fmt.Printf("Successfully deleted pairs '%s'\n", strings.Join(deleted, "', '"))
		}

	default:
		flag.Parse()

		if len(os.Args) < 2 || internal.Help {
			flag.Usage()
			os.Exit(0)
		}

		aliases := os.Args[1:]
		pairs, err := configurator.GetPairs(aliases...)
		if err != nil {
			fmt.Println(err.Error())
		}

		for _, pair := range pairs {
			fmt.Println(pair)
		}
	}

	os.Exit(0)
}

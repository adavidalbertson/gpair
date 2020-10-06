package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adavidalbertson/gpair/internal/config"

	"github.com/adavidalbertson/gpair/internal"
	"github.com/adavidalbertson/gpair/internal/subcommands"
)

func main() {
	help := flag.Bool("help", false, "Display usage information")
	flag.BoolVar(help, "\nh", false, "Display usage information (shorthand)")
	flag.BoolVar(&internal.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&internal.Verbose, "\nv", false, "Enable verbose output (shorthand)")
	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Println("gpair is a utility that makes it easier to share credit for collaboration using git.")
		fmt.Println("It stores the contact info of your frequent collaborators and outputs a 'Co-author' clause for your git commit messages.")
		fmt.Println("Run `gpair ALIAS` to retrieve the 'Co-Author' clause for the collaborator saved under 'ALIAS'.")
		fmt.Println("To add a collaborator, use the 'add' subcommand. For information on using 'add', run `gpair add -h`.")
		oldUsage()
	}

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addCmd.String("alias", "", "A short name for the collaborator, used in the `gpair ALIAS` command")
	addCmd.String("name", "", "The git username for the collaborator")
	addCmd.String("email", "", "The email for the collaborator")
	addCmd.Bool("help", false, "Display usage information")
	addCmd.Bool("h", false, "\nDisplay usage information (shorthand)")
	addCmd.BoolVar(&internal.Verbose, "verbose", false, "Enable verbose output")
	addCmd.BoolVar(&internal.Verbose, "v", false, "\nEnable verbose output (shorthand)")
	oldAddUsage := addCmd.Usage
	addCmd.Usage = func() {
		fmt.Println("The 'add' subcommand is used to save your collaborators' git contact info.")
		fmt.Println("It can take positional arguments in the following order: `gpair add [ALIAS] USERNAME EMAIL`")
		fmt.Println("The 'ALIAS' field is optional. If omitted, it will be the same as the username.")
		fmt.Println("You can also set fields explicitly as shown below.")
		oldAddUsage()
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	configurator := config.NewConfigurator()

	switch os.Args[1] {
	case addCmd.Name():
		subcommands.Add(os.Args[2:], *addCmd, configurator)

	default:
		flag.Parse()

		if len(os.Args) < 2 || *help {
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

package subcommands

import (
	"flag"
	"fmt"
	"os"

	"github.com/adavidalbertson/gpair/internal/config"

	"github.com/adavidalbertson/gpair/internal"
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

// Pair is the function executed if no subcommand is passed in
// It prints the git pairing clauses for the collaborators with the given aliases
func Pair() {
	flag.Parse()

	if len(os.Args) < 2 || internal.Help {
		flag.Usage()
		os.Exit(0)
	}

	configurator, err := config.NewConfigurator()
	if err != nil {
		panic(err)
	}

	aliases := os.Args[1:]
	collaborators, err := configurator.GetCollaborators(aliases...)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, collaborator := range collaborators {
		fmt.Println(collaborator)
	}
}

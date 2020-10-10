package subcommands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/adavidalbertson/gpair/internal"
	"github.com/adavidalbertson/gpair/internal/config"
)

// RemoveCmd is the flagset for the 'remove' subcommand
var RemoveCmd flag.FlagSet

func init() {
	RemoveCmd = *flag.NewFlagSet("remove", flag.ExitOnError)
	RemoveCmd.BoolVar(&internal.Help, "help", false, "Display usage information")
	RemoveCmd.BoolVar(&internal.Help, "h", false, "\nDisplay usage information (shorthand)")
	RemoveCmd.BoolVar(&internal.Verbose, "verbose", false, "Enable verbose output")
	RemoveCmd.BoolVar(&internal.Verbose, "v", false, "\nEnable verbose output (shorthand)")
	oldUsage := RemoveCmd.Usage
	RemoveCmd.Usage = func() {
		fmt.Println()
		fmt.Println("The 'remove' subcommand is used to remove a collaborator's git contact info.")
		fmt.Println("It can be run with one or more alias as 'gpair remove ALIAS_1 [ALIAS_2 ...]'.")
		fmt.Println()
		oldUsage()
		RemoveCmd.PrintDefaults()
		fmt.Println()
	}
}

func parseRemoveArgs(args []string) (aliases []string, err error) {
	err = RemoveCmd.Parse(args)

	internal.PrintVerbose("Got aliases: %s", strings.Join(RemoveCmd.Args(), ", "))

	return RemoveCmd.Args(), err
}

// Remove is the function executed by the 'remove' subcommand
// It removes the collaborators with the aliases given in the args
func Remove() {
	aliases, err := parseRemoveArgs(os.Args[2:])
	if err != nil {
		panic(err)
	}

	if internal.Help {
		RemoveCmd.Usage()
		os.Exit(0)
	}

	configurator, err := config.NewConfigurator()
	if err != nil {
		panic(err)
	}

	deleted, err := configurator.DeleteCollaborators(aliases...)
	if err != nil {
		fmt.Println(err.Error())
	}

	if len(deleted) > 0 {
		fmt.Printf("Successfully deleted collaborators '%s'\n", strings.Join(deleted, "', '"))
	}
}

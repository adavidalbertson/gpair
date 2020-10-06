package subcommands

import (
	"flag"
	"fmt"
	"strings"

	"github.com/adavidalbertson/gpair/internal"
)

var RemoveCmd flag.FlagSet

func init() {
	RemoveCmd = *flag.NewFlagSet("remove", flag.ExitOnError)
	RemoveCmd.BoolVar(&internal.Help, "help", false, "Display usage information")
	RemoveCmd.BoolVar(&internal.Help, "h", false, "\nDisplay usage information (shorthand)")
	RemoveCmd.BoolVar(&internal.Verbose, "verbose", false, "Enable verbose output")
	RemoveCmd.BoolVar(&internal.Verbose, "v", false, "\nEnable verbose output (shorthand)")
	oldUsage := RemoveCmd.Usage
	RemoveCmd.Usage = func() {
		fmt.Println("The 'remove' subcommand is used to remove a collaborator's git contact info.")
		fmt.Println("It can be run with one or more alias as `gpair remove ALIAS_1 [ALIAS_2 ...]`.")
		oldUsage()
	}
}

func ParseRemoveArgs(args []string) (aliases []string, err error) {
	err = RemoveCmd.Parse(args)

	internal.PrintVerbose("Got aliases: %s", strings.Join(RemoveCmd.Args(), ", "))

	return RemoveCmd.Args(), err
}

package subcommands

import (
	"os"
	"github.com/adavidalbertson/gpair/internal/git"
	"fmt"
	"github.com/adavidalbertson/gpair/internal"
	"flag"
)

// UnpairCmd is the flagset for the 'unpair' subcommand
var UnpairCmd flag.FlagSet

func init() {
	UnpairCmd = *flag.NewFlagSet("unpair", flag.ExitOnError)
	UnpairCmd.BoolVar(&internal.Help, "help", false, "Display usage information")
	UnpairCmd.BoolVar(&internal.Help, "h", false, "\nDisplay usage information (shorthand)")
	UnpairCmd.BoolVar(&internal.Verbose, "verbose", false, "Enable verbose output")
	UnpairCmd.BoolVar(&internal.Verbose, "v", false, "\nEnable verbose output (shorthand)")
	oldUsage := UnpairCmd.Usage
	UnpairCmd.Usage = func() {
		fmt.Println()
		fmt.Println("The 'unpair' subcommand removes co-author lines from the default commit message.")
		fmt.Println("Use this subcommand when you are done pairing.")
		fmt.Println()
		oldUsage()
		UnpairCmd.PrintDefaults()
		fmt.Println()
	}
}

// Unpair is the function executed by the 'unpair' subcommand
// It unsets the default commit message template
func Unpair() {
	err := UnpairCmd.Parse(os.Args[2:])
	if err != nil {
		panic(err)
	}

	if internal.Help {
		UnpairCmd.Usage()
		os.Exit(0)
	}

	err = git.UnsetTemplate()
	if err != nil {
		panic(err)
	}

	internal.PrintVerbose("Successfully unpaired!")
}
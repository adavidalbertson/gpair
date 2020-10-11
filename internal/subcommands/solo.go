package subcommands

import (
	"os"
	"github.com/adavidalbertson/gpair/internal/git"
	"fmt"
	"github.com/adavidalbertson/gpair/internal"
	"flag"
)

// SoloCmd is the flagset for the 'unpair' subcommand
var SoloCmd flag.FlagSet

func init() {
	SoloCmd = *flag.NewFlagSet("solo", flag.ExitOnError)
	SoloCmd.BoolVar(&internal.Help, "help", false, "Display usage information")
	SoloCmd.BoolVar(&internal.Help, "h", false, "\nDisplay usage information (shorthand)")
	SoloCmd.BoolVar(&internal.Verbose, "verbose", false, "Enable verbose output")
	SoloCmd.BoolVar(&internal.Verbose, "v", false, "\nEnable verbose output (shorthand)")
	oldUsage := SoloCmd.Usage
	SoloCmd.Usage = func() {
		fmt.Println()
		fmt.Println("The 'solo' subcommand removes co-author lines from the default commit message.")
		fmt.Println("Use this subcommand when you are done pairing.")
		fmt.Println()
		oldUsage()
		SoloCmd.PrintDefaults()
		fmt.Println()
	}
}

// Solo is the function executed by the 'solo' subcommand
// It unsets the default commit message template
func Solo() {
	err := SoloCmd.Parse(os.Args[2:])
	if err != nil {
		panic(err)
	}

	if internal.Help {
		SoloCmd.Usage()
		os.Exit(0)
	}

	if !git.IsInstalled() {
		fmt.Println("git needs to be installed for gpair to work.")
		os.Exit(0)
	}

	isCustomTemplate, err := git.IsCustomTemplate()
	if err != nil {
		panic(err)
	}

	if isCustomTemplate {
		fmt.Println("It looks like you are using a custom git commit template already.")
		os.Exit(0)
	}

	_, err = git.GetRepoName()
	if err != nil {
		fmt.Println("gpair must be run inside a git repository")
		os.Exit(0)
	}

	err = git.UnsetTemplate()
	if err != nil {
		panic(err)
	}

	internal.PrintVerbose("Successfully unpaired!")
}
package subcommands

import (
	"text/tabwriter"
	"sort"
	"github.com/adavidalbertson/gpair/internal/config"
	"os"
	"fmt"
	"github.com/adavidalbertson/gpair/internal"
	"flag"
)

// ListCmd is the flagset for the 'list' subcommand
var ListCmd flag.FlagSet

func init() {
	ListCmd = *flag.NewFlagSet("list", flag.ExitOnError)
	ListCmd.BoolVar(&internal.Help, "help", false, "Display usage information")
	ListCmd.BoolVar(&internal.Help, "h", false, "\nDisplay usage information (shorthand)")
	ListCmd.BoolVar(&internal.Verbose, "verbose", false, "Enable verbose output")
	ListCmd.BoolVar(&internal.Verbose, "v", false, "\nEnable verbose output (shorthand)")
	oldUsage := ListCmd.Usage
	ListCmd.Usage = func() {
		fmt.Println()
		fmt.Println("The 'list' subcommand lists available coauthors.")
		fmt.Println()
		oldUsage()
		ListCmd.PrintDefaults()
		fmt.Println()
	}
}


// List is the function executed by the 'list' subcommand
// It lists all configured coauthors
func List() {
	err := ListCmd.Parse(os.Args[2:])
	if err != nil || internal.Help {
		ListCmd.Usage()
		os.Exit(0)
	}

	configurator, err := config.NewConfigurator()
	if err != nil {
		panic(err)
	}

	list, err := configurator.GetCollaborators()
	if err != nil {
		panic(err)
	}

	sort.Slice(list, func(i, j int) bool {
		return config.Less(list[i], list[j])
	})

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0x0)
	for _, collab := range list {
		line := fmt.Sprintf("\t%s\t<%s>", collab.Name, collab.Email)
		if (collab.Alias != collab.Name) {
			line = fmt.Sprintf("%s:%s", collab.Alias, line)
		}
		fmt.Fprintln(tw, line)
	}
	tw.Flush()
}
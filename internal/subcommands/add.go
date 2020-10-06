package subcommands

import (
	"flag"
	"fmt"

	"github.com/adavidalbertson/gpair/internal"
	"github.com/adavidalbertson/gpair/internal/config"
)

var AddCmd flag.FlagSet

func init() {
	AddCmd = *flag.NewFlagSet("add", flag.ExitOnError)
	AddCmd.String("alias", "", "A short name for the collaborator, used in the `gpair ALIAS` command")
	AddCmd.String("name", "", "The git username for the collaborator")
	AddCmd.String("email", "", "The email for the collaborator")
	AddCmd.Bool("help", false, "Display usage information")
	AddCmd.Bool("h", false, "\nDisplay usage information (shorthand)")
	AddCmd.BoolVar(&internal.Verbose, "verbose", false, "Enable verbose output")
	AddCmd.BoolVar(&internal.Verbose, "v", false, "\nEnable verbose output (shorthand)")
	oldAddUsage := AddCmd.Usage
	AddCmd.Usage = func() {
		fmt.Println("The 'add' subcommand is used to save your collaborators' git contact info.")
		fmt.Println("It can take positional arguments in the following order: `gpair add [ALIAS] USERNAME EMAIL`")
		fmt.Println("The 'ALIAS' field is optional. If omitted, it will be the same as the username.")
		fmt.Println("You can also set fields explicitly as shown below.")
		oldAddUsage()
	}
}

func ParseAddArgs(args []string) (alias, name, email string, err error) {
	err = AddCmd.Parse(args)
	if err != nil {
		return
	}

	alias = AddCmd.Lookup("alias").Value.String()
	name = AddCmd.Lookup("name").Value.String()
	email = AddCmd.Lookup("email").Value.String()

	internal.PrintVerbose("-alias='%s' -name='%s' -email='%s'\n", alias, name, email)

	missingArgs := 0
	if name == "" {
		missingArgs++
	}
	if email == "" {
		missingArgs++
	}

	argIndex := 0
	if alias == "" && AddCmd.NArg() > argIndex && AddCmd.NArg() > missingArgs {
		alias = AddCmd.Arg(argIndex)
		argIndex++
		internal.PrintVerbose("alias not set explicitly, we have enough args to use positional argument '%s'\n", alias)
	}
	if name == "" && AddCmd.NArg() > argIndex {
		name = AddCmd.Arg(argIndex)
		argIndex++
		internal.PrintVerbose("name not set explicitly, using positional argument '%s'\n", name)
	}
	if email == "" && AddCmd.NArg() > argIndex {
		email = AddCmd.Arg(argIndex)
		internal.PrintVerbose("email not set explicitly, using positional argument '%s'\n", email)
	}

	if alias == "" {
		alias = name
	}

	return
}

func Add(alias, name, email string, configurator config.Configurator) error {
	if name == "" || email == "" {
		AddCmd.Usage()
		return fmt.Errorf("name and email are required arguments")
	}

	addPair := config.Pair{Name: name, Email: email}
	err := configurator.AddPair(alias, addPair)
	if err != nil {
		return err
	}

	fmt.Printf("Added pair '%s': %s\n", alias, addPair)

	return nil
}

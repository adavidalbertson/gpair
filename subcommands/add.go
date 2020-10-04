package subcommands

import (
	"flag"
	"fmt"

	"github.com/adavidalbertson/gpair/internal"
)

func Add(addCmd flag.FlagSet, alias, name, email string) {
	internal.PrintVerbose("-alias='%s' -name='%s' -email='%s'\n", alias, name, email)

	missingArgs := 0
	if name == "" {
		missingArgs++
	}
	if email == "" {
		missingArgs++
	}

	argIndex := 0
	if alias == "" && addCmd.NArg() > argIndex && addCmd.NArg() > missingArgs {
		alias = addCmd.Arg(argIndex)
		argIndex++
		internal.PrintVerbose("alias not set explicitly, we have enough args to use positional argument '%s'\n", alias)
	}
	if name == "" && addCmd.NArg() > argIndex {
		name = addCmd.Arg(argIndex)
		argIndex++
		internal.PrintVerbose("name not set explicitly, using positional argument '%s'\n", name)
	}
	if email == "" && addCmd.NArg() > argIndex {
		email = addCmd.Arg(argIndex)
		internal.PrintVerbose("email not set explicitly, using positional argument '%s'\n", email)
	}

	if alias == "" {
		alias = name
	}

	addPair := internal.Pair{Name: name, Email: email}
	err := internal.AddPair(alias, addPair)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Added pair '%s': %s\n", alias, addPair)
}

package subcommands

import (
	"flag"
	"fmt"
	"os"

	"github.com/adavidalbertson/gpair/internal/git"
	"github.com/adavidalbertson/gpair/internal/store"

	"github.com/adavidalbertson/gpair/internal/config"

	"github.com/adavidalbertson/gpair/internal"
)

func init() {
	flag.BoolVar(&internal.Help, "help", false, "Display usage information")
	flag.BoolVar(&internal.Help, "h", false, "\nDisplay usage information (shorthand)")
	flag.BoolVar(&internal.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&internal.Verbose, "v", false, "\nEnable verbose output (shorthand)")
	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Println()
		fmt.Println("gpair is a utility that makes it easier to share credit for collaboration on GitHub.")
		fmt.Println("It stores the contact info of your frequent collaborators and adds a 'Co-author' clause to your default commit message.")
		fmt.Println("Run `gpair ALIAS` to retrieve the 'Co-Author' clause for the collaborator saved under 'ALIAS'.")
		fmt.Println("For multiple collaborators, run 'gpair ALIAS_1 [ALIAS_2 ...]'")
		fmt.Println("To add a collaborator, use the 'add' subcommand. For more information, run 'gpair add -h'.")
		fmt.Println("To remove a collaborator, use the 'remove' subcommand. For more information, run 'gpair remove -h'")
		fmt.Println("To stop pairing, run the 'gpair unpair' subcommand. For more information, run 'gpair unpair -h'")
		fmt.Println()
		oldUsage()
		fmt.Println()
	}
}

// Pair is the function executed if no subcommand is passed in
// It prints the git pairing clauses for the collaborators with the given aliases
func Pair() {
	flag.Parse()

	if internal.Help {
		flag.Usage()
		os.Exit(0)
	}

	configurator, err := config.NewConfigurator()
	if err != nil {
		panic(err)
	}

	aliases := flag.Args()
	collaborators, err := configurator.GetCollaborators(aliases...)
	if err != nil {
		fmt.Println(err.Error())
	}

	if len(collaborators) == 0 {
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

	repoName, err := git.GetRepoName()
	if err != nil {
		fmt.Println("gpair must be run inside a git repository")
		os.Exit(0)
	}

	templatePath, err := git.CreateTemplate(repoName, collaborators...)
	if err != nil {
		if efi, ok := err.(*store.ErrFileInaccessible); ok {
			fmt.Printf("Failed to create template file at %s. Make sure appropriate permissions are set.\n", efi.Path)
			os.Exit(0)
		}

		panic(err)
	}

	err = git.SetTemplate(templatePath)
	if err != nil {
		panic(err)
	}

	for _, collaborator := range collaborators {
		internal.PrintVerbose(collaborator.String())
	}
}

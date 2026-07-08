package altversion

import (
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"

	prettyprint "gocmd/testfiles/PrettyPrint"
	altver "gocmd/testfiles/alternateVersions"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "branch",
		Aliases: []string{"br"},
		Short:   "Manage branches",
	}
	cmd.AddCommand(newCreateCommand())
	cmd.AddCommand(newListCommand())
	return cmd
}

func create_alt_ver(cmd *cobra.Command, args []string) error {
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("creating branch with name: %s\n", args[0])

	name, err := altver.CreateAltVer(*repo, args[0])
	if err != nil {
		panic(err)
	}
	prettyprint.PrintMessage(fmt.Sprintf("Created branch %q from current version", name), "-", "Use |ezg switch <name>| to switch to this branch")
	return nil
}

func list_alt_ver(cmd *cobra.Command, args []string) error {
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		panic(err)
	}

	err = altver.ListAltVer(*repo)

	if err != nil {
		panic(err)
	}

	return nil

}

func newCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new branch from current HEAD",
		Args:  cobra.ExactArgs(1),
		RunE:  create_alt_ver,
	}
}

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List all branches",
		Args:    cobra.NoArgs,
		RunE:    list_alt_ver,
	}
}

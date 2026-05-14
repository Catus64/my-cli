package altversion

import (
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"

	altver "gocmd/testfiles/alternateVersions"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "savefile",
		Short: "Manage Savefiles for alternate versions",
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

	fmt.Printf("creating savefile with name: %s\n", args[0])

	err = altver.CreateAltVer(*repo, args[0])
	if err != nil {
		panic(err)
	}

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
		Short: "Create a new alternate version from current HEAD",
		Args:  cobra.ExactArgs(1),
		RunE:  create_alt_ver,
	}
}

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all alternate versions",
		Args:  cobra.NoArgs,
		RunE:  list_alt_ver,
	}
}

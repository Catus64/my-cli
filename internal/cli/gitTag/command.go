package gittag

import (
	// "fmt"
	//githashread "gocmd/testfiles/GitHashRead"

	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func show_ref(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

	_ = gitobj.Create_Ref(*repo, "newtag", "2b984f9b038e1e62ae67dc32bf6665c34d305eed")

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-name [object-hash]",
		Aliases: []string{"set", "name"},
		Short:   "Set a name for a specific object hash",
		Args:    cobra.MaximumNArgs(1),
		RunE:    show_ref,
	}

	return cmd
}

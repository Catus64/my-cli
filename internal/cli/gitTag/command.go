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

	sha, err := gitobj.Ref_Resolve(*repo, "HEAD")

	_ = gitobj.Create_Ref(*repo, args[0], *sha)

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-name [object-hash]",
		Aliases: []string{"set", "name", "tag"},
		Short:   "Set a name for current active version",
		Args:    cobra.MaximumNArgs(1),
		RunE:    show_ref,
	}

	return cmd
}

package showref

import (
	// "fmt"
	//githashread "gocmd/testfiles/GitHashRead"

	gitobj "gocmd/testfiles/GitObject"
	gitshowref "gocmd/testfiles/GitShowRef"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func show_ref(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

	_, err = gitobj.Ref_Resolve(*repo, "HEAD")
	if err != nil {
		panic(err)
	}

	refs := gitshowref.Ref_list(*repo, "", "")
	for i, ref := range refs {
		println(ref, i)
	}

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show-versions [object-hash]",
		Aliases: []string{"show-ver"},
		Short:   "Show tracked versions in the repository",
		Args:    cobra.MaximumNArgs(0),
		RunE:    show_ref,
	}

	return cmd
}

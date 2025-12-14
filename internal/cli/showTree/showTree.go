package showTree

import (
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func LsTree(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}
	obj := githashread.Object_Read(*repo, args[0])
	leafs := gitobj.Tree_Parse(obj.Deserialize())

	//formatting required later

	for _, leaf := range leafs {
		println(leaf.String())
	}
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-tree [tree-hash]",
		Short: "Show contents of Tree object",
		Args:  cobra.MaximumNArgs(1),
		RunE:  LsTree,
	}

	return cmd
}

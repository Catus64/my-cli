package showObject

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func show(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}
	obj := githashread.Object_Read(*repo, args[0])

	// fmt.Println("Object format:", obj.Get_Format())

	if obj.Get_Format() == "tree" || obj.Get_Format() == "leaf" {
		println("This is a tree object. please use the 'ls-tree' command to view its contents.")
		return nil
	}
	fmt.Println(string(obj.Deserialize()))
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-object [object-hash]",
		Short: "Print a friendly greeting",
		Args:  cobra.MaximumNArgs(1),
		RunE:  show,
	}

	return cmd
}

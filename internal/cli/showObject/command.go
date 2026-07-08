package showObject

import (
	githashread "gocmd/testfiles/GitHashRead"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"

	"github.com/spf13/cobra"
)

func show(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}
	obj, err := githashread.Object_Read(*repo, args[0])
	if err != nil {
		panic(err)
	}

	// fmt.Println("Object format:", obj.Get_Format())

	if obj.Get_Format() == "tree" || obj.Get_Format() == "leaf" {
		println("This is a tree object. please use the 'show-tree' command to view its contents.")
		return nil
	}
	// fmt.Println(string(obj.Deserialize()))
	content := string(obj.Deserialize())
	prettyprint.PrintObjectContent(args[0], []byte(content))
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show-object [object-hash]",
		Aliases: []string{"obj"},
		Short:   "Print a friendly greeting",
		Args:    cobra.MaximumNArgs(1),
		RunE:    show,
	}

	return cmd
}

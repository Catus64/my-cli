package add

import (
	Add "gocmd/testfiles/GitAddRemove"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func add(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

	// git add somefile.go
	Add.Add(repo, args, Add.Options{All: false})
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [file]",
		Short: "Add file contents to the Savelist",
		RunE:  add,
	}

	return cmd
}

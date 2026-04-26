package add

import (
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func add(cmd *cobra.Command, args []string) error {
	required := true
	_, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

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

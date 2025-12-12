package showtree

import (
	// "fmt"
	//githashread "gocmd/testfiles/GitHashRead"
	gitLog "gocmd/testfiles/GitLog"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func log(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}
	err = gitLog.Log(*repo)
	if err != nil {
		panic(err)
	}
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log [object-hash]",
		Short: "Print a friendly greeting",
		Args:  cobra.MaximumNArgs(0),
		RunE:  log,
	}

	return cmd
}

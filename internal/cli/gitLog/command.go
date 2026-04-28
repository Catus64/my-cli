package gitLog

import (
	// "fmt"
	//githashread "gocmd/testfiles/GitHashRead"
	gitLog "gocmd/testfiles/GitLog"
	gitpath "gocmd/testfiles/Gitrepostruct"
	unpack "gocmd/testfiles/Unpack"

	"github.com/spf13/cobra"
)

func log(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}
	unpack.UnpackPackfiles(*repo)
	err = gitLog.Log(*repo)
	if err != nil {
		panic(err)
	}
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history [object-hash]",
		Short: "Print a friendly greeting",
		Args:  cobra.MaximumNArgs(0),
		RunE:  log,
	}

	return cmd
}

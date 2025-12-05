package gitinit

import (
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func Git_init(cmd *cobra.Command, args []string) error {
	fmt.Println("making repo")
	rootpath := gitpath.Get_Os_Dir()
	gitrepo := gitpath.Repo_create(rootpath)
	fmt.Println("Initialized empty Git repository in", gitrepo.GitDir)
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a new git repository",
		RunE:  Git_init,
	}
	return cmd
}

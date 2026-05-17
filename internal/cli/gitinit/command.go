package gitinit

import (
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"

	"github.com/spf13/cobra"
)

func Git_init(cmd *cobra.Command, args []string) error {
	// fmt.Println("making repo")
	rootpath := gitpath.Get_Os_Dir()
	gitrepo, _ := gitpath.Repo_create(rootpath)
	prettyprint.PrintMessage(
		"Repository Initialized Successfully",
		gitrepo.GitDir,
		"An empty repository has been created and\nis ready for use.\n If you have files you want to save\nuse ezg git -a then ezg save",
	)
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-repo",
		Short: "Create a new git repository",
		RunE:  Git_init,
	}
	return cmd
}

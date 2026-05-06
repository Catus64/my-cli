package getcurrent

import (
	"fmt"
	gitCurrent "gocmd/testfiles/GitCurrent"
	gitobject "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func get_current(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

	branch, err := gitCurrent.Get_Active_Branch(*repo)
	if err != nil {
		return err
	}
	fmt.Println("Current active branch:", branch)

	index, err := gitobject.Index_Read2(*repo)
	if err != nil {
		panic(err)
	}

	err = gitCurrent.StatusHeadIndex(*repo, *index)

	_, untracked, err := gitCurrent.StatusIndexWorktree(*repo, *index)
	if err != nil {
		return err
	}

	fmt.Println("\nUntracked files:")
	for _, f := range untracked {
		fmt.Println(" ", f)
	}

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Get the current active branch name and modified files in the working directory",
		RunE:  get_current,
	}

	return cmd
}

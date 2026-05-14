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

	// staged changes
	headResult, err := gitCurrent.StatusHeadIndex(*repo, *index)
	if err != nil {
		return err
	}
	fmt.Println("\nChanges to be committed:")
	for _, f := range headResult.Added {
		fmt.Println("  added:   ", f)
	}
	for _, f := range headResult.Modified {
		fmt.Println("  modified:", f)
	}
	for _, f := range headResult.Deleted {
		fmt.Println("  deleted: ", f)
	}

	// unstaged + untracked
	worktreeResult, err := gitCurrent.StatusIndexWorktree(*repo, *index)
	if err != nil {
		return err
	}
	fmt.Println("\nChanges not staged for commit:")
	for _, f := range worktreeResult.Modified {
		fmt.Println("  modified:", f)
	}
	for _, f := range worktreeResult.Deleted {
		fmt.Println("  deleted: ", f)
	}

	fmt.Println("\nUntracked files:")
	for _, f := range worktreeResult.Untracked {
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

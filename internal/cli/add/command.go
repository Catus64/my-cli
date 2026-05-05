package add

import (
	"fmt"
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

	all, _ := cmd.Flags().GetBool("all")

	// --all and specific files are mutually exclusive
	if all && len(args) > 0 {
		return fmt.Errorf("cannot use --all with specific files")
	}
	if !all && len(args) == 0 {
		return fmt.Errorf("specify files to add or use --all")
	}

	Add.Add(repo, args, Add.Options{All: all})
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [file]",
		Short: "Add file contents to the Savelist",
		Long: `Add file contents to the index (staging area).

Examples:
  ezgit add main.go                  add a single file
  ezgit add main.go internal/cli/    add multiple files
  ezgit add --all                    add all tracked/untracked files`,
		RunE: add,
	}

	cmd.Flags().BoolP("all", "a", false, "Add all files in the working directory")

	return cmd
}

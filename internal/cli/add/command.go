package add

import (
	"fmt"
	Add "gocmd/testfiles/GitAddRemove"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"

	"github.com/spf13/cobra"
)

func add(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

	all, _ := cmd.Flags().GetBool("all")
	showIgnored, _ := cmd.Flags().GetBool("show-ignored")

	if all && len(args) > 0 {
		return fmt.Errorf("cannot use --all with specific files")
	}
	if !all && len(args) == 0 {
		return fmt.Errorf("specify files to add or use --all")
	}

	//copy git add . style behaviour
	if len(args) == 1 && args[0] == "." {
		all = true
		args = nil
	}

	result, err := Add.Add(repo, args, Add.Options{All: all})
	if err != nil {
		return err
	}

	// for _, file := range result.Added {
	// 	fmt.Println(file.Path)
	// }

	// build table rows
	var rows []prettyprint.TableRow
	for _, f := range result.Added {
		rows = append(rows, prettyprint.TableRow{
			File:   f.Path,
			SHA:    f.SHA,
			Status: f.Status,
		})
	}

	prettyprint.PrintAddTable(rows, result.Ignored, all)

	if showIgnored && len(result.Ignored) > 0 {
		prettyprint.PrintIgnoredFiles(result.Ignored)
	}

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [file]",
		Short: "Add file contents to the Index / Entry List",
		Long: `Add file contents to the Index / Entry List (other common names include: Staging Area)
		This command adds files into the index which is a list of files that will be saved in the next version.

Examples:
  ezgit add main.go                  add a single file
  ezgit add main.go internal/cli/    add multiple files
  ezgit add --all                    add all tracked/untracked files`,
		RunE: add,
	}

	cmd.Flags().BoolP("all", "a", false, "Add all files in the working directory")

	return cmd
}

package combine

import (
	"fmt"
	gitsave "gocmd/testfiles/GitSave"
	save "gocmd/testfiles/GitSave"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	prettyprint "gocmd/testfiles/PrettyPrint"
	"strings"

	"github.com/spf13/cobra"
)

func merge(cmd *cobra.Command, args []string) error {
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		panic(err)
	}

	info, err := save.Combine(*repo, args[0])
	if err != nil {
		return fmt.Errorf("failed to merge branch: %w", err)
	}
	if info == nil {
		// already up to date, or user cancelled — nothing to name/print
		return nil
	}

	var versionNum int
	var versionName string
	entry, err := gitsave.WriteVersionRef(*repo, info.Branch, info.CommitSHA)
	if err != nil {
		fmt.Println("warning:", err)
		logger.L().Warn("failed to write version ref", "error", err)
	} else {
		versionNum = entry.Number
		versionName = entry.Name
	}

	prettyprint.PrintMergeResult(prettyprint.SaveResult{
		Branch:      info.Branch,
		CommitSHA:   info.CommitSHA,
		TreeSHA:     info.TreeSHA,
		ParentSHA:   strings.Join(info.Parents, " + "),
		Author:      info.Author,
		Timestamp:   info.Timestamp,
		Message:     info.Message,
		VersionNum:  versionNum,
		VersionName: versionName,
	})

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge [branch-name]",
		Short: "Merge changes from another branch into the current branch",
		Args:  cobra.MaximumNArgs(1),
		RunE:  merge,
	}

	return cmd
}

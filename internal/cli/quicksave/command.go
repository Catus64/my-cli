package quicksave

import (
	"fmt"
	setconfig "gocmd/internal/cli/set-config"
	Add "gocmd/testfiles/GitAddRemove"
	gitobj "gocmd/testfiles/GitObject"
	gitsave "gocmd/testfiles/GitSave"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	prettyprint "gocmd/testfiles/PrettyPrint"
	"time"

	"github.com/spf13/cobra"
)

func runQuicksave(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

	// simulate ezg add -a
	_, err = Add.Add(repo, nil, Add.Options{All: true})
	if err != nil {
		return fmt.Errorf("quicksave: failed to add files: %w", err)
	}

	// read index after add
	index, err := gitobj.Index_Read2(*repo)
	if err != nil {
		panic(err)
	}

	// check if index is clean
	is_clean, err := gitsave.CheckCommitReady(*repo, *index)
	if err != nil {
		panic(err)
	}
	if !is_clean {
		return nil
	}

	// Build tree from index(same as save)
	treeSHA, err := gitobj.TreeFromIndex(*repo, *index)
	if err != nil {
		panic(err)
	}
	logger.L().Debug("Quicksave tree built", "treeSHA", treeSHA)

	// Get parent from head (same as save)
	var parents []string
	var parentSHAStr string
	parentSHA, err := gitobj.Ref_Resolve(*repo, "HEAD")
	if err == nil && parentSHA != nil {
		parents = []string{*parentSHA}
		parentSHAStr = *parentSHA
	} else {
		fmt.Println("No parent — this is the first commit")
	}

	// get user config
	user_config, err := setconfig.GetOrPromptConfig()
	if err != nil {
		panic(err)
	}

	author := user_config.Format()
	timestamp := time.Now()
	message := "This Version is Quicksaved"

	// Create Commit
	commitSHA, err := gitsave.Version_Create(*repo, treeSHA, parents, author, timestamp, message)
	if err != nil {
		panic(err)
	}
	logger.L().Info("Quicksave commit created", "sha", commitSHA)

	// Update branch ref
	branchName, err := gitsave.Update_Branch_Ref(*repo, commitSHA)
	if err != nil {
		panic(err)
	}

	// Refresh index again
	err = gitsave.RefreshIndex(*repo, index)
	if err != nil {
		panic(err)
	}

	// ezgit custom version ref naming
	var versionNum int
	var versionName string
	entry, err := gitsave.WriteVersionRef(*repo, branchName, commitSHA)
	if err != nil {
		logger.L().Warn("failed to write version ref", "error", err)
	} else {
		versionNum = entry.Number
		versionName = entry.Name
	}

	//Pretty print version ref
	prettyprint.PrintSaveResult(prettyprint.SaveResult{
		Branch:      branchName,
		CommitSHA:   commitSHA,
		TreeSHA:     treeSHA,
		ParentSHA:   parentSHAStr,
		Author:      author,
		Timestamp:   timestamp,
		Message:     message,
		VersionNum:  versionNum,
		VersionName: versionName,
	})

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "quicksave",
		Short:   "Quickly add all files and save with a generic message",
		Aliases: []string{"qsave", "qs"},
		Args:    cobra.NoArgs,
		RunE:    runQuicksave,
	}
	return cmd
}

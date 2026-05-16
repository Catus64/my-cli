package gitsave

import (
	"bufio"
	"fmt"
	setconfig "gocmd/internal/cli/set-config"
	gitobj "gocmd/testfiles/GitObject"
	gitsave "gocmd/testfiles/GitSave"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func promptMessage() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type your commit message and press Enter(Not advisable to leave empty):")
	fmt.Print("Commit message: ")
	message, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(message) == "" {
		message = "No commit message provided"
	}
	return strings.TrimSpace(message), nil
}

func save(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

	index, err := gitobj.Index_Read2(*repo)
	if err != nil {
		panic(err)
	}

	is_clean, err := gitsave.CheckCommitReady(*repo, *index)
	if err != nil {
		panic(err)
	}
	if !is_clean {
		return nil
	}

	treeSHA, err := gitobj.TreeFromIndex(*repo, *index)
	if err != nil {
		panic(err)
	}
	logger.L().Debug("Tree SHA:", treeSHA)

	var parents []string
	parentSHA, err := gitobj.Ref_Resolve(*repo, "HEAD")
	if err == nil && parentSHA != nil {
		logger.L().Debug(*parentSHA)
		parents = []string{*parentSHA}
	} else {
		fmt.Println("No parent this is the first commit")
	}

	//Get commit message from user input
	message, err := promptMessage()
	if err != nil {
		panic(err)
	}

	logger.L().Info("Creating new commit", "message", message, "treeSHA", treeSHA, "parents", parents)

	//Get author from config
	user_config, err := setconfig.GetOrPromptConfig()
	if err != nil {
		panic(err)
	}

	logger.L().Info("Using author", "author", user_config.Format())

	author := user_config.Format()
	timestamp := time.Now()

	commitSHA, err := gitsave.Version_Create(*repo, treeSHA, parents, author, timestamp, message)
	if err != nil {
		panic(err)
	}
	fmt.Println("New commit SHA:", commitSHA)

	logger.L().Info("Updating branch ref to new commit", "branch", "HEAD", "commitSHA", commitSHA)

	branchName, err := gitsave.Update_Branch_Ref(*repo, commitSHA)
	if err != nil {
		panic(err)
	}
	logger.L().Debug("Saved on branch ", branchName, commitSHA)

	logger.L().Info("Save completed successfully", "branch", branchName, "commitSHA", commitSHA)

	err = gitsave.RefreshIndex(*repo, index)
	if err != nil {
		panic(err)
	}

	logger.L().Info("Index refreshed after save")

	entry, err := gitsave.WriteVersionRef(*repo, branchName, commitSHA)
	if err != nil {
		fmt.Println("warning:", err)
	} else {
		fmt.Printf("Version: %s\n", entry.ShortName())
	}

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save",
		Short: "Save a new version to the repository",
		Args:  cobra.MaximumNArgs(0),
		RunE:  save,
	}

	return cmd
}

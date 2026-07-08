package viewcommand

import (
	"bufio"
	"fmt"
	gitcurrent "gocmd/testfiles/GitCurrent"
	githashread "gocmd/testfiles/GitHashRead"
	gitlog "gocmd/testfiles/GitLog"
	gitobj "gocmd/testfiles/GitObject"
	gitsave "gocmd/testfiles/GitSave"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Interactively browse ezgit saved versions",
		Args:  cobra.NoArgs,
		RunE:  runVersions,
	}
}

func runVersions(cmd *cobra.Command, args []string) error {
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		return err
	}

	// Step 1 — pick a branch
	headsDir := gitpath.Repo_Path(*repo, "refs", "heads")
	entries, err := os.ReadDir(headsDir)
	if err != nil {
		return fmt.Errorf("no branches found")
	}

	activeBranch, _ := gitcurrent.Get_Active_Branch(*repo)

	fmt.Println("Select a branch to browse versions:")
	var branches []string
	for i, e := range entries {
		if e.IsDir() {
			continue
		}
		marker := "  "
		if e.Name() == activeBranch {
			marker = "* "
		}
		fmt.Printf("%s%d. %s\n", marker, i+1, e.Name())
		branches = append(branches, e.Name())
	}

	fmt.Print("\nenter number: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	num, err := strconv.Atoi(input)
	if err != nil || num < 1 || num > len(branches) {
		return fmt.Errorf("invalid selection")
	}
	branch := branches[num-1]

	// Step 2 — load versions for branch
	versions, err := gitsave.ListVersionRefs(*repo, branch)
	if err != nil || len(versions) == 0 {
		fmt.Printf("No versions found for branch %s\n", branch)
		return nil
	}

	// Step 3 — build tree items from versions
	var items []prettyprint.TreeItem
	for _, v := range versions {
		items = append(items, prettyprint.TreeItem{
			Name:   v.ShortName(),
			SHA:    v.SHA,
			IsTree: false,
		})
	}

	// Step 4 — on select, show save result box then offer tree
	onSelect := func(sha string) error {
		// read version info
		entry, err := gitsave.ReadVersionRef(*repo, branch, sha)
		if err != nil {
			return err
		}

		// read commit for metadata
		obj, err := githashread.Object_Read(*repo, sha)
		if err != nil {
			return err
		}
		commit, ok := obj.(*gitobj.GitCommit)
		if !ok {
			return fmt.Errorf("not a commit object")
		}
		commit.Deserialize()

		date, author := gitlog.Format_Date_Author(string(commit.KvlmDict.Dict["author"]))
		treeSHA := strings.TrimSpace(string(commit.KvlmDict.Dict["tree"]))

		parents := gitobj.GetKvlmValues(commit.KvlmDict.Dict, "parent")
		parentDisplay := ""
		if len(parents) > 0 {
			parentDisplay = parents[0]
		}

		// parse timestamp from author line
		timestamp, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", date)

		prettyprint.PrintSaveResult(prettyprint.SaveResult{
			Branch:      branch,
			CommitSHA:   sha,
			TreeSHA:     treeSHA,
			ParentSHA:   parentDisplay,
			Author:      author,
			Timestamp:   timestamp,
			Message:     string(commit.KvlmDict.Dict["data"]),
			VersionNum:  entry.Number,
			VersionName: entry.Name,
		})

		fmt.Print("\npress v to view tree, any other key to go back: ")
		reader := bufio.NewReader(os.Stdin)
		key, _ := reader.ReadByte()
		if key == 'v' {
			header := fmt.Sprintf("v%d · %s | %s", entry.Number, entry.Name, branch)
			return viewTree(*repo, treeSHA, header)
		}
		return nil
	}

	header := fmt.Sprintf("Versions · %s  (%d total)", branch, len(versions))
	return prettyprint.RunRefsViewer(
		items,
		func(sha string) ([]prettyprint.TreeItem, error) { return nil, nil },
		onSelect,
		prettyprint.DefaultViewerConfig,
		header,
	)
}

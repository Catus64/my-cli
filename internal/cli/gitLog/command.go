package gitLog

import (
	// "fmt"
	//githashread "gocmd/testfiles/GitHashRead"

	"fmt"
	"strings"

	gitCurrent "gocmd/testfiles/GitCurrent"
	githashread "gocmd/testfiles/GitHashRead"
	gitlog "gocmd/testfiles/GitLog"
	gitobj "gocmd/testfiles/GitObject"
	unpack "gocmd/testfiles/GitPacketExtractor"
	gitsave "gocmd/testfiles/GitSave"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"

	"github.com/spf13/cobra"
)

func log(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}
	unpack.Extract(*repo) //extract for possible loose files

	// resolve HEAD to starting commit
	headSHA, err := gitobj.Ref_Resolve(*repo, "HEAD")
	if err != nil || headSHA == nil {
		fmt.Println("No commits yet.")
		return nil
	}

	branch, _ := gitCurrent.Get_Active_Branch(*repo)

	// build LogCommit from a SHA
	buildCommit := func(sha string) (prettyprint.LogCommit, error) {
		obj, err := githashread.Object_Read(*repo, sha)
		if err != nil {
			return prettyprint.LogCommit{}, err
		}
		commit, ok := obj.(*gitobj.GitCommit)
		if !ok {
			return prettyprint.LogCommit{}, fmt.Errorf("not a commit: %s", sha)
		}
		commit.Deserialize()

		date, author := gitlog.Format_Date_Author(string(commit.KvlmDict.Dict["author"]))
		treeSHA := strings.TrimSpace(string(commit.KvlmDict.Dict["tree"]))
		parents := gitobj.GetKvlmValues(commit.KvlmDict.Dict, "parent")
		message := strings.TrimSpace(string(commit.KvlmDict.Dict["data"]))

		// try version lookup
		lc := prettyprint.LogCommit{
			SHA:      sha,
			ShortSHA: sha[:7],
			Author:   author,
			Date:     date,
			Message:  message,
			TreeSHA:  treeSHA,
			Parents:  parents,
		}

		if branch != "" {
			entry, err := gitsave.ReadVersionRef(*repo, branch, sha)
			if err == nil {
				lc.HasVersion = true
				lc.VersionNum = entry.Number
				lc.VersionName = entry.Name
			}
		}

		return lc, nil
	}

	// fetchTreeItems for tree viewer
	fetchTreeItems := func(sha string) ([]prettyprint.TreeItem, error) {
		obj, err := githashread.Object_Read(*repo, sha)
		if err != nil {
			return nil, err
		}
		leafs := gitobj.Tree_Parse(obj.Deserialize())
		var items []prettyprint.TreeItem
		for _, leaf := range leafs {
			items = append(items, prettyprint.TreeItem{
				Name:   string(leaf.Path),
				SHA:    leaf.Sha,
				IsTree: strings.HasPrefix(string(leaf.Mode), "04"),
			})
		}
		return items, nil
	}

	// fetchBlob for tree viewer
	fetchBlob := func(sha string) ([]byte, error) {
		obj, err := githashread.Object_Read(*repo, sha)
		if err != nil {
			return nil, err
		}
		return obj.Deserialize(), nil
	}

	// build initial commit
	initial, err := buildCommit(*headSHA)
	if err != nil {
		return fmt.Errorf("failed to read HEAD commit: %w", err)
	}

	return prettyprint.RunLogViewer(
		initial,
		buildCommit,
		fetchTreeItems,
		fetchBlob,
		prettyprint.DefaultViewerConfig,
	)
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history [object-hash]",
		Short: "Print a friendly greeting",
		Args:  cobra.MaximumNArgs(0),
		RunE:  log,
	}

	return cmd
}

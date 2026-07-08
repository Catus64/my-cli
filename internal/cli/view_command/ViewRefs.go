package viewcommand

import (
	"fmt"
	gitshowref "gocmd/testfiles/GitShowRef"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"

	"github.com/spf13/cobra"
)

func newRefsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "refs",
		Short: "Interactively browse all repository latest and tagged versions",
		Args:  cobra.NoArgs,
		RunE:  runRefs,
	}
}

func runRefs(cmd *cobra.Command, args []string) error {
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		return err
	}

	// load all refs into TreeItem list
	refs := gitshowref.Ref_list(*repo, "", "")
	var items []prettyprint.TreeItem
	for refName, sha := range refs {
		shortSHA := sha
		if len(shortSHA) > 7 {
			shortSHA = shortSHA[:7]
		}
		items = append(items, prettyprint.TreeItem{
			Name:   refName,
			SHA:    sha,
			IsTree: false, // refs are leaf nodes
		})
	}

	if len(items) == 0 {
		fmt.Println("No refs found.")
		return nil
	}

	// use a custom viewer that calls viewCommit on select
	return runRefsViewer(*repo, items)
}

func runRefsViewer(repo gitpath.GitRepository, items []prettyprint.TreeItem) error {
	onSelect := func(sha string) error {
		return resolveAndView(repo, sha)
	}
	return prettyprint.RunRefsViewer(items, func(sha string) ([]prettyprint.TreeItem, error) {
		return nil, nil
	}, onSelect, prettyprint.DefaultViewerConfig, "") // no header for refs list itself
}

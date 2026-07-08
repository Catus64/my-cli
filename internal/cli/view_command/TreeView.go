package viewcommand

import (
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"
	"strings"

	"github.com/spf13/cobra"
)

func newTreeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tree <sha>",
		Short: "Interactively browse a tree object",
		Args:  cobra.ExactArgs(1),
		RunE:  runTree,
	}
}

func runTree(cmd *cobra.Command, args []string) error {
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		return err
	}
	return viewTree(*repo, args[0], "") // empty header
}

func viewTree(repo gitpath.GitRepository, sha string, header string) error {
	rootItems, err := buildTreeItems(repo, sha)
	if err != nil {
		return err
	}
	fetchChildren := func(sha string) ([]prettyprint.TreeItem, error) {
		return buildTreeItems(repo, sha)
	}
	fetchBlob := func(sha string) ([]byte, error) {
		obj, err := githashread.Object_Read(repo, sha)
		if err != nil {
			return nil, err
		}
		return obj.Deserialize(), nil
	}
	return prettyprint.RunTreeViewer(rootItems, fetchChildren, fetchBlob, prettyprint.DefaultViewerConfig, header)
}

func buildTreeItems(repo gitpath.GitRepository, sha string) ([]prettyprint.TreeItem, error) {
	obj, err := githashread.Object_Read(repo, sha)
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

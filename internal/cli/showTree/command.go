package showTree

import (
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"
	"strings"

	"github.com/spf13/cobra"
)

func LsTree(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

	// build root items from the given SHA
	rootItems, err := buildTreeItems(*repo, args[0])
	if err != nil {
		panic(err)
	}

	// fetchChildren — called when user enters a subtree
	fetchChildren := func(sha string) ([]prettyprint.TreeItem, error) {
		return buildTreeItems(*repo, sha)
	}

	// fetchBlob — called when user selects a file
	fetchBlob := func(sha string) ([]byte, error) {
		obj, err := githashread.Object_Read(*repo, sha)
		if err != nil {
			return nil, err
		}
		return obj.Deserialize(), nil
	}

	return prettyprint.RunTreeViewer(
		rootItems,
		fetchChildren,
		fetchBlob,
		prettyprint.DefaultViewerConfig,
		"",
	)
}

// buildTreeItems converts GitTreeLeaf slice into TreeItem slice
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

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show-tree [tree-hash]",
		Aliases: []string{"tree"},
		Short:   "Show contents of Tree object",
		Args:    cobra.MaximumNArgs(1),
		RunE:    LsTree,
	}

	return cmd
}

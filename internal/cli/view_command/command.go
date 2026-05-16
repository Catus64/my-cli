package viewcommand

import (
	"bufio"
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitlog "gocmd/testfiles/GitLog"
	gitobj "gocmd/testfiles/GitObject"
	gitshowref "gocmd/testfiles/GitShowRef"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View repository objects, trees, refs and versions",
	}

	cmd.AddCommand(newResolveCommand())
	cmd.AddCommand(newTreeCommand())
	cmd.AddCommand(newObjectCommand())
	cmd.AddCommand(newRefsCommand())
	cmd.AddCommand(newVersionCommand())

	return cmd
}

func newResolveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sha <sha>",
		Short: "View any object by SHA — auto detects type",
		Args:  cobra.ExactArgs(1),
		RunE:  runResolve,
	}
}

func runResolve(cmd *cobra.Command, args []string) error {
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		return err
	}
	return resolveAndView(*repo, args[0])
}

func resolveAndView(repo gitpath.GitRepository, sha string) error {
	obj, err := githashread.Object_Read(repo, sha)
	if err != nil {
		return fmt.Errorf("could not read object %s: %w", sha, err)
	}

	switch obj.Get_Format() {
	case "commit":
		return viewCommit(repo, sha)
	case "tree":
		return viewTree(repo, sha)
	case "blob":
		content := obj.Deserialize()
		prettyprint.PrintObjectContent(sha, content)
		return nil
	case "tag":
		tag, ok := obj.(*gitobj.GitTag)
		if !ok {
			return fmt.Errorf("failed to cast tag")
		}
		tag.Deserialize()
		// show tag fields
		prettyprint.PrintObjectContent(sha, []byte(fmt.Sprintf(
			"object: %s\ntype:   %s\ntag:    %s\ntagger: %s\n\n%s",
			string(tag.KvlmDict.Dict["object"]),
			string(tag.KvlmDict.Dict["type"]),
			string(tag.KvlmDict.Dict["tag"]),
			string(tag.KvlmDict.Dict["tagger"]),
			string(tag.KvlmDict.Dict["data"]),
		)))
		return nil
	default:
		return fmt.Errorf("unknown object type: %s", obj.Get_Format())
	}
}

// ─── view tree <sha> ─────────────────────────────────────────────────

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
	return viewTree(*repo, args[0])
}

func viewTree(repo gitpath.GitRepository, sha string) error {
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

	return prettyprint.RunTreeViewer(rootItems, fetchChildren, fetchBlob, prettyprint.DefaultViewerConfig)
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

// ─── view object <sha> ───────────────────────────────────────────────

func newObjectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "object <sha>",
		Short: "View raw content of any object",
		Args:  cobra.ExactArgs(1),
		RunE:  runObject,
	}
}

func runObject(cmd *cobra.Command, args []string) error {
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		return err
	}
	return viewObject(*repo, args[0])
}

func viewObject(repo gitpath.GitRepository, sha string) error {
	obj, err := githashread.Object_Read(repo, sha)
	if err != nil {
		return fmt.Errorf("could not read object: %w", err)
	}
	prettyprint.PrintObjectContent(sha, obj.Deserialize())
	return nil
}

// ─── view commit (internal) ──────────────────────────────────────────

func viewCommit(repo gitpath.GitRepository, sha string) error {
	obj, err := githashread.Object_Read(repo, sha)
	if err != nil {
		return err
	}
	commit, ok := obj.(*gitobj.GitCommit)
	if !ok {
		return fmt.Errorf("not a commit object")
	}
	commit.Deserialize()

	shortSHA := sha
	if len(sha) > 7 {
		shortSHA = sha[:7]
	}

	parents := gitobj.GetKvlmValues(commit.KvlmDict.Dict, "parent")
	parentDisplay := ""
	if len(parents) > 0 {
		parentDisplay = parents[0][:7]
		if len(parents) > 1 {
			parentDisplay += fmt.Sprintf(" (+%d more)", len(parents)-1)
		}
	}

	date, author := gitlog.Format_Date_Author(string(commit.KvlmDict.Dict["author"]))
	treeSHA := strings.TrimSpace(string(commit.KvlmDict.Dict["tree"]))

	prettyprint.PrintCommit(
		shortSHA,
		author,
		date,
		treeSHA[:7],
		parentDisplay,
		string(commit.KvlmDict.Dict["data"]),
	)

	// prompt to view tree
	fmt.Print("\npress v to view tree, any other key to exit: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadByte()
	if input == 'v' {
		return viewTree(repo, treeSHA)
	}
	return nil
}

// ─── view refs ───────────────────────────────────────────────────────

func newRefsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "refs",
		Short: "Interactively browse all repository refs",
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

	// fetchChildren unused for refs — refs have no children
	fetchChildren := func(sha string) ([]prettyprint.TreeItem, error) {
		return nil, nil
	}

	// selecting a ref → resolve to commit → show commit view + tree option
	fetchBlob := func(sha string) ([]byte, error) {
		// we hijack fetchBlob to show commit view instead
		// return nil content — handled by onSelect override
		return nil, nil
	}

	// use a custom viewer that calls viewCommit on select
	return runRefsViewer(*repo, items, fetchChildren, fetchBlob)
}

func runRefsViewer(
	repo gitpath.GitRepository,
	items []prettyprint.TreeItem,
	fetchChildren func(string) ([]prettyprint.TreeItem, error),
	fetchBlob func(string) ([]byte, error),
) error {
	// override fetchBlob to show commit view
	fetchCommit := func(sha string) ([]byte, error) {
		// restore normal terminal before showing commit
		return nil, nil
	}
	_ = fetchCommit

	onSelect := func(sha string) error {
		return resolveAndView(repo, sha)
	}

	return prettyprint.RunRefsViewer(items, fetchChildren, onSelect, prettyprint.DefaultViewerConfig)
}

// ─── view version (placeholder) ──────────────────────────────────────

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "View versions (coming soon)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Version browser coming soon.")
			return nil
		},
	}
}

package viewcommand

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"

	"github.com/spf13/cobra"
)

func newResolveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sha <sha>",
		Short: "View any object of any type by SHA",
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
		return viewTree(repo, sha, "")
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

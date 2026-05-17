package viewcommand

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"

	"github.com/spf13/cobra"
)

func newObjectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "object <sha>",
		Short: "View content of any BLOB object",
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

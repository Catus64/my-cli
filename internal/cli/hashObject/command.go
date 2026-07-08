package hashObject

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"

	"github.com/spf13/cobra"
)

func hash(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}
	if len(args) == 0 {
		return fmt.Errorf("specify files to hash")
	}
	sha, _ := githashread.Hash_Object(args[0], "blob", *repo)
	prettyprint.PrintObjectStored("blob", args[0], fmt.Sprintf("%x", sha))
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "hash-file [file]",
		Aliases: []string{"hash", "store"},
		Short:   "Hash a file and store it in the repository",
		Args:    cobra.MaximumNArgs(1),
		RunE:    hash,
	}

	return cmd
}

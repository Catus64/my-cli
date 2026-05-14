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
	sha, _ := githashread.Hash_Object(args[0], "blob", *repo)
	// fmt.Printf("file: {%s} has been compressed and stored in the objects folder\n", args[0])
	// fmt.Printf("SHA: %x \n", sha)
	// fmt.Println("To search for the file in the repository please use this hash to find it later")
	prettyprint.PrintObjectStored("blob", args[0], fmt.Sprintf("%x", sha))
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hash-file [object-hash]",
		Short: "Print a friendly greeting",
		Args:  cobra.MaximumNArgs(1),
		RunE:  hash,
	}

	return cmd
}

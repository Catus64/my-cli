package checkout

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func checkout(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

	tree := githashread.Object_Read(*repo, args[0])
	tree.Deserialize()

	// exprerimenting checkout of a single file
	concreteTree, ok := tree.(gitobj.GitTree)
	if !ok {
		panic("not a tree object")
	}
	concreteTree.DeserializeData(tree.Deserialize())
	fmt.Println(concreteTree.Items[6].String())

	path := gitpath.Get_Os_Dir()
	sister_path := filepath.Join(path, "..", "sister_checkout_dir")
	err = os.Mkdir(sister_path, 0o755)
	if err != nil {
	}
	filepath := filepath.Join(sister_path, string(concreteTree.Items[6].Path))

	blob_obj := githashread.Object_Read(*repo, concreteTree.Items[6].Sha)
	blob_data := blob_obj.Deserialize()

	err = os.WriteFile(filepath, blob_data, 0o644)
	if err != nil {
		panic(err)
	}

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load [commit-object]",
		Short: "Load Version",
		Args:  cobra.MaximumNArgs(1),
		RunE:  checkout,
	}

	return cmd
}

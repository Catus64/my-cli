package load

import (
	"fmt"
	loading "gocmd/testfiles/GitCheckout"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

// for testing : f662436f828f1cbe6dd8dff9f2f9935e5af52b3b

func load(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}
	if len(args) < 1 {
		panic("missing required argument <commit-object> sha")
	}
	name := ""
	if len(args) == 2 {
		name = args[1]
	} else {
		name = args[0]
	}

	loaded_path, err := loading.Load(args[0], name, *repo)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading Successful!!")
	fmt.Println("Loaded to:", loaded_path)

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load [commit-object]",
		Short: "Load Version",
		Args:  cobra.MaximumNArgs(2),
		RunE:  load,
	}

	return cmd
}

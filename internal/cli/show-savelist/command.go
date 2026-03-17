package showsavelist

import (
	gitobject "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func show_savelist(cmd *cobra.Command, args []string) error {

	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		panic(err)
	}

	_, err = gitobject.Index_Read(*repo)
	if err != nil {
		panic(err)
	}
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-savelist",
		Short: "Show all files that is about to be saved in the next version",
		Args:  cobra.MaximumNArgs(1),
		RunE:  show_savelist,
	}

	return cmd
}

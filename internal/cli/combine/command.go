package combine

import (
	"fmt"
	save "gocmd/testfiles/GitSave"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func combine(cmd *cobra.Command, args []string) error {

	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		panic(err)
	}

	err = save.Combine(*repo, args[0])
	if err != nil {
		return fmt.Errorf("failed to combine savefiles: %w", err)
	}

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "combine",
		Short: "Combine current version with latest version from another savefile",
		Args:  cobra.MaximumNArgs(1),
		RunE:  combine,
	}

	return cmd
}

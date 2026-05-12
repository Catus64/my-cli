package combine

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func combine(cmd *cobra.Command, args []string) error {

	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		panic(err)
	}
	checkSha := "a067aeb2f5ab8668edc04a9ad5de60ddac1e71d8"
	commit, err := githashread.Object_Read(*repo, checkSha)

	if err != nil {
		return fmt.Errorf("failed to read object: %w", err)
	}

	concreteCommit, ok := commit.(*gitobj.GitCommit)
	if !ok {
		return fmt.Errorf("failed to cast commit to GitCommit")
	}
	concreteCommit.Deserialize()
	parent_sha := concreteCommit.Dict["parent"]

	fmt.Println(string(parent_sha))

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

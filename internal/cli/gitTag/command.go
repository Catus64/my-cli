package gittag

import (
	// "fmt"
	//githashread "gocmd/testfiles/GitHashRead"

	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func show_ref(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

	name := args[0]

	var sha string
	if len(args) > 1 {
		// hash explicitly provided - validate it resolves to a real object
		candidate := args[1]
		if _, err := githashread.Object_Read(*repo, candidate); err != nil {
			return fmt.Errorf("invalid hash: %s", candidate)
		}
		sha = candidate
	} else {
		// no hash provided - default to HEAD
		resolved, err := gitobj.Ref_Resolve(*repo, "HEAD")
		if err != nil || resolved == nil {
			return fmt.Errorf("could not resolve HEAD")
		}
		sha = *resolved
	}

	if err := gitobj.Create_Ref(*repo, name, sha); err != nil {
		return fmt.Errorf("failed to create ref: %w", err)
	}

	fmt.Printf("Tagged %s -> %s\n", name, sha[:7])
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-name <name> [object-hash]",
		Aliases: []string{"set", "name", "tag"},
		Short:   "Set a name for a version (defaults to latest version if no object hash)",
		Args:    cobra.RangeArgs(1, 2),
		RunE:    show_ref,
	}
	return cmd
}

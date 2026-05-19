package checkignore

import (
	"fmt"
	gitcheckignore "gocmd/testfiles/GitCheckIgnore"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"

	"github.com/spf13/cobra"
)

func check_ignore(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

	rules, err := gitcheckignore.ReadGitIgnore(*repo)
	if err != nil {
		return err
	}
	logger.L().Debug("GitIgnore read",
		"absolute", rules.Absolute,
		"scoped", rules.Scoped,
	)

	for _, arg := range args {
		fmt.Println("Checking if", arg, "is ignored...")
		fmt.Println(gitcheckignore.CheckIgnore(rules, arg))
	}

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "check-ignore [object-hash]",
		Aliases: []string{"check", "ci"},
		Short:   "Check if the specified object is ignored by .gitignore",
		RunE:    check_ignore,
	}

	return cmd
}

package checkignore

import (
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

	raw, include, ok := gitcheckignore.Gitignore_Parse1("#asdasdas")

	logger.L().Debug("gitignore parsed",
		"raw", raw,
		"include", include,
		"ok", ok,
	)

	read, err := gitcheckignore.ReadGitIgnore(*repo)
	if err != nil {
		return err
	}
	logger.L().Debug("GitIgnore read",
		"absolute", read.Absolute,
		"scoped", read.Scoped,
	)

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-ignore [object-hash]",
		Short: "Check if the specified object is ignored by .gitignore",
		Args:  cobra.MaximumNArgs(1),
		RunE:  check_ignore,
	}

	return cmd
}

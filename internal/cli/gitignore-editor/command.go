package gitignoreeditor

import (
	gitignoreeditor "gocmd/testfiles/GitCheckIgnore"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func ignoreCmd(cmd *cobra.Command, args []string) error {
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		return err
	}
	return gitignoreeditor.EditGitignore(repo.WorkTree)
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-ignore",
		Aliases: []string{"ignore"},
		Short:   "Interactively edit .gitignore rules",
		RunE:    ignoreCmd,
	}
	return cmd
}

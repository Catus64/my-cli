package switchver

import (
	gitpath "gocmd/testfiles/Gitrepostruct"
	alt_ver "gocmd/testfiles/alternateVersions"

	"github.com/spf13/cobra"
)

func testsomething() error {
	return nil
}

func switch_alt_ver(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		return nil
	}

	alt_ver.SwitchAltVer(*repo, "what")
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [object-hash]",
		Short: "Remove the specified object from the save list",
		RunE:  switch_alt_ver,
	}

	return cmd
}

package switchver

import (
	gitpath "gocmd/testfiles/Gitrepostruct"
	alt_ver "gocmd/testfiles/alternateVersions"

	"github.com/spf13/cobra"
)

func switch_alt_ver(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}

	err = alt_ver.SwitchAltVer(*repo, args[0])
	if err != nil {
		panic(err)
	}
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch-ver [version-name]",
		Short: "Switch to the specified alternate version savefile",
		RunE:  switch_alt_ver,
		Args:  cobra.ExactArgs(1),
	}

	return cmd
}

package remove

import (
	gitaddremove "gocmd/testfiles/GitAddRemove"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func remove(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		return err
	}
	// index, err := gitobject.Index_Read2(*repo)
	// if err != nil {
	// 	return err
	// }

	// original, _ := os.ReadFile(gitpath.Repo_Path(*repo, "index"))

	// gitobject.Index_Write(*repo, *index)

	// result, _ := os.ReadFile(gitpath.Repo_Path(*repo, "index"))

	// logger.L().Debug("Length of bytes",
	// 	"original", len(original), "result", len(result))

	// fmt.Println("match:", bytes.Equal(result, original))

	err = gitaddremove.Remove(repo, []string{"somefile.txt"}, gitaddremove.RemoveOptions{
		Delete:          true,
		SkipMissingFile: false,
	})
	if err != nil {
		return err
	}

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [object-hash]",
		Short: "Remove the specified object from the save list",
		Args:  cobra.MaximumNArgs(1),
		RunE:  remove,
	}

	return cmd
}

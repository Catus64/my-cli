package remove

import (
	"bytes"
	"fmt"
	gitobject "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	"os"

	"github.com/spf13/cobra"
)

func remove(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		return err
	}
	index, err := gitobject.Index_Read2(*repo)
	if err != nil {
		return err
	}

	original, _ := os.ReadFile(gitpath.Repo_Path(*repo, "index"))

	gitobject.Index_Write(*repo, *index)

	result, _ := os.ReadFile(gitpath.Repo_Path(*repo, "index"))

	logger.L().Debug("Length of bytes",
		"original", len(original), "result", len(result))

	fmt.Println("match:", bytes.Equal(result, original))
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [object-hash]",
		Short: "Remove the specified object from the save list",
		RunE:  remove,
	}

	return cmd
}

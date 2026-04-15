package gitaddremove

import (
	"fmt"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
)

func Remove_entry(repo gitpath.GitRepository, paths []string, del bool, missingSkip bool) error {
	_, err := gitobj.Index_Read2(repo)
	if err != nil {
		return err
	}

	worktree := repo.WorkTree + string(os.PathSeparator)
	fmt.Println(worktree)

	return nil
}

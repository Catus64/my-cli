package gitaddremove

import (
	"fmt"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
	"path/filepath"
	"strings"
)

type RemoveOptions struct {

	// for deleting the file from worktree
	Delete bool

	// skip missing file -> no error
	SkipMissingFile bool
}

func Remove(repo *gitpath.GitRepository, paths []string, options RemoveOptions) error {

	// read index
	index, err := gitobj.Index_Read2(*repo)
	if err != nil {
		return err
	}

	worktree := repo.WorkTree + string(os.PathSeparator)

	fmt.Println("index entries before remove:", len(index.Entries))
	// fmt.Println("index: ", index)

	fmt.Println("worktree: ", worktree)

	// make absolute path

	absolutePaths := make(map[string]struct{})

	for _, path := range paths {
		absolutePath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("unable to resolve path %s: %w", path, err)
		}

		if !strings.HasPrefix(absolutePath, worktree) {
			return fmt.Errorf("path %s is outside the repository worktree", path)
		}

		absolutePaths[absolutePath] = struct{}{}

	}

	return nil
}

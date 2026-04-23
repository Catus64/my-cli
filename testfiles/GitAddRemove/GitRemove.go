package gitaddremove

import (
	"fmt"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	log "gocmd/testfiles/Helper"
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

	var keptEntries []gitobj.GitIndexEntry
	var toRemove []string

	for _, entry := range index.Entries {
		fullpath := filepath.Join(repo.WorkTree, entry.Name)

		if _, found := absolutePaths[fullpath]; found {
			toRemove = append(toRemove, fullpath)
			delete(absolutePaths, fullpath)
		} else {
			keptEntries = append(keptEntries, entry)
		}
	}

	if len(absolutePaths) > 0 && !options.SkipMissingFile {
		missing := make([]string, 0, len(absolutePaths))
		for path := range absolutePaths {
			missing = append(missing, path)
		}

		return fmt.Errorf("Cannot remove paths not in the index: %v", missing)
	}

	if options.Delete {
		for _, path := range toRemove {
			log.L().Debug("deleting file from filesystem", "path", path)

			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to delete file %s, %w", path, err)
			}
		}
	}

	index.Entries = keptEntries

	fmt.Println("Index after remove", len(keptEntries))

	for _, index := range keptEntries {
		fmt.Println(index.Name)
	}

	return nil
}

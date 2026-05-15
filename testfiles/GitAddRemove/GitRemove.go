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

func Remove(repo *gitpath.GitRepository, paths []string, options RemoveOptions) (*gitobj.GitIndex, error) {

	// read index
	index, err := gitobj.Index_Read2(*repo)
	if err != nil {
		return nil, err
	}

	worktree := filepath.Clean(repo.WorkTree) + string(os.PathSeparator)

	// fmt.Println("index: ", index)

	// make absolute path

	absolutePaths := make(map[string]struct{})

	for _, path := range paths {
		absolutePath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("unable to resolve path %s: %w", path, err)
		}

		// Normalize absolute path
		absolutePath = filepath.Clean(absolutePath)

		fmt.Println("worktree:", worktree)
		fmt.Println("absolutePath:", absolutePath)
		fmt.Println("hasPrefix:", strings.HasPrefix(absolutePath, worktree))

		if !strings.HasPrefix(absolutePath, worktree) {
			return nil, fmt.Errorf("path %s is outside the repository worktree", path)
		}

		absolutePaths[absolutePath] = struct{}{}

	}

	var keptEntries []gitobj.GitIndexEntry
	var toRemove []string

	for _, entry := range index.Entries {
		fullpath := filepath.Clean(filepath.Join(repo.WorkTree, entry.Name))

		// for ap := range absolutePaths {
		// 	fmt.Printf("absolutePath   : %q\n", ap)
		// 	fmt.Printf("match?         : %v\n", fullpath == ap)
		// }

		if _, found := absolutePaths[fullpath]; found {
			toRemove = append(toRemove, fullpath)
			// fmt.Printf("%s has been removed from the Savelist \n\n", fullpath)
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

		return nil, fmt.Errorf("Cannot remove paths not in the index: %v", missing)
	}

	if options.Delete {
		for _, path := range toRemove {
			log.L().Debug("deleting file from filesystem", "path", path)
			fmt.Printf("%s has been deleted from the filesystem\n\n", path)
			if err := os.Remove(path); err != nil {
				return nil, fmt.Errorf("failed to delete file %s, %w", path, err)
			}
		}
	}

	// index.Entries = keptEntries

	// for _, index := range keptEntries {
	// 	fmt.Println(index.Name)
	// }

	index.Entries = keptEntries
	if err := gitobj.Index_Write(*repo, *index); err != nil {
		return nil, fmt.Errorf("failed to write index: %w", err)
	}

	return index, nil
}

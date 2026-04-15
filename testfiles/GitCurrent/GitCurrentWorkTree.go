package gitcurrent

import (
	"fmt"
	gitcheckignore "gocmd/testfiles/GitCheckIgnore"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

func StatusIndexWorktree(repo gitpath.GitRepository, index gitobj.GitIndex) ([]string, error) {
	fmt.Println("Changes not staged for commit:")

	// Load gitignore rules
	rules, err := gitcheckignore.ReadGitIgnore(repo)
	if err != nil {
		return nil, err
	}

	// Walk filesystem and collect all relative file paths
	var allFiles []string
	err = filepath.WalkDir(repo.WorkTree, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip .git directory entirely
		if d.IsDir() && path == repo.GitDir {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			rel, err := filepath.Rel(repo.WorkTree, path)
			if err != nil {
				return err
			}
			allFiles = append(allFiles, rel)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Compare index entries against real files on disk
	for _, entry := range index.Entries {
		fullPath := filepath.Join(repo.WorkTree, entry.Name)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// File in index but gone from disk
			fmt.Println("  deleted:  ", entry.Name)
		} else {
			// File exists — check timestamps first
			info, _ := os.Stat(fullPath)
			stat := info.Sys().(*syscall.Stat_t)

			entryMtime := int64(entry.MtimeSec)*1e9 + int64(entry.MtimeNano)
			actualMtime := stat.Mtim.Sec*1e9 + stat.Mtim.Nsec

			logger.L().Debug("Comparing index entry to disk",
				"entry", entry.Name,
				"indexMtime", entryMtime,
				"actualMtime", actualMtime,
			)
			if entryMtime != actualMtime {
				// Timestamps differ — hash to confirm
				newSHA, err := githashread.Hash_Object_NoWrite(fullPath, "blob")
				if err != nil {
					return nil, err
				}
				logger.L().Debug("Hashing file to confirm modification",
					"entry", entry.Name,
					"indexSHA", entry.SHA,
					"newSHA", fmt.Sprintf("%x", newSHA),
				)
				if fmt.Sprintf("%x", newSHA) != entry.SHA {
					fmt.Println("  modified: ", entry.Name)
				}
			}
		}

		// Remove from allFiles — it's accounted for
		for i, f := range allFiles {
			if f == entry.Name {
				allFiles = append(allFiles[:i], allFiles[i+1:]...)
				break
			}
		}
	}
	logger.L().Debug("All files on disk after processing index", "allFiles", allFiles)

	var untrackedFiled []string
	// Whatever remains in allFiles is untracked
	for _, f := range allFiles {
		ignored, err := gitcheckignore.CheckIgnore(rules, f)
		if err != nil || ignored {
			continue
		}
		untrackedFiled = append(untrackedFiled, f)
	}

	return untrackedFiled, nil
}

package gitcurrent

import (
	"fmt"
	gitcheckignore "gocmd/testfiles/GitCheckIgnore"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

func StatusIndexWorktree(repo gitpath.GitRepository, index gitobj.GitIndex) (*IndexWorkTree, error) {
	result := IndexWorkTree{}

	rules, err := gitcheckignore.ReadGitIgnore(repo)
	if err != nil {
		return &result, err
	}

	// collect all files on disk
	var allFiles []string
	err = filepath.WalkDir(repo.WorkTree, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
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
		return &result, err
	}

	// compare index against disk
	for _, entry := range index.Entries {
		fullPath := filepath.Join(repo.WorkTree, entry.Name)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			result.Deleted = append(result.Deleted, entry.Name)
		} else {
			info, _ := os.Stat(fullPath)
			stat := info.Sys().(*syscall.Stat_t)
			entryMtime := int64(entry.MtimeSec)*1e9 + int64(entry.MtimeNano)
			actualMtime := stat.Mtim.Sec*1e9 + stat.Mtim.Nsec

			if entryMtime != actualMtime {
				newSHA, err := githashread.Hash_Object_NoWrite(fullPath, "blob")
				if err != nil {
					return &result, err
				}
				if fmt.Sprintf("%x", newSHA) != entry.SHA {
					result.Modified = append(result.Modified, entry.Name)
				}
			}
		}

		// remove from allFiles — accounted for
		for i, f := range allFiles {
			if f == entry.Name {
				allFiles = append(allFiles[:i], allFiles[i+1:]...)
				break
			}
		}
	}

	// remaining files are untracked
	for _, f := range allFiles {
		ignored, err := gitcheckignore.CheckIgnore(rules, f)
		if err != nil || ignored {
			continue
		}
		result.Untracked = append(result.Untracked, f)
	}

	return &result, nil
}

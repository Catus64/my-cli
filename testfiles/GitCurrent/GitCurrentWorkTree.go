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
)

func StatusIndexWorktree(repo gitpath.GitRepository, index gitobj.GitIndex) (*IndexWorkTree, error) {
    result := IndexWorkTree{}

    rules, err := gitcheckignore.ReadGitIgnore(repo)
    if err != nil {
        return &result, err
    }

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
            // Normalize to forward slashes to match git index format
            allFiles = append(allFiles, filepath.ToSlash(rel))
        }
        return nil
    })
    if err != nil {
        return &result, err
    }

    for _, entry := range index.Entries {
        // entry.Name is always forward-slashed (git format)
        // build OS path only for disk operations
        fullPath := filepath.Join(repo.WorkTree, filepath.FromSlash(entry.Name))

        if _, err := os.Stat(fullPath); os.IsNotExist(err) {
            result.Deleted = append(result.Deleted, entry.Name)
        } else {
            info, _ := os.Stat(fullPath)
            actualMtime := info.ModTime().UnixNano()
            entryMtime := int64(entry.MtimeSec)*1e9 + int64(entry.MtimeNano)

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

        // Now this comparison works — both sides are forward-slashed
        for i, f := range allFiles {
            if f == entry.Name {
                allFiles = append(allFiles[:i], allFiles[i+1:]...)
                break
            }
        }
    }

    for _, f := range allFiles {
        ignored, err := gitcheckignore.CheckIgnore(rules, f)
        if err != nil || ignored {
            continue
        }
        result.Untracked = append(result.Untracked, f)
    }

    return &result, nil
}
package gitaddremove

import (
	"fmt"
	gitcheckignore "gocmd/testfiles/GitCheckIgnore"
	hashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
	"path/filepath"
	"strings"
)

// option if user wants to add all
type Options struct {
	All bool
}

// absolute path is for the real system path
// clean is relative to .git parent (worktree)
type cleanPath struct {
	absolute string
	relative string
}

//context switch to return to command.go

type AddedFile struct {
	Path   string
	SHA    string
	Status string // New | Modified
}

type AddResult struct {
	Added   []AddedFile
	Ignored []string
}

func Add(repo *gitpath.GitRepository, paths []string, opts Options) (*AddResult, error) {

	result := AddResult{}

	// if all collect all files in working directory
	if opts.All {
		Allpaths, err := getAllFiles(*repo)
		if err != nil {
			return nil, err
		}
		// Filter out gitignored files
		rules, err := gitcheckignore.ReadGitIgnore(*repo)
		if err != nil {
			return nil, fmt.Errorf("failed to read gitignore: %w", err)
		}

		var filtered []string
		for _, path := range Allpaths {
			// CheckIgnore needs relative path
			rel, err := filepath.Rel(repo.WorkTree, path)
			if err != nil {
				fmt.Printf("Error processing %s: %v\n", path, err)
				continue
			}
			ignored, err := gitcheckignore.CheckIgnore(rules, rel)
			if err != nil {
				fmt.Printf("Error checking ignore for %s: %v\n", rel, err)
				continue
			}
			if ignored {
				continue
			}
			filtered = append(filtered, path)
		}
		paths = filtered
	}

	//dest -> dest/ explicit end and work with both os
	worktree := repo.WorkTree + string(os.PathSeparator)

	curIndex, err := gitobj.Index_Read2(*repo)
	if err != nil {
		return nil, fmt.Errorf("failed ")
	}

	curSHAs := make(map[string]string)
	for _, entry := range curIndex.Entries {
		curSHAs[entry.Name] = entry.SHA
	}

	//remove first before adding
	//easy way to avoid duplicates
	_, err = Remove(repo, paths, RemoveOptions{
		Delete:          false,
		SkipMissingFile: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to clear existing index entries: %w", err)
	}

	var cleanPaths []cleanPath

	for _, path := range paths {
		absolutePath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve path %s: %w", path, err)
		}

		// check if path is inside repositroy worktree
		if !strings.HasPrefix(absolutePath, worktree) {
			return nil, fmt.Errorf("path %s is outside the working directory", path)
		}

		info, err := os.Stat(absolutePath)
		if err != nil || info.IsDir() {
			return nil, fmt.Errorf("not a file or does not exist %s: %w", path, err)
		}

		//Rel is getting rel path from repo worktree
		relativePath, err := filepath.Rel(repo.WorkTree, absolutePath)
		if err != nil {
			return nil, fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}

		cleanPaths = append(cleanPaths, cleanPath{
			absolute: absolutePath,
			relative: relativePath,
		})
	}

	index, err := gitobj.Index_Read2(*repo)
	if err != nil {
		return nil, fmt.Errorf("failed to read index: %w", err)
	}

	for _, clean_path := range cleanPaths {
		// store blobs in /object
		sha, err := hashread.Hash_Object(clean_path.absolute, "blob", *repo)
		if err != nil {
			return nil, fmt.Errorf("failed to hash file %s: %w", clean_path.relative, err)
		}
		shaHex := fmt.Sprintf("%x", sha)

		stat, err := os.Stat(clean_path.absolute)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file %s: %w", clean_path.relative, err)
		}

		//build index entry the append
		entry, err := buildIndexEntry(clean_path.relative, fmt.Sprintf("%x", sha), stat)
		if err != nil {
			return nil, fmt.Errorf("failed to build index entry for %s: %w", clean_path.relative, err)
		}

		index.Entries = append(index.Entries, entry)
		// fmt.Printf("Adding file: %s (relative path: %s)\n", clean_path.absolute, clean_path.relative)
		// fmt.Printf("Hash for %s: %x\n", clean_path.relative, sha)

		//new or modified
		status := "new"
		if oldSHA, exists := curSHAs[clean_path.relative]; exists && oldSHA != shaHex {
			status = "modified"
		} else if exists {
			// SHA unchanged entries
			continue
		}

		result.Added = append(result.Added, AddedFile{
			Path:   clean_path.relative,
			SHA:    shaHex,
			Status: status,
		})

	}

	if err := gitobj.Index_Write(*repo, *index); err != nil {
		return nil, fmt.Errorf("failed to write index: %w", err)
	}

	return &result, nil
}

func buildIndexEntry(relpath, sha string, stat os.FileInfo) (gitobj.GitIndexEntry, error) {
	mod := stat.ModTime()

	var ctimeSec, ctimeNsec, mtimeSec, mtimeNsec uint32
	var dev, ino, uid, gid uint32

	mtimeSec = uint32(mod.Unix())
	mtimeNsec = uint32(mod.Nanosecond())

	ctimeSec = uint32(mod.Unix())
	ctimeNsec = uint32(mod.Nanosecond())
	dev = 0
	ino = 0
	uid = 0
	gid = 0

	mode := stat.Mode()
	modeType := uint16(0b1000) // blob file
	modePerms := uint16(mode.Perm()) & 0o777

	return gitobj.GitIndexEntry{
		CtimeSec:    ctimeSec,
		CtimeNano:   ctimeNsec,
		MtimeSec:    mtimeSec,
		MtimeNano:   mtimeNsec,
		Dev:         dev,
		Ino:         ino,
		ModeType:    modeType,
		ModePerms:   modePerms,
		UID:         uid,
		GID:         gid,
		FSize:       uint32(stat.Size()),
		SHA:         sha,
		AssumeValid: false,
		Stage:       0,
		Name:        relpath,
	}, nil
}

func getAllFiles(repo gitpath.GitRepository) ([]string, error) {
	var files []string

	//anonymous function here is to deal with walk result
	err := filepath.Walk(repo.WorkTree, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

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
	"syscall"
	"time"
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

func Add(repo *gitpath.GitRepository, paths []string, opts Options) error {
	// if all collect all files in working directory
	if opts.All {
		Allpaths, err := getAllFiles(*repo)
		if err != nil {
			return err
		}
		// Filter out gitignored files
		rules, err := gitcheckignore.ReadGitIgnore(*repo)
		if err != nil {
			return fmt.Errorf("failed to read gitignore: %w", err)
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

	//remove first before adding
	//easy way to avoid duplicates
	_, err := Remove(repo, paths, RemoveOptions{
		Delete:          false,
		SkipMissingFile: true,
	})
	if err != nil {
		return fmt.Errorf("failed to clear existing index entries: %w", err)
	}

	var cleanPaths []cleanPath

	for _, path := range paths {
		absolutePath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to resolve path %s: %w", path, err)
		}

		// check if path is inside repositroy worktree
		if !strings.HasPrefix(absolutePath, worktree) {
			return fmt.Errorf("path %s is outside the working directory", path)
		}

		info, err := os.Stat(absolutePath)
		if err != nil || info.IsDir() {
			return fmt.Errorf("not a file or does not exist %s: %w", path, err)
		}

		//Rel is getting rel path from repo worktree
		relativePath, err := filepath.Rel(repo.WorkTree, absolutePath)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}

		cleanPaths = append(cleanPaths, cleanPath{
			absolute: absolutePath,
			relative: relativePath,
		})
	}

	index, err := gitobj.Index_Read2(*repo)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	for _, clean_path := range cleanPaths {
		// store blobs in /object
		sha, err := hashread.Hash_Object(clean_path.absolute, "blob", *repo)
		if err != nil {
			return fmt.Errorf("failed to hash file %s: %w", clean_path.relative, err)
		}

		stat, err := os.Stat(clean_path.absolute)
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %w", clean_path.relative, err)
		}

		//build index entry the append
		entry, err := buildIndexEntry(clean_path.relative, fmt.Sprintf("%x", sha), stat)
		if err != nil {
			return fmt.Errorf("failed to build index entry for %s: %w", clean_path.relative, err)
		}

		index.Entries = append(index.Entries, entry)
		// fmt.Printf("Adding file: %s (relative path: %s)\n", clean_path.absolute, clean_path.relative)
		// fmt.Printf("Hash for %s: %x\n", clean_path.relative, sha)

	}

	if err := gitobj.Index_Write(*repo, *index); err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}

	return nil
}

func buildIndexEntry(relpath, sha string, stat os.FileInfo) (gitobj.GitIndexEntry, error) {
	mod := stat.ModTime()
	ctime := stat.Sys()

	var ctimeSec, ctimeNsec, mtimeSec, mtimeNsec uint32
	var dev, ino, uid, gid uint32

	mtimeSec = uint32(mod.Unix())
	mtimeNsec = uint32(mod.Nanosecond())

	// Read platform-specific fields from syscall.Stat_t
	// This is mainly to handle cases with windows vs linux(unix)
	// windows doesn't have ctime and bunch other stuff
	if sysStat, ok := ctime.(*syscall.Stat_t); ok {
		ctimeSec = uint32(sysStat.Ctim.Sec)
		ctimeNsec = uint32(sysStat.Ctim.Nsec)
		dev = uint32(sysStat.Dev)
		ino = uint32(sysStat.Ino)
		uid = sysStat.Uid
		gid = sysStat.Gid
	} else {
		// values for windows lmao
		ctimeSec = uint32(time.Now().Unix())
		ctimeNsec = 0
		dev = 0
		ino = 0
		uid = 0
		gid = 0
	}

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

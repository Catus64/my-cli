package gitcheckignore

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gitpath "gocmd/testfiles/Gitrepostruct"
)

// GitIgnore holds all ignore rules return rules from all sources of ignore
type GitIgnore struct {
	Absolute []string            // local + global ignores
	Scoped   map[string][]string // directory -> patterns from staged .gitignore
}

// ReadGitIgnore reads all ignore rules for a repo
func ReadGitIgnore(repo gitpath.GitRepository) (*GitIgnore, error) {
	ret := &GitIgnore{
		Absolute: []string{},
		Scoped:   make(map[string][]string),
	}

	// Scoped ignore rules from .git/info/exclude
	localFile := filepath.Join(repo.GitDir, "info", "exclude")
	if _, err := os.Stat(localFile); err == nil {
		data, err := os.ReadFile(localFile)
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(data), "\n")
		ret.Absolute = append(ret.Absolute, Gitignore_Parse(lines)...)
	}

	// Global ignore rules from XDG config home
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		configHome = filepath.Join(home, ".config")
	}

	// Global rules from ~/.config/git/ignore
	globalFile := filepath.Join(configHome, "git", "ignore")
	if _, err := os.Stat(globalFile); err == nil {
		data, err := os.ReadFile(globalFile)
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(data), "\n")
		ret.Absolute = append(ret.Absolute, Gitignore_Parse(lines)...)
	}

	// Walk the Directory to search for .gitignores
	err := filepath.WalkDir(repo.WorkTree, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// error reading files will terminate walk
			return err
		}
		if d.IsDir() && path == repo.GitDir {
			// skip .git folder
			return filepath.SkipDir
		}
		if !d.IsDir() && d.Name() == ".gitignore" {
			// if file is gitignore
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// find relative path to where the .gitignore is
			rel, err := filepath.Rel(repo.WorkTree, filepath.Dir(path))
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel) // OS compatibiility
			lines := strings.Split(string(data), "\n")
			// store the relative path (rel) alongside the parsed rules
			ret.Scoped[rel] = append(ret.Scoped[rel], Gitignore_Parse(lines)...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

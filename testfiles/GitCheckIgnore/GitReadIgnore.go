package gitcheckignore

import (
	"os"
	"path/filepath"
	"strings"

	read "gocmd/testfiles/GitHashRead"
	object "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
)

// GitIgnore holds all ignore rules
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

	// Scoped ignore rules from staged .gitignore files in the index
	idx, err := object.Index_Read2(repo)
	if err != nil {
		return nil, err
	}

	for _, entry := range idx.Entries {
		if entry.Name == ".gitignore" || strings.HasSuffix(entry.Name, "/.gitignore") {
			dirName := filepath.Dir(entry.Name)

			obj, err := read.Object_Read(repo, entry.SHA)
			if err != nil {
				return nil, err
			}

			lines := strings.Split(string(obj.Deserialize()), "\n")
			ret.Scoped[dirName] = Gitignore_Parse(lines)
		}
	}

	return ret, nil
}

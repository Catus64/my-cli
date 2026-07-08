package gitsave

import (
	_ "embed"
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

//go:embed words.txt
var wordFile string
var wordList []string

type VersionEntry struct {
	Number int
	Name   string
	SHA    string
}

func (v VersionEntry) ShortName() string {
	return fmt.Sprintf("v%d · %s", v.Number, v.Name)
}

// World file and world list is embedded to this program
// to make excecutable self contained and not rely on external files

func loadWords() {
	if len(wordList) > 0 {
		return
	}
	lines := strings.Split(strings.TrimSpace(wordFile), "\n")
	for _, line := range lines {
		word := strings.TrimSpace(line)
		if word != "" {
			wordList = append(wordList, word)
		}
	}
}

func generateName() string {
	loadWords()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	w1 := wordList[rng.Intn(len(wordList))]
	w2 := wordList[rng.Intn(len(wordList))]
	return w1 + "-" + w2
}

// pathing related stuff

// .git/ezgit/saves
func savesDir(repo gitpath.GitRepository) string {
	return filepath.Join(repo.GitDir, "ezgit", "saves")
}

// .git/ezgit/saves/branch/sha
func branchSavesDir(repo gitpath.GitRepository, branch string) string {
	return filepath.Join(savesDir(repo), branch)
}

// .git/ezgit/saves/branch_counter kinda dumb way to keep track
// but it works
func counterFile(repo gitpath.GitRepository, branch string) string {
	return filepath.Join(savesDir(repo), branch+"_ver.count")
}

//counter stuff

func readCounter(repo gitpath.GitRepository, branch string) (int, error) {
	data, err := os.ReadFile(counterFile(repo, branch))
	if os.IsNotExist(err) {
		return 0, nil // first commit counter starts at 0
	}
	if err != nil {
		return 0, fmt.Errorf("failed to read counter: %w", err)
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, fmt.Errorf("malformed counter file: %w", err)
	}
	return n, nil
}

func writeCounter(repo gitpath.GitRepository, branch string, n int) error {
	return os.WriteFile(counterFile(repo, branch), []byte(strconv.Itoa(n)), 0644)
}

func nextVersion(repo gitpath.GitRepository, branch string) (int, error) {
	current, err := readCounter(repo, branch)
	if err != nil {
		return 0, err
	}
	next := current + 1
	if err := writeCounter(repo, branch, next); err != nil {
		return 0, err
	}
	return next, nil
}

// Actual writing logic

func WriteVersionRef(repo gitpath.GitRepository, branch string, sha string) (VersionEntry, error) {
	// if branch dont exist = make one
	dir := branchSavesDir(repo, branch)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return VersionEntry{}, fmt.Errorf("failed to create saves dir: %w", err)
	}

	num, err := nextVersion(repo, branch)
	if err != nil {
		return VersionEntry{}, err
	}

	//name generation result will be word-word
	name := generateName()

	// file named after SHA, content is "number\nname\n"
	content := fmt.Sprintf("%d\n%s\n", num, name)
	if err := os.WriteFile(filepath.Join(dir, sha), []byte(content), 0644); err != nil {
		return VersionEntry{}, fmt.Errorf("failed to write version ref: %w", err)
	}

	entry := VersionEntry{Number: num, Name: name, SHA: sha}
	logger.L().Info("Version reference created", "version", entry.ShortName(), "sha", sha[:7])
	return entry, nil
}

func ReadVersionRef(repo gitpath.GitRepository, branch string, sha string) (VersionEntry, error) {
	path := filepath.Join(branchSavesDir(repo, branch), sha)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return VersionEntry{}, fmt.Errorf("no version ref for %s on branch %s", sha[:7], branch)
		}
		return VersionEntry{}, fmt.Errorf("failed to read version ref: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) < 2 {
		return VersionEntry{}, fmt.Errorf("malformed version ref for %s", sha[:7])
	}

	num, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return VersionEntry{}, fmt.Errorf("invalid version number: %w", err)
	}

	return VersionEntry{
		Number: num,
		Name:   strings.TrimSpace(lines[1]),
		SHA:    sha,
	}, nil
}

func ListVersionRefs(repo gitpath.GitRepository, branch string) ([]VersionEntry, error) {
	dir := branchSavesDir(repo, branch)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []VersionEntry{}, nil // no versions yet
		}
		return nil, fmt.Errorf("failed to read saves dir: %w", err)
	}

	var versions []VersionEntry
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		v, err := ReadVersionRef(repo, branch, entry.Name())
		if err != nil {
			logger.L().Debug("skipping malformed version ref", "file", entry.Name())
			continue
		}
		versions = append(versions, v)
	}

	// sort by number ascending when returning list
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Number < versions[j].Number
	})

	return versions, nil
}

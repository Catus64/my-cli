package alternateversions

import (
	"fmt"
	gitCurrent "gocmd/testfiles/GitCurrent"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"io/fs"
	"os"
	"path/filepath"
)

func CreateAltVer(repo gitpath.GitRepository, name string) error {
	headSHA, err := gitobj.Ref_Resolve(repo, "HEAD")
	if err != nil || headSHA == nil {
		return fmt.Errorf("no commits yet — save a version first before creating an alternate")
	}

	refPath := gitpath.Repo_Path(repo, "refs", "heads", name)

	// handle nested dirs like feature/commit
	if err := os.MkdirAll(filepath.Dir(refPath), 0755); err != nil {
		return fmt.Errorf("failed to create ref directory: %w", err)
	}

	if _, err := os.Stat(refPath); err == nil {
		return fmt.Errorf("alternate version %q already exists", name)
	}

	if err := os.WriteFile(refPath, []byte(*headSHA+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to create alternate version: %w", err)
	}

	fmt.Printf("Created alternate version %q from current version (%s)\n", name, (*headSHA)[:7])
	return nil
}

func ListAltVer(repo gitpath.GitRepository) error {
	activeBranch, _ := gitCurrent.Get_Active_Branch(repo)

	headsDir := gitpath.Repo_Path(repo, "refs", "heads")

	fmt.Println("Alternate versions:")
	// WalkDir handles nested branches like feature/commit
	err := filepath.WalkDir(headsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil // if directory, recurse into it
		}

		// Get branch name relative to heads dir
		name, err := filepath.Rel(headsDir, path)
		if err != nil {
			return err
		}
		// window compatibility: convert backslashes to slashes in branch names
		name = filepath.ToSlash(name)

		marker := "  "
		if name == activeBranch {
			marker = "* "
		}

		ref := "refs/heads/" + name
		sha, err := gitobj.Ref_Resolve(repo, ref)
		if err != nil || sha == nil {
			return nil // skip unresolvable refs
		}

		shortSHA := *sha
		if len(shortSHA) >= 7 {
			shortSHA = shortSHA[:7]
		}

		fmt.Printf("%s%-30s  %s\n", marker, name, shortSHA)

		return nil
	})
	fmt.Println("")
	fmt.Println(`If you feel like any Branches/Alternate Versions are missing 
try pulling from remote first with "git pull"(if you are using git along with this tool)`)

	return err
}

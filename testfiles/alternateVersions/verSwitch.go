package alternateversions

import (
	"bufio"
	"fmt"
	add "gocmd/testfiles/GitAddRemove"
	loading "gocmd/testfiles/GitCheckout"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func SwitchAltVer(repo gitpath.GitRepository, name string) error {
	// 1. Check branch exists
	refPath := gitpath.Repo_Path(repo, "refs", "heads", name)
	if _, err := os.Stat(refPath); os.IsNotExist(err) {
		return fmt.Errorf("alternate version %q does not exist", name)
	}

	// 2. Resolve branch SHA
	branchSHA, err := gitobj.Ref_Resolve(repo, "refs/heads/"+name)
	if err != nil || branchSHA == nil {
		return fmt.Errorf("failed to resolve alternate version %q: %w", name, err)
	}

	// 3. Dirty check — warn user
	index, err := gitobj.Index_Read2(repo)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}
	if isDirty(repo, *index) {
		fmt.Println("WARNING : You have unsaved changes that will be lost if you switch.")
		fmt.Print("    Continue anyway? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" {
			fmt.Println("Switch cancelled.")
			return nil
		}
	}

	// 4. Backup worktree to ../ezgit_backup
	backupPath := filepath.Join(filepath.Dir(repo.WorkTree), "ezgit_backup")
	fmt.Println("Backing up current worktree to", backupPath)
	if err := copyDir(repo.WorkTree, backupPath, repo.GitDir); err != nil {
		return fmt.Errorf("failed to backup worktree: %w", err)
	}

	// 5. Clear worktree (except .git)
	if err := clearWorktree(repo); err != nil {
		// restore backup before returning error
		_ = copyDir(backupPath, repo.WorkTree, "")
		return fmt.Errorf("failed to clear worktree: %w", err)
	}

	// 6. Checkout branch tree into worktree
	commit, err := githashread.Object_Read(repo, *branchSHA)
	if err != nil {
		_ = copyDir(backupPath, repo.WorkTree, "")
		return fmt.Errorf("failed to read commit: %w", err)
	}
	concreteCommit, ok := commit.(*gitobj.GitCommit)
	if !ok {
		_ = copyDir(backupPath, repo.WorkTree, "")
		return fmt.Errorf("not a commit object")
	}
	concreteCommit.Deserialize()
	treeSHA := string(concreteCommit.KvlmDict.Dict["tree"])

	tree, err := githashread.Object_Read(repo, treeSHA)
	if err != nil {
		_ = copyDir(backupPath, repo.WorkTree, "")
		return fmt.Errorf("failed to read tree: %w", err)
	}
	concreteTree, ok := tree.(gitobj.GitTree)
	if !ok {
		_ = copyDir(backupPath, repo.WorkTree, "")
		return fmt.Errorf("not a tree object")
	}
	concreteTree.DeserializeData(tree.Deserialize())

	if err := loading.TreeCheckout(repo, concreteTree, repo.WorkTree); err != nil {
		// restore backup on failure
		_ = copyDir(backupPath, repo.WorkTree, "")
		return fmt.Errorf("failed to checkout tree: %w", err)
	}

	// 7. Update HEAD to point to new branch
	headPath := gitpath.Repo_Path(repo, "HEAD")
	headContent := fmt.Sprintf("ref: refs/heads/%s\n", name)
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	// 8. Rebuild index via add --all
	if err := add.Add(&repo, nil, add.Options{All: true}); err != nil {
		return fmt.Errorf("failed to rebuild index: %w", err)
	}

	// 9. Delete backup on success
	os.RemoveAll(backupPath)

	fmt.Printf("Switched to alternate version %q\n", name)
	return nil
}

// isDirty checks if any index entry has a different mtime than the real file
func isDirty(repo gitpath.GitRepository, index gitobj.GitIndex) bool {
	for _, entry := range index.Entries {
		fullPath := filepath.Join(repo.WorkTree, entry.Name)
		info, err := os.Stat(fullPath)
		if err != nil {
			return true // file deleted = dirty
		}
		if uint32(info.ModTime().Unix()) != entry.MtimeSec {
			return true
		}
	}
	return false
}

// clearWorktree removes all files/dirs in worktree except .git
func clearWorktree(repo gitpath.GitRepository) error {
	entries, err := os.ReadDir(repo.WorkTree)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}
		fullPath := filepath.Join(repo.WorkTree, entry.Name())
		if err := os.RemoveAll(fullPath); err != nil {
			return err
		}
	}
	return nil
}

// copyDir recursively copies src to dst, skipping skipPath
func copyDir(src, dst, skipPath string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if skipPath != "" && path == skipPath {
			return filepath.SkipDir
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, data, 0644)
	})
}

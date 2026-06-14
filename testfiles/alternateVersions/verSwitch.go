package alternateversions

import (
	"fmt"
	add "gocmd/testfiles/GitAddRemove"
	gitcheckignore "gocmd/testfiles/GitCheckIgnore"
	loading "gocmd/testfiles/GitCheckout"
	gitCurrent "gocmd/testfiles/GitCurrent"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"io/fs"
	"os"
	"path/filepath"
)

func SwitchAltVer(repo gitpath.GitRepository, name string) error {
	//check if Savefile is same as current Savefile, if so do nothing
	activeBranch, err := gitCurrent.Get_Active_Branch(repo)
	if err != nil {
		return fmt.Errorf("failed to get active branch: %w", err)
	}
	if activeBranch == name {
		fmt.Printf("Already on Savefile: %q\n", name)
		return nil
	}

	// check Savefile exists
	refPath := gitpath.Repo_Path(repo, "refs", "heads", name)
	if _, err := os.Stat(refPath); os.IsNotExist(err) {
		return fmt.Errorf("Savefile %q does not exist", name)
	}

	// Resolve Savefile SHA
	branchSHA, err := gitobj.Ref_Resolve(repo, "refs/heads/"+name)
	if err != nil || branchSHA == nil {
		return fmt.Errorf("failed to resolve Savefile %q: %w", name, err)
	}

	// Dirty check is to warn user for unsaved changes
	// index, err := gitobj.Index_Read2(repo)
	// if err != nil {
	// 	return fmt.Errorf("failed to read Savelist: %w", err)
	// }
	// if isDirty(repo, *index) {
	// 	fmt.Println("WARNING : You have unsaved changes that will be lost if you switch.")
	// 	fmt.Print("    Continue anyway? (y/n): ")
	// 	reader := bufio.NewReader(os.Stdin)
	// 	input, _ := reader.ReadString('\n')
	// 	input = strings.TrimSpace(strings.ToLower(input))
	// 	if input != "y" {
	// 		fmt.Println("Switch cancelled.")
	// 		return nil
	// 	}
	// }

	// Get rules so we dont delete gitignored stuff
	// solution is not perfect but will do for now
	rules, err := gitcheckignore.ReadGitIgnore(repo)
	if err != nil {
		return fmt.Errorf("failed to read gitignore: %w", err)
	}

	// Backup worktree, will restore if error
	backupPath := filepath.Join(filepath.Dir(repo.WorkTree), "ezgit_backup")
	// fmt.Println("Backing up current worktree to", backupPath)
	err = copyDir(repo.WorkTree, backupPath, repo.GitDir, rules)
	if err != nil {
		return fmt.Errorf("failed to backup worktree: %w", err)
	}

	// Clear worktree (except .git)
	if err := clearWorktree(repo, rules); err != nil {
		// restore backup logic
		_ = copyDir(backupPath, repo.WorkTree, "", rules)
		return fmt.Errorf("failed to clear worktree: %w", err)
	}

	// Checkout Savefile tree into worktree
	commit, err := githashread.Object_Read(repo, *branchSHA)
	if err != nil {
		_ = copyDir(backupPath, repo.WorkTree, "", rules)
		return fmt.Errorf("failed to read commit: %w", err)
	}
	concreteCommit, ok := commit.(*gitobj.GitCommit)
	if !ok {
		_ = copyDir(backupPath, repo.WorkTree, "", rules)
		return fmt.Errorf("not a commit object")
	}
	concreteCommit.Deserialize()
	treeSHA := string(concreteCommit.KvlmDict.Dict["tree"])

	tree, err := githashread.Object_Read(repo, treeSHA)
	if err != nil {
		_ = copyDir(backupPath, repo.WorkTree, "", rules)
		return fmt.Errorf("failed to read tree: %w", err)
	}
	concreteTree, ok := tree.(gitobj.GitTree)
	if !ok {
		_ = copyDir(backupPath, repo.WorkTree, "", rules)
		return fmt.Errorf("not a tree object")
	}
	concreteTree.DeserializeData(tree.Deserialize())

	err = loading.TreeCheckout(repo, concreteTree, repo.WorkTree)
	if err != nil {
		// restore backup on failure
		_ = copyDir(backupPath, repo.WorkTree, "", rules)
		return fmt.Errorf("failed to checkout tree: %w", err)
	}

	// Update HEAD to point to Savefile ref
	headPath := gitpath.Repo_Path(repo, "HEAD")
	headContent := fmt.Sprintf("ref: refs/heads/%s\n", name)
	err = os.WriteFile(headPath, []byte(headContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	// Clear index in case files does not exist in switched Savefile
	emptyIndex := &gitobj.GitIndex{
		Version: 2,
		Entries: []gitobj.GitIndexEntry{},
	}
	if err := gitobj.Index_Write(repo, *emptyIndex); err != nil {
		return fmt.Errorf("failed to clear index: %w", err)
	}

	// Rebuild index via add --all
	_, err = add.Add(&repo, nil, add.Options{All: true})
	if err != nil {
		return fmt.Errorf("failed to rebuild index: %w", err)
	}

	// Delete backup on success
	os.RemoveAll(backupPath)

	fmt.Printf("Switched to Savefile: %q\n", name)
	return nil
}

// check if any file is not in the savelist. Checks Mtime
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

// IsDirty exposes the dirty check for the GUI to call before switching
func IsDirty(repo gitpath.GitRepository) (bool, error) {
	index, err := gitobj.Index_Read2(repo)
	if err != nil {
		return false, err
	}

	// check modified tracked files (mtime-based)
	if isDirty(repo, *index) {
		return true, nil
	}

	// check untracked / unstaged modified / deleted files in worktree
	worktreeStatus, err := gitCurrent.StatusIndexWorktree(repo, *index)
	if err != nil {
		return false, err
	}
	if len(worktreeStatus.Untracked) > 0 || len(worktreeStatus.Modified) > 0 || len(worktreeStatus.Deleted) > 0 {
		return true, nil
	}

	// check staged but not committed changes
	headStatus, err := gitCurrent.StatusHeadIndex(repo, *index)
	if err != nil {
		return false, err
	}
	if headStatus.HasChanges() {
		return true, nil
	}

	return false, nil
}

// delete all files/dirs in worktree except .git
func clearWorktree(repo gitpath.GitRepository, rules *gitcheckignore.GitIgnore) error {
	entries, err := os.ReadDir(repo.WorkTree)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}

		fullPath := filepath.Join(repo.WorkTree, entry.Name())
		rel, err := filepath.Rel(repo.WorkTree, fullPath)
		if err != nil {
			continue
		}

		// skip ignored files — they aren't tracked so don't touch them
		ignored, err := gitcheckignore.CheckIgnore(rules, rel)
		if err == nil && ignored {
			continue
		}

		if err := os.RemoveAll(fullPath); err != nil {
			return err
		}
	}
	return nil
}

// copyDir recursively copies src to dst, skipping skipPath
func copyDir(src, dst, skipPath string, rules *gitcheckignore.GitIgnore) error {
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

		// skip ignored files — no need to back them up
		if rules != nil {
			ignored, err := gitcheckignore.CheckIgnore(rules, rel)
			if err == nil && ignored {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
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

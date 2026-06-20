package alternateversions

import (
	"fmt"
	add "gocmd/testfiles/GitAddRemove"
	gitCurrent "gocmd/testfiles/GitCurrent"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitsave "gocmd/testfiles/GitSave"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
)

type MergeConflict struct {
	FilePath    string
	CurrentContent string // current branch version
	IncomingContent string // branch being merged in
}

type MergeResult struct {
	Conflicts []MergeConflict
	AutoMerged []string // files merged without conflict
}

// MergeBranch only detects what would change — it does NOT write anything
// to the working tree. Call ApplyAutoMergedFiles separately once the user
// confirms the merge (so Cancel can leave the working tree untouched).
func MergeBranch(repo gitpath.GitRepository, targetBranch string) (*MergeResult, error) {
	activeBranch, err := gitCurrent.Get_Active_Branch(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get active branch: %w", err)
	}
	if activeBranch == targetBranch {
		return nil, fmt.Errorf("cannot merge a branch into itself")
	}

	// Get current branch tree
	currentTreeSHA, err := gitCurrent.Get_Tree_SHA(repo, "HEAD")
	if err != nil {
		return nil, fmt.Errorf("failed to get current tree: %w", err)
	}
	currentFiles := map[string]string{}
	gitCurrent.TreeToMap(repo, currentTreeSHA, "", currentFiles)

	// Get target branch SHA and tree
	targetSHA, err := gitobj.Ref_Resolve(repo, "refs/heads/"+targetBranch)
	if err != nil || targetSHA == nil {
		return nil, fmt.Errorf("failed to resolve branch %q: %w", targetBranch, err)
	}

	targetObj, err := githashread.Object_Read(repo, *targetSHA)
	if err != nil {
		return nil, fmt.Errorf("failed to read target commit: %w", err)
	}
	targetCommit, ok := targetObj.(*gitobj.GitCommit)
	if !ok {
		return nil, fmt.Errorf("not a commit object")
	}
	targetCommit.Deserialize()
	targetTreeSHA := strings.TrimSpace(string(targetCommit.KvlmDict.Dict["tree"]))

	targetFiles := map[string]string{}
	gitCurrent.TreeToMap(repo, targetTreeSHA, "", targetFiles)

	result := &MergeResult{}

	// Check each file in target branch
	for filePath, targetSHAVal := range targetFiles {
		currentSHAVal, existsInCurrent := currentFiles[filePath]

		if !existsInCurrent {
			// File only in target — safe to add
			result.AutoMerged = append(result.AutoMerged, filePath)
			continue
		}

		if currentSHAVal == targetSHAVal {
			// Same content — no conflict
			continue
		}

		// Different content — conflict
		currentContent := readBlobContent(repo, currentSHAVal)
		targetContent := readBlobContent(repo, targetSHAVal)

		result.Conflicts = append(result.Conflicts, MergeConflict{
			FilePath:        filePath,
			CurrentContent:  currentContent,
			IncomingContent: targetContent,
		})
	}

	return result, nil
}

// ApplyAutoMergedFiles writes the auto-merged (non-conflicting) target
// branch files to the working tree. Call this only after the user has
// confirmed they want to proceed with the merge.
func ApplyAutoMergedFiles(repo gitpath.GitRepository, targetBranch string) error {
	currentTreeSHA, err := gitCurrent.Get_Tree_SHA(repo, "HEAD")
	if err != nil {
		return fmt.Errorf("failed to get current tree: %w", err)
	}
	currentFiles := map[string]string{}
	gitCurrent.TreeToMap(repo, currentTreeSHA, "", currentFiles)

	targetSHA, err := gitobj.Ref_Resolve(repo, "refs/heads/"+targetBranch)
	if err != nil || targetSHA == nil {
		return fmt.Errorf("failed to resolve branch %q: %w", targetBranch, err)
	}

	targetObj, err := githashread.Object_Read(repo, *targetSHA)
	if err != nil {
		return fmt.Errorf("failed to read target save: %w", err)
	}
	targetCommit, ok := targetObj.(*gitobj.GitCommit)
	if !ok {
		return fmt.Errorf("not a save object")
	}
	targetCommit.Deserialize()
	targetTreeSHA := strings.TrimSpace(string(targetCommit.KvlmDict.Dict["tree"]))

	targetFiles := map[string]string{}
	gitCurrent.TreeToMap(repo, targetTreeSHA, "", targetFiles)

	for filePath, targetSHAVal := range targetFiles {
		currentSHAVal, existsInCurrent := currentFiles[filePath]
		if existsInCurrent && currentSHAVal == targetSHAVal {
			continue
		}
		content := readBlobContent(repo, targetSHAVal)
		fullPath := filepath.Join(repo.WorkTree, filepath.FromSlash(filePath))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func readBlobContent(repo gitpath.GitRepository, sha string) string {
	obj, err := githashread.Object_Read(repo, sha)
	if err != nil {
		return ""
	}
	return string(obj.Deserialize())
}

func ApplyConflictResolution(repo gitpath.GitRepository, filePath string, content string) error {
	fullPath := filepath.Join(repo.WorkTree, filepath.FromSlash(filePath))
	return os.WriteFile(fullPath, []byte(content), 0644)
}

func CompleteMerge(repo gitpath.GitRepository, targetBranch string, message string, window fyne.Window) error {
	// Get current HEAD SHA (first parent)
	currentSHA, err := gitobj.Ref_Resolve(repo, "HEAD")
	if err != nil || currentSHA == nil {
		return fmt.Errorf("failed to resolve HEAD: %w", err)
	}

	// Get target branch SHA (second parent)
	targetSHA, err := gitobj.Ref_Resolve(repo, "refs/heads/"+targetBranch)
	if err != nil || targetSHA == nil {
		return fmt.Errorf("failed to resolve target branch: %w", err)
	}

	// Read current index and build tree
	index, err := gitobj.Index_Read2(repo)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	// Re-add all worktree files to index to capture conflict resolutions
	_, err = add.Add(&repo, nil, add.Options{All: true})
	if err != nil {
		return fmt.Errorf("failed to stage merged files: %w", err)
	}

	// Re-read index after add
	index, err = gitobj.Index_Read2(repo)
	if err != nil {
		return fmt.Errorf("failed to re-read index: %w", err)
	}

	// Build tree from index
	treeSHA, err := gitobj.TreeFromIndex(repo, *index)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	// Two parents: current HEAD + target branch
	parents := []string{*currentSHA, *targetSHA}

	// Get author from config
	userConfig, err := gitpath.Load()
	if err != nil {
		return fmt.Errorf("no user config found, please save first to set up your name/email: %w", err)
	}

	commitSHA, err := gitsave.Version_Create(repo, treeSHA, parents, userConfig.Format(), time.Now(), message)
	if err != nil {
		return fmt.Errorf("failed to create merge commit: %w", err)
	}

	// Update current branch ref
	_, err = gitsave.Update_Branch_Ref(repo, commitSHA)
	if err != nil {
		return fmt.Errorf("failed to update branch ref: %w", err)
	}

	// Refresh index
	err = gitsave.RefreshIndex(repo, index)
	if err != nil {
		return fmt.Errorf("failed to refresh index: %w", err)
	}

	return nil
}
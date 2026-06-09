package gitsave

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	setconfig "gocmd/internal/cli/set-config"
	add "gocmd/testfiles/GitAddRemove"
	gitCurrent "gocmd/testfiles/GitCurrent"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
)

type conflict struct {
	path       string
	currentSHA string
	targetSHA  string
}

// getCommitTreeSHA resolves a commit SHA to its tree SHA
func getCommitTreeSHA(repo gitpath.GitRepository, commitSHA string) (string, error) {
	obj, err := githashread.Object_Read(repo, commitSHA)
	if err != nil {
		return "", fmt.Errorf("failed to read commit %s: %w", commitSHA, err)
	}
	commit, ok := obj.(*gitobj.GitCommit)
	if !ok {
		return "", fmt.Errorf("object %s is not a commit", commitSHA)
	}
	commit.Deserialize()
	return strings.TrimSpace(string(commit.KvlmDict.Dict["tree"])), nil
}

// buildFileMap builds a flat path→sha map from a commit SHA
func buildFileMap(repo gitpath.GitRepository, commitSHA string) (map[string]string, error) {
	treeSHA, err := getCommitTreeSHA(repo, commitSHA)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	if err := gitCurrent.TreeToMap(repo, treeSHA, "", result); err != nil {
		return nil, err
	}
	return result, nil
}

// checkDirty returns true if index has uncommitted changes
func checkDirty(repo gitpath.GitRepository, index gitobj.GitIndex) bool {
	for _, entry := range index.Entries {
		fullPath := filepath.Join(repo.WorkTree, entry.Name)
		info, err := os.Stat(fullPath)
		if err != nil {
			return true
		}
		if uint32(info.ModTime().Unix()) != entry.MtimeSec {
			return true
		}
	}
	return false
}

// viewFile prints file contents with a border
func viewFile(label, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("  could not read file: %v\n", err)
		return
	}
	border := strings.Repeat("─", 60)
	fmt.Printf("┌%s┐\n", border)
	fmt.Printf("│ %-58s │\n", label)
	fmt.Printf("├%s┤\n", border)
	for _, line := range strings.Split(string(data), "\n") {
		for _, wrapped := range wrapText(line, 58) {
			fmt.Printf("│ %-58s │\n", wrapped)
		}
	}
	fmt.Printf("└%s┘\n", border)
}

// resolveConflict handles the interactive prompt for one conflicting file
// returns the chosen resolution: "current", "target", or "both"
func resolveConflict(
	repo gitpath.GitRepository,
	filePath string,
	currentSHA string,
	targetSHA string,
	currentBranch string,
	targetBranch string,
) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fullPath := filepath.Join(repo.WorkTree, filePath)

	// write both versions to temp files for viewing
	currentTemp := fullPath + ".__current__"
	targetTemp := fullPath + ".__target__"

	// read blob content for current
	currentObj, err := githashread.Object_Read(repo, currentSHA)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(currentTemp, currentObj.Deserialize(), 0644); err != nil {
		return "", err
	}
	defer os.Remove(currentTemp)

	// read blob content for target
	targetObj, err := githashread.Object_Read(repo, targetSHA)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(targetTemp, targetObj.Deserialize(), 0644); err != nil {
		return "", err
	}
	defer os.Remove(targetTemp)

	border := strings.Repeat("─", 60)
	for {
		fmt.Printf("\n┌%s┐\n", border)
		fmt.Printf("│ CONFLICT: %-48s │\n", filePath)
		fmt.Printf("├%s┤\n", border)
		fmt.Printf("│  [1]  Keep from %-42s │\n", currentBranch+" (current)")
		fmt.Printf("│  [2]  Keep from %-42s │\n", targetBranch+" (target)")
		fmt.Printf("│  [3]  Write  both as separate files%-24s│\n", "") // Fixed padding
		fmt.Printf("│  [v1] View %-47s │\n", currentBranch+" version")
		fmt.Printf("│  [v2] View %-47s │\n", targetBranch+" version")
		fmt.Printf("│  [b]  View both%-44s│\n", "") // Fixed padding
		fmt.Printf("└%s┘\n", border)
		fmt.Print("  choice > ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "1":
			return "current", nil
		case "2":
			return "target", nil
		case "3":
			return "both", nil
		case "v1":
			viewFile(currentBranch+" version of "+filePath, currentTemp)
		case "v2":
			viewFile(targetBranch+" version of "+filePath, targetTemp)
		case "b":
			viewFile(currentBranch+" version of "+filePath, currentTemp)
			viewFile(targetBranch+" version of "+filePath, targetTemp)
		default:
			fmt.Println("  invalid choice, try again.")
		}
	}
}

// applyResolution writes the resolved file to the worktree
func applyResolution(
	repo gitpath.GitRepository,
	filePath string,
	choice string,
	currentSHA string,
	targetSHA string,
	currentBranch string,
	targetBranch string,
) error {
	fullPath := filepath.Join(repo.WorkTree, filePath)

	readBlob := func(sha string) ([]byte, error) {
		obj, err := githashread.Object_Read(repo, sha)
		if err != nil {
			return nil, err
		}
		return obj.Deserialize(), nil
	}

	ext := filepath.Ext(filePath)
	base := strings.TrimSuffix(filePath, ext)

	switch choice {
	case "current":
		data, err := readBlob(currentSHA)
		if err != nil {
			return err
		}
		return os.WriteFile(fullPath, data, 0644)

	case "target":
		data, err := readBlob(targetSHA)
		if err != nil {
			return err
		}
		return os.WriteFile(fullPath, data, 0644)

	case "both":
		// write current branch version
		currentData, err := readBlob(currentSHA)
		if err != nil {
			return err
		}
		currentPath := filepath.Join(repo.WorkTree, base+"_"+currentBranch+ext)
		if err := os.WriteFile(currentPath, currentData, 0644); err != nil {
			return err
		}
		// write target branch version
		targetData, err := readBlob(targetSHA)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(repo.WorkTree, base+"_"+targetBranch+ext)
		if err := os.WriteFile(targetPath, targetData, 0644); err != nil {
			return err
		}
		// remove original conflicted file
		return os.Remove(fullPath)
	}
	return nil
}

func Combine(repo gitpath.GitRepository, targetBranch string) error {
	reader := bufio.NewReader(os.Stdin)

	// get current branch from HEAD
	currentBranch, err := gitCurrent.Get_Active_Branch(repo)
	if err != nil || currentBranch == "" {
		return fmt.Errorf("could not determine current branch")
	}
	//resolve head to get sha
	currentCommitSHA, err := gitobj.Ref_Resolve(repo, "HEAD")
	if err != nil || currentCommitSHA == nil {
		return fmt.Errorf("could not resolve HEAD")
	}

	// check target branch exists
	targetRef := "refs/heads/" + targetBranch
	targetCommitSHA, err := gitobj.Ref_Resolve(repo, targetRef)
	if err != nil || targetCommitSHA == nil {
		return fmt.Errorf("alternate version %q not found", targetBranch)
	}

	// ast forward check if HEAD is ancestor of target
	if *currentCommitSHA == *targetCommitSHA {
		fmt.Println("Already up to date.")
		return nil
	}

	// dirty index check
	index, err := gitobj.Index_Read2(repo)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}
	if checkDirty(repo, *index) {
		fmt.Println(" You have unsaved changes.")
		fmt.Print("   Continue anyway? (y/n): ")
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(input)) != "y" {
			fmt.Println("Combine cancelled.")
			return nil
		}
	}

	// build file maps from both commits
	fmt.Println("Comparing versions...")
	currentFiles, err := buildFileMap(repo, *currentCommitSHA)
	if err != nil {
		return fmt.Errorf("failed to read current branch files: %w", err)
	}
	targetFiles, err := buildFileMap(repo, *targetCommitSHA)
	if err != nil {
		return fmt.Errorf("failed to read target branch files: %w", err)
	}

	// categorize files
	var conflicts []conflict

	// files only in target, auto add
	for path, sha := range targetFiles {
		if _, exists := currentFiles[path]; !exists {
			// fmt.Printf("  adding:   %s\n", path)
			obj, err := githashread.Object_Read(repo, sha)
			if err != nil {
				return err
			}
			fullPath := filepath.Join(repo.WorkTree, path)
			if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(fullPath, obj.Deserialize(), 0644); err != nil {
				return err
			}
		}
	}

	// files only in current — keep as is (target deleted, keep current)
	for path := range currentFiles {
		if _, exists := targetFiles[path]; !exists {
			fmt.Printf("  keeping:  %s (not in target)\n", path)
		}
	}

	// files in both — check SHA
	for path, currentSHA := range currentFiles {
		if targetSHA, exists := targetFiles[path]; exists {
			if currentSHA != targetSHA {
				conflicts = append(conflicts, conflict{path, currentSHA, targetSHA})
			}
			// same SHA — no action needed
		}
	}

	// resolve conflicts interactively
	if len(conflicts) == 0 {
		fmt.Println("No conflicts found.")
	} else {
		fmt.Printf("\n%d conflict(s) to resolve:\n", len(conflicts))
		for _, c := range conflicts {
			choice, err := resolveConflict(repo, c.path, c.currentSHA, c.targetSHA, currentBranch, targetBranch)
			if err != nil {
				return fmt.Errorf("failed to resolve conflict for %s: %w", c.path, err)
			}
			if err := applyResolution(repo, c.path, choice, c.currentSHA, c.targetSHA, currentBranch, targetBranch); err != nil {
				return fmt.Errorf("failed to apply resolution for %s: %w", c.path, err)
			}
		}
	}

	// 9. rebuild index
	emptyIndex := &gitobj.GitIndex{Version: 2, Entries: []gitobj.GitIndexEntry{}}
	if err := gitobj.Index_Write(repo, *emptyIndex); err != nil {
		return fmt.Errorf("failed to clear index: %w", err)
	}
	if _, err := add.Add(&repo, nil, add.Options{All: true}); err != nil {
		return fmt.Errorf("failed to rebuild index: %w", err)
	}

	// 10. auto commit merge
	index, err = gitobj.Index_Read2(repo)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}
	treeSHA, err := gitobj.TreeFromIndex(repo, *index)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	cfg, err := setconfig.GetOrPromptConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	mergeMessage := fmt.Sprintf("Merge alternate Savefile '%s' into '%s'", targetBranch, currentBranch)
	commitSHA, err := Version_Create(
		repo,
		treeSHA,
		[]string{*currentCommitSHA, *targetCommitSHA}, // two parents
		cfg.Format(),
		time.Now(),
		mergeMessage,
	)
	if err != nil {
		return fmt.Errorf("failed to create combined version: %w", err)
	}

	// 11. update branch ref
	_, err = Update_Branch_Ref(repo, commitSHA)
	if err != nil {
		return fmt.Errorf("failed to update savefile: %w", err)
	}

	fmt.Printf("\nCombined '%s' into '%s' → %s\n", targetBranch, currentBranch, commitSHA[:7])
	return nil
}

func wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}
	var lines []string
	for len(text) > width {
		lines = append(lines, text[:width])
		text = text[width:]
	}
	if text != "" {
		lines = append(lines, text)
	}
	return lines
}

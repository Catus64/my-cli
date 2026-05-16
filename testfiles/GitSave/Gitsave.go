package gitsave

import (
	"bufio"
	"fmt"
	gitCurrent "gocmd/testfiles/GitCurrent"
	gitcur "gocmd/testfiles/GitCurrent"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// format timezone from timestamp
func get_timezone(timestamp time.Time) string {
	_, offset := timestamp.Zone() // offset in seconds
	hours := offset / 3600
	minutes := (offset % 3600) / 60
	sign := "+"
	if offset < 0 {
		sign = "-"
		hours = -hours
		minutes = -minutes
	}
	time_zone := fmt.Sprintf("%s%02d%02d", sign, hours, minutes)

	return time_zone
}

func Version_Create(
	repo gitpath.GitRepository,
	treeSHA string,
	parentSHA []string, //for multiple parents
	author string,
	timestamp time.Time,
	message string,
) (string, error) {

	time_zone := get_timezone(timestamp)

	// Build author string: "Name email unix timestamp timezone"
	author_line := fmt.Sprintf("%s %d %s", author, timestamp.Unix(), time_zone)

	message = strings.TrimSpace(message) + "\n"

	// Git expects: tree, parent (optional), author, committer, then message
	dict := make(map[string][]byte)
	dict["tree"] = []byte(treeSHA)
	if len(parentSHA) > 0 {
		dict["parent"] = []byte(strings.Join(parentSHA, "\x00"))
	}
	dict["author"] = []byte(author_line)
	dict["committer"] = []byte(author_line)
	dict["data"] = []byte(message)

	commit := gitobj.MakeGitCommit(dict)

	//write commit object to /objects
	sha, err := githashread.Object_Write(commit, &repo)
	if err != nil {
		return "", err
	}

	sha_str := fmt.Sprintf("%x", sha)

	return sha_str, nil
}

func Update_Branch_Ref(repo gitpath.GitRepository, verSHA string) (string, error) {
	branchName, err := gitcur.Get_Active_Branch(repo)
	if err != nil {
		return "", err
	}
	if branchName == "" {
		return "", fmt.Errorf("detached HEAD, ezgit does not support detached HEAD: type `git checkout <branch>` to switch to a branch")
	}

	// Update the branch ref to point to the new commit SHA
	ref_path := gitpath.Repo_Path(repo, "refs", "heads", branchName)

	// in case file does not exist
	if err := os.MkdirAll(filepath.Dir(ref_path), 0755); err != nil {
		return "", fmt.Errorf("failed to create ref directory: %w", err)
	}

	err = os.WriteFile(ref_path, []byte(verSHA+"\n"), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to update Savefile ref :%w", err)
	}

	logger.L().Info("Savefile reference updated", "Savefile", branchName, "new_sha", verSHA)

	return branchName, nil
}

func RefreshIndex(repo gitpath.GitRepository, index *gitobj.GitIndex) error {
	for i, entry := range index.Entries {
		fullPath := filepath.Join(repo.WorkTree, entry.Name)
		info, err := os.Stat(fullPath)
		if err != nil {
			continue // file might not exist
		}
		// update stat data to current
		index.Entries[i].MtimeSec = uint32(info.ModTime().Unix())
		index.Entries[i].MtimeNano = uint32(info.ModTime().Nanosecond())
		index.Entries[i].FSize = uint32(info.Size())

		index.Entries[i].MtimeSec = uint32(info.ModTime().Unix())
		index.Entries[i].MtimeNano = uint32(info.ModTime().Nanosecond())
		index.Entries[i].FSize = uint32(info.Size())
		updateStatTimes(index, i, info)
	}
	// fmt.Println("index refresh")
	return gitobj.Index_Write(repo, *index)
}

func CheckCommitReady(repo gitpath.GitRepository, index gitobj.GitIndex) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	// Check if nothing changed
	headResult, err := gitCurrent.StatusHeadIndex(repo, index)
	if err != nil {
		return false, err
	}
	if !headResult.HasChanges() {
		fmt.Println("Nothing to commit your savelist is the same as the version's content.")
		return false, nil
	}

	// Check for unstaged changes
	worktreeResult, err := gitCurrent.StatusIndexWorktree(repo, index)
	if err != nil {
		return false, err
	}
	if worktreeResult.HasUnstaged() {
		fmt.Println("!!WARNING!!  These changes are NOT included in this save:")
		for _, f := range worktreeResult.Modified {
			fmt.Println("  modified (not added):", f)
		}
		for _, f := range worktreeResult.Deleted {
			fmt.Println("  deleted  (not added):", f)
		}
		fmt.Println("\n!!WARNING!!  Progress on these files can be lost if you switch branches.")
		fmt.Print("  Continue saving without them? (y/n): ")
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(input)) != "y" {
			fmt.Println("Save cancelled.")
			return false, nil
		}
	}

	return true, nil
}

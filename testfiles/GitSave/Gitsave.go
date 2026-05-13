package gitsave

import (
	"fmt"
	gitcur "gocmd/testfiles/GitCurrent"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func Version_Create(
	repo gitpath.GitRepository,
	treeSHA string,
	parentSHA []string, //for multiple parents
	author string,
	timestamp time.Time,
	message string,
) (string, error) {

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

	sha, err := githashread.Object_Write(commit, &repo)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha), nil
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
	err = os.WriteFile(ref_path, []byte(verSHA+"\n"), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to update branch ref :%w", err)
	}

	logger.L().Info("Branch reference updated", "branch", branchName, "new_sha", verSHA)

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

		if sysStat, ok := info.Sys().(*syscall.Stat_t); ok {
			index.Entries[i].CtimeSec = uint32(sysStat.Ctim.Sec)
			index.Entries[i].CtimeNano = uint32(sysStat.Ctim.Nsec)
		}
	}
	fmt.Println("index refresh")
	return gitobj.Index_Write(repo, *index)
}

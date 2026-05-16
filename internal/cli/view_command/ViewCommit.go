package viewcommand

import (
	"bufio"
	"fmt"
	gitcurrent "gocmd/testfiles/GitCurrent"
	githashread "gocmd/testfiles/GitHashRead"
	gitlog "gocmd/testfiles/GitLog"
	gitobj "gocmd/testfiles/GitObject"
	gitsave "gocmd/testfiles/GitSave"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"
	"os"
	"strings"
)

func viewCommit(repo gitpath.GitRepository, sha string) error {
	obj, err := githashread.Object_Read(repo, sha)
	if err != nil {
		return err
	}
	commit, ok := obj.(*gitobj.GitCommit)
	if !ok {
		return fmt.Errorf("not a commit object")
	}
	commit.Deserialize()

	shortSHA := sha
	if len(sha) > 7 {
		shortSHA = sha[:7]
	}

	parents := gitobj.GetKvlmValues(commit.KvlmDict.Dict, "parent")
	parentDisplay := ""
	if len(parents) > 0 {
		parentDisplay = parents[0][:7]
		if len(parents) > 1 {
			parentDisplay += fmt.Sprintf(" (+%d more)", len(parents)-1)
		}
	}

	date, author := gitlog.Format_Date_Author(string(commit.KvlmDict.Dict["author"]))
	treeSHA := strings.TrimSpace(string(commit.KvlmDict.Dict["tree"]))

	prettyprint.PrintCommit(
		shortSHA,
		author,
		date,
		treeSHA[:7],
		parentDisplay,
		string(commit.KvlmDict.Dict["data"]),
	)

	// prompt to view tree
	fmt.Print("\npress v to view tree, any other key to exit: ")
	branch, _ := gitcurrent.Get_Active_Branch(repo)
	header := ""
	if branch != "" {
		entry, err := gitsave.ReadVersionRef(repo, branch, sha)
		if err == nil {
			header = fmt.Sprintf("v%d · %s | %s", entry.Number, entry.Name, branch)
		}
	}

	prettyprint.PrintCommit(shortSHA, author, date, treeSHA[:7], parentDisplay, string(commit.KvlmDict.Dict["data"]))

	fmt.Print("\npress v to view tree, any other key to exit: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadByte()
	if input == 'v' {
		return viewTree(repo, treeSHA, header) // pass header
	}
	return nil
}

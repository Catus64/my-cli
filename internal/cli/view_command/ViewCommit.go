package viewcommand

import (
	"bufio"
	"fmt"
	gitcurrent "gocmd/testfiles/GitCurrent"
	githashread "gocmd/testfiles/GitHashRead"
	gitlog "gocmd/testfiles/GitLog"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	prettyprint "gocmd/testfiles/PrettyPrint"
	"os"
	"strings"
)

func viewCommit(repo gitpath.GitRepository, sha string) error {
	logger.L().Debug("View commit running")

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

	// resolve best label — tag > version > nothing
	branch, _ := gitcurrent.Get_Active_Branch(repo)
	header := prettyprint.ResolveCommitLabel(repo, sha, branch)
	// logger.L().Debug("Checking Header Existence", "Header", header)
	// print header line if we have one
	if header != "" {
		fmt.Println(header)
	}

	prettyprint.PrintCommit(
		shortSHA,
		author,
		date,
		treeSHA[:7],
		parentDisplay,
		header,
		string(commit.KvlmDict.Dict["data"]),
	)

	fmt.Print("\npress v to view tree, any other key to exit: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadByte()
	if input == 'v' {
		return viewTree(repo, treeSHA, header)
	}
	return nil
}

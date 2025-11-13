package gitlog

import (
	"bytes"
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
	"strconv"
	"strings"
	"time"
)

func Read_Master(repo gitpath.GitRepository) string {
	head := gitpath.Repo_Path(repo, "refs", "heads", "master")
	data, err := os.ReadFile(head)
	if err != nil {
		panic(err)
	}
	data = bytes.ReplaceAll(data, []byte(" "), []byte(""))
	data = bytes.ReplaceAll(data, []byte("\n"), []byte(""))

	return string(data)
}

func Format_Date_Author(text string) (string, string) {
	temp_string := strings.SplitN(text, " ", -1)
	temp_time := temp_string[2]
	num, err := strconv.ParseInt(temp_time, 10, 64)
	if err != nil {
		panic(err)
	}
	loc, _ := time.LoadLocation("Local")
	t := time.Unix(num, 0).In(loc)

	author := temp_string[0] + temp_string[1]

	date := t.Format(time.RFC1123Z)

	return date, author
}

func Log_One(kvlm gitobj.KvlmDict) {

}

func Recurse_Log(repo *gitpath.GitRepository, commit gitobj.GitCommit, sha string) {
	//log one obj
	fmt.Println("commit ", sha)
	//fmt.Println("Author:", string(commit.Dict["author"]))
	date, author := Format_Date_Author(string(commit.Dict["author"]))

	fmt.Println("Date: ", date)
	fmt.Println("Author: ", author)

	//fmt.Println("\n", string(commit.Dict["data"]))

	//parent = current pbj = end func
	if string(commit.Dict["parent"]) == sha {
		println("end")
		return
	}

	//recurse
	Parent_Obj := githashread.Object_Read(*repo, string(commit.Dict["parent"]))
	Concrete_Parent_Commit, ok := Parent_Obj.(*gitobj.GitCommit)
	if !ok {
		panic("not a commit object")
	}
	Concrete_Parent_Commit.Deserialize()
	Recurse_Log(repo, *Concrete_Parent_Commit, string(commit.Dict["parent"]))
}

func Log(repo gitpath.GitRepository) error {
	head := Read_Master(repo)
	Commit_Object := githashread.Object_Read(repo, head)
	Commit_Object.Deserialize()
	Concrete_Commit, ok := Commit_Object.(*gitobj.GitCommit)
	if !ok {
		panic("not a commit object")
	}
	//fmt.Println(string(Concrete_Commit.Dict["data"]))
	Recurse_Log(&repo, *Concrete_Commit, head)
	return nil
}

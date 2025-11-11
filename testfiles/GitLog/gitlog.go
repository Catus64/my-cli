package gitlog

import (
	"fmt"
	read "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
)

func Format_date_time(time string) string {
	return " "
}

func Log_One(kvlm gitobj.KvlmDict) {

}

func Recurse_Git(commit gitobj.GitCommit) {
	//log one obj

	//parent = current pbj = end func

	//recurse
}

func Log(repo gitpath.GitRepository, path string) error {
	Commit_Object := read.Object_Read(repo, path)
	Commit_Object.Deserialize()
	Concrete_Commit, ok := Commit_Object.(*gitobj.GitCommit)
	if !ok {
		panic(fmt.Sprintf("%v is not a commit object", path))
	}
	fmt.Println(Concrete_Commit.KvlmDict.Dict["parent"])
	return nil
}

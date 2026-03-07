package GitObjLib

import (
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
	"time"
)

func Ref_Resolve(repo gitpath.GitRepository, ref string) *string {
	path := gitpath.Repo_Path(repo, ref)
	//fmt.Println(path)
	_, err := os.Stat(path)
	if err != nil {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	data = data[:len(data)-1] // remove newline
	//fmt.Println(string(data))
	if string(data[0:4]) == string("ref:") {
		return Ref_Resolve(repo, string(data[5:]))
	}

	ret := string(data)

	return &ret
}

type GitTag struct {
	GitObjectData
	Tag_Name    string
	Tag_Email   string
	Tag_Date    time.Time
	Tag_Message string
	Tag_Object  string
	Tag_Type    string
}

func Create_Ref(repo gitpath.GitRepository, name string, sha string) error {
	path := gitpath.Repo_Path(repo, "refs", "tags", name)
	err := os.WriteFile(path, []byte(sha+"\n"), 0644)
	if err != nil {
		return err
	}
	return nil
}

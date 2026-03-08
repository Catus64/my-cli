package GitObjLib

import (
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
	"time"
)

// take a ref eg, ".git/HEAD" and resolve it to the sha
// it points to, if it is a ref to another ref then resolve that one instead

func Ref_Resolve(repo gitpath.GitRepository, ref string) (*string, error) {
	path := gitpath.Repo_Path(repo, ref)

	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data = data[:len(data)-1] // remove newline from end of file
	if string(data[0:4]) == string("ref:") {
		return Ref_Resolve(repo, string(data[5:]))
	}

	ret := string(data)
	return &ret, nil
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

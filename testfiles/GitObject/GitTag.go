package GitObjLib

import (
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
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
	KvlmDict
	format []byte
}

func (t *GitTag) Get_Format() string {
	return "tag"
}

func (t *GitTag) Deserialize() []byte {
	if t.data == nil {
		return Kvlm_Serialize(t.KvlmDict)
	}
	temp := make(map[string][]byte)
	kvlm := KvlmDict{Dict: temp}
	t.KvlmDict = Kvlm_Parse(t.data, 0, kvlm)
	return t.data
}

func (t *GitTag) Serialize() *[]byte {
	ret := Kvlm_Serialize(t.KvlmDict)
	return &ret
}

func Create_Ref(repo gitpath.GitRepository, name string, sha string) error {
	path := gitpath.Repo_Path(repo, "refs", "tags", name)
	err := os.WriteFile(path, []byte(sha+"\n"), 0644)
	if err != nil {
		return err
	}
	return nil
}

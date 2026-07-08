package GitObjLib

import (
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
	"strings"
)

// take a ref eg, ".git/HEAD" and resolve it to the sha
// it points to, if it is a ref to another ref then resolve that one instead

func Ref_Resolve(repo gitpath.GitRepository, ref string) (*string, error) {
	path := gitpath.Repo_Path(repo, ref)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// fall back to packed-refs
			if sha, ok := lookupPackedRef(repo, ref); ok {
				return &sha, nil
			}
		}
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data = data[:len(data)-1]
	if string(data[0:4]) == "ref:" {
		return Ref_Resolve(repo, strings.TrimSpace(string(data[5:])))
	}
	ret := string(data)
	return &ret, nil
}

func lookupPackedRef(repo gitpath.GitRepository, ref string) (string, bool) {
	packedPath := gitpath.Repo_Path(repo, "packed-refs") // .git/packed-refs, not under refs/
	data, err := os.ReadFile(packedPath)
	if err != nil {
		return "", false
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "^") {
			continue // comments and peeled-tag lines
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		sha, name := parts[0], parts[1]
		if name == ref {
			return sha, true
		}
	}
	return "", false
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
	if t.Data == nil {
		return Kvlm_Serialize(t.KvlmDict)
	}
	temp := make(map[string][]byte)
	kvlm := KvlmDict{Dict: temp}
	t.KvlmDict = Kvlm_Parse(t.Data, 0, kvlm)
	return t.Data
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

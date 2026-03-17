package GitObjLib

import (
	"bytes"
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type GitObjectData struct {
	data []byte
}

type GitObject interface {
	Serialize() *[]byte
	Deserialize() []byte
	Get_Format() string
}

type GitBlob struct {
	GitObjectData
	format []byte
}

func (blob GitBlob) Serialize() *[]byte {
	return nil
}

func (blob GitBlob) Deserialize() []byte {
	return blob.data
}

func (blob GitBlob) Get_Format() string {
	return string(blob.format)
}

// factory class to make git object from binary blob
func MakeGitObj(Byte_data []byte) GitObject {
	parts := bytes.SplitN(Byte_data, []byte{0}, -1)
	header := bytes.SplitN(parts[0], []byte{32}, -1)

	temp_fmt := header[0]

	var obj GitObject = nil

	switch string(temp_fmt) {
	case "blob":
		// fmt.Println("returning blob")
		obj = GitBlob{GitObjectData: GitObjectData{parts[1]}, format: []byte("blob")}
	case "commit":
		//fmt.Println("returning commit object")
		obj = &GitCommit{GitObjectData: GitObjectData{parts[1]}, format: []byte("commit")}
	case "tag":
		fmt.Println("returning tag")
	case "tree":
		//fmt.Println("returning tree")
		x := bytes.IndexByte(Byte_data[0:], 0x00)
		GitTreeData := Byte_data[x+1:]
		obj = GitTree{GitObjectData: GitObjectData{GitTreeData}, format: []byte("tree")}
	case "ref":
		fmt.Println("returning ref")
	default:
	}

	return obj
}

// return an object from its data and specified format (mostly for writing)
func MakeGitObjWithFormat(Byte_data []byte, Obj_format string) GitObject {
	var obj GitObject = nil

	switch string(Obj_format) {
	case "blob":
		// fmt.Println("returning blob")
		obj = GitBlob{GitObjectData: GitObjectData{Byte_data}, format: []byte("blob")}
	case "commit":
		fmt.Println("returning commit")
	case "tag":
		fmt.Println("returning tag")
	case "tree":
		fmt.Println("returning tree")
	case "ref":
		fmt.Println("returning ref")
	default:
	}

	return obj
}

func Object_Resolve(repo gitpath.GitRepository, name string) ([]string, error) {

	var candidates []string

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, nil
	}

	// HEAD
	if name == "HEAD" {
		sha, err := Ref_Resolve(repo, "HEAD")
		if err != nil {
			return nil, err
		}
		return []string{*sha}, nil
	}

	// compile regex to check if name is a hash
	hashRegex := regexp.MustCompile(`^[0-9A-Fa-f]{4,40}$`)

	// check whether name is a hash
	if hashRegex.MatchString(name) {

		name = strings.ToLower(name)
		prefix := name[:2]

		path := filepath.Join(repo.GitDir, "objects", prefix)

		files, err := os.ReadDir(path)
		if err == nil {

			rem := name[2:]

			for _, f := range files {
				if strings.HasPrefix(f.Name(), rem) {
					candidates = append(candidates, prefix+f.Name())
				}
			}
		}
	}

	//if input is not a hash then check refs folder

	// refs/tags
	if sha, _ := Ref_Resolve(repo, "refs/tags/"+name); sha != nil {
		candidates = append(candidates, *sha)
	}

	// refs/heads
	if sha, _ := Ref_Resolve(repo, "refs/heads/"+name); sha != nil {
		candidates = append(candidates, *sha)
	}

	// refs/remotes
	if sha, _ := Ref_Resolve(repo, "refs/remotes/"+name); sha != nil {
		candidates = append(candidates, *sha)
	}

	return candidates, nil
}

package main

import (
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"
)

type GitObjectData struct {
	data []byte
}

type GitObject interface {
	Serialize() []byte
	Deserialize() []string
}

type GitBlob struct {
	obj GitObjectData
	fmt []byte
}

func (blob GitBlob) serialize() {
}

func (blob GitBlob) deserialize() {

}

func Object_Read(repo gitpath.GitRepository, sha string) GitObject {
	first := sha[:2]
	rest := sha[2:]

	fmt.Println(first)
	fmt.Println(rest)

	obj_path := gitpath.Repo_Path(repo, first, rest)
	fmt.Println(obj_path)
	return nil
}

func main() {
	required := true
	repo, _ := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)

	fmt.Println(repo.GitDir)
	Object_Read(*repo, "0e5108686cccc75435e3515a916fb7ceaf9b248d")
}

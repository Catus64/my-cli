package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"io"
	"os"
)

type GitObjectData struct {
	data []byte
}

type GitObject interface {
	Serialize() *[]byte
	Deserialize() *[]string
	Get_Format() string
}

type GitBlob struct {
	obj GitObjectData
	fmt []byte
}

func (blob GitBlob) Serialize() *[]byte {
	return nil
}

func (blob GitBlob) Deserialize() *[]string {
	return nil
}

func (blob GitBlob) Get_Format() string {
	return string(blob.fmt)
}

func Read_Blob(path string) *[]byte {

	//check whether file exist
	_, err := os.Stat(path)
	if err != nil {
		panic(fmt.Sprintf("%v does not exist/error", path))
	}

	fmt.Println("File exist")

	//read file
	filedata, err := os.ReadFile(path)
	if err != nil {
		panic("File cannot be read")
	}

	//decompressing via zlib
	r, err := zlib.NewReader(bytes.NewReader(filedata))
	if err != nil {
		panic("not zlib format")
	}
	defer r.Close()

	//reading decompressed data
	decompressed, err := io.ReadAll(r)
	if err != nil {
		panic("read fail")
	}

	return &decompressed
}

func MakeGitObj(data []byte) GitObject {
	parts := bytes.SplitN(data, []byte{0}, -1)
	header := bytes.SplitN(parts[0], []byte{32}, -1)

	temp_fmt := header[0]

	switch string(temp_fmt) {
	case "blob":
		fmt.Println("returning blob")
	case "commit":
		fmt.Println("returning commit")
	case "tag":
		fmt.Println("returning tag")
	case "tree":
		fmt.Println("returning tree")
	case "ref":
		fmt.Println("returning ref")
	}

	return nil
}

func Object_Read(repo gitpath.GitRepository, sha string) GitObject {

	//splitting strings for path check
	first := sha[:2]
	rest := sha[2:]

	//getting obj_path
	obj_path := gitpath.Repo_Path(repo, "objects", first, rest)
	fmt.Println(obj_path)

	decompressed := Read_Blob(obj_path)

	o := MakeGitObj(*decompressed)
	fmt.Println(o)

	/*
		fmt.Println("decomp:\n", ((*decompressed)[4]))

		parts := bytes.SplitN(*decompressed, []byte{0}, -1)
		headers := bytes.SplitN(parts[0], []byte{32}, -1)

		obj := GitBlob{
			obj: GitObjectData{data: *decompressed},
			fmt: headers[0],
		}
	*/

	return nil
}

func main() {
	required := true
	repo, _ := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)

	fmt.Println(repo.GitDir)

	//66726975a39cdb5babb3a02e8761ecd37e1c7c49

	obj := Object_Read(*repo, "66726975a39cdb5babb3a02e8761ecd37e1c7c49")
	fmt.Println(obj)
	//fmt.Printf("format: %v\n", obj.Get_Format())
}

package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"io"
	"os"
	"strconv"
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

// reading blob data to print them out
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

// factory class to make git object
func MakeGitObj(Byte_data []byte) GitObject {
	parts := bytes.SplitN(Byte_data, []byte{0}, -1)
	header := bytes.SplitN(parts[0], []byte{32}, -1)

	temp_fmt := header[0]

	var obj GitObject = nil

	switch string(temp_fmt) {
	case "blob":
		fmt.Println("returning blob")
		//fmt.Println(string(Byte_data))
		obj = GitBlob{GitObjectData: GitObjectData{Byte_data}, format: []byte("blob")}
	case "commit":
		fmt.Println("returning commit")
		//fmt.Println(string(Byte_data))
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

func MakeGitObjWithFormat(Byte_data []byte, Obj_format string) GitObject {
	var obj GitObject = nil

	switch string(Obj_format) {
	case "blob":
		fmt.Println("returning blob")
		//fmt.Println(string(Byte_data))
		obj = GitBlob{GitObjectData: GitObjectData{Byte_data}, format: []byte("blob")}
	case "commit":
		fmt.Println("returning commit")
		//fmt.Println(string(Byte_data))
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

func Object_Read(repo gitpath.GitRepository, sha string) GitObject {

	//splitting strings for path check
	first := sha[:2]
	rest := sha[2:]

	//getting obj_path
	obj_path := gitpath.Repo_Path(repo, "objects", first, rest)
	fmt.Println(obj_path)

	//reading blob with zlib
	decompressed := Read_Blob(obj_path)

	//make obj based on format
	obj := MakeGitObj(*decompressed)
	fmt.Println(obj)

	return nil
}

func BuildGitObjectToWrite(obj GitObject) []byte {
	size := []byte(strconv.Itoa(len(obj.Deserialize())))

	result := append([]byte(obj.Get_Format()), ' ')
	result = append(result, size...)
	result = append(result, 0x00)
	result = append(result, obj.Deserialize()...)

	return result
}

func Object_Write(obj GitObject, repo *gitpath.GitRepository) ([]byte, error) {

	result := BuildGitObjectToWrite(obj)

	sha := sha1.Sum(result)

	hexstring := fmt.Sprintf("%x", sha)

	fmt.Printf("SHA: %x \n", sha)

	if repo != nil {
		mkdir := true
		location_1 := hexstring[:2]
		location_2 := hexstring[2:]
		_ = gitpath.Repo_File(*repo, mkdir, "objects", location_1, location_2)
	}

	return sha[:], nil
}

func Hash_Object(file_path string, format string, repo gitpath.GitRepository) ([]byte, error) {

	data, err := os.ReadFile(file_path)
	if err != nil {
		panic(fmt.Sprintf("error opening file %v", file_path))
	}

	//fmt.Println(string(data))

	obj := MakeGitObjWithFormat(data, format)

	sha, _ := Object_Write(obj, &repo)

	//fmt.Println(string(obj.Deserialize()))
	//fmt.Println(obj.Get_Format())

	return sha, nil
}

func main() {
	required := true
	repo, _ := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)

	//fmt.Println(repo.GitDir)

	//66726975a39cdb5babb3a02e8761ecd37e1c7c49 blob

	//064a74c46df49ea6afb685b1004433026c81c152 commit

	//obj := Object_Read(*repo, "064a74c46df49ea6afb685b1004433026c81c152")
	//fmt.Println(obj)

	sha, _ := Hash_Object("settings", "blob", *repo)

	fmt.Printf("%x \n", sha)
}

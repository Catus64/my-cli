package GitObjLib

import (
	"bytes"
	"fmt"
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
		fmt.Println("returning blob")
		obj = GitBlob{GitObjectData: GitObjectData{parts[1]}, format: []byte("blob")}
	case "commit":
		fmt.Println("returning commit object")
		obj = &GitCommit{GitObjectData: GitObjectData{parts[1]}, format: []byte("commit")}
	case "tag":
		fmt.Println("returning tag")
	case "tree":
		fmt.Println("returning tree")
		x := bytes.IndexByte(Byte_data[0:], 0x00)
		GitTreeData := Byte_data[x+1:]
		obj = GitTreeLeaf{GitObjectData: GitObjectData{GitTreeData}, format: []byte("tree")}
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
		fmt.Println("returning blob")
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

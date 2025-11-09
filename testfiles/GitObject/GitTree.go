package GitObjLib

import (
	"bytes"
	"fmt"
)

type GitTreeLeaf struct {
	GitObjectData
	format []byte
	mode   []byte
	sha    []byte
	path   string
}

func (leaf GitTreeLeaf) String() string {
	mode := string(leaf.mode)
	sha := string(leaf.sha)
	return fmt.Sprintf("%s %s %s", mode, leaf.path, sha)
}

func (leaf GitTreeLeaf) Serialize() *[]byte {
	return nil
}

func (leaf GitTreeLeaf) Deserialize() []byte {
	return leaf.data
}

func (leaf GitTreeLeaf) Get_Format() string {
	return string(leaf.format)
}

func Tree_Parse_One(raw []byte, start int) (int, *GitTreeLeaf, error) {

	fmt.Printf("%x \n", raw)

	x := bytes.IndexByte(raw[start:], ' ')
	if x != 5 && x != 6 {
		return 0, nil, fmt.Errorf("bad tree object")
	}

	mode := raw[start:x]
	var temp = []byte{0x00}
	if len(mode) == 5 {
		mode = append(temp, mode...)
	}

	y := bytes.IndexByte(raw[x:], 0x00) + x

	fmt.Println("y: ", y)
	fmt.Println("x: ", x)
	fmt.Println("mode: ", string(mode))

	path := raw[x+1 : y]

	fmt.Println("path: ", string(path))
	return 0, nil, nil
}

func Tree_Parse(raw []byte) []byte {
	return nil
}

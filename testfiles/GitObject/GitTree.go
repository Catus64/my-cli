package GitObjLib

import (
	"bytes"
	"fmt"
)

type GitTreeLeaf struct {
	GitObjectData
	format []byte
	mode   []byte
	sha    string
	path   []byte
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
	return string("leaf")
}

func Tree_Parse_One(raw []byte, start int) (int, *GitTreeLeaf, error) {

	//fmt.Printf("%x \n", raw)

	x := bytes.IndexByte(raw[start:], ' ') + start
	if x-start != 5 && x-start != 6 {
		//panic("bad tree object")
		return 0, nil, fmt.Errorf("bad tree object")
	}

	mode := raw[start:x]
	var temp = []byte{'0'}
	if len(mode) == 5 {
		mode = append(temp, mode...)
	}

	//fmt.Println("x: ", x)
	y := bytes.IndexByte(raw[x:], 0x00) + x

	//fmt.Println("y: ", y)
	fmt.Println("mode: ", string(mode))

	path := raw[x+1 : y]

	fmt.Println("path: ", string(path))

	rawsha := raw[y+1 : y+21]

	sha := fmt.Sprintf("%040x", rawsha)

	fmt.Printf("sha: %v \n", sha)

	leaf := GitTreeLeaf{
		mode: mode,
		sha:  sha,
		path: path,
	}

	//y + 21 is when the first entry ends
	return y + 21, &leaf, nil
}

func Tree_Parse(raw []byte) []GitTreeLeaf {

	var leaves []GitTreeLeaf
	position := 0
	max := len(raw)
	for position < max {
		new_position, leaf, err := Tree_Parse_One(raw, position)
		if err != nil {
			break
		}
		leaves = append(leaves, *leaf)
		position = new_position
	}

	return leaves
}

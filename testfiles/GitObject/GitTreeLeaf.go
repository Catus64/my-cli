package GitObjLib

import (
	"bytes"
	"fmt"
	"sort"
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

/*
func (leaf GitTreeLeaf) Tree_Leaf_Sort_Key() []byte {
	mode := string(leaf.mode[0:2])
	if mode == "10" {
		return leaf.path
	} else {
		return append(leaf.path, 0x2f) //adding '/' for directories
	}
}
*/

func Tree_Serialize(items []GitTreeLeaf) []byte {
	ret := []byte{}

	sort.Slice(items, func(i, j int) bool {

		modei := string(items[i].mode[0:2])
		modej := string(items[j].mode[0:2])

		if modej == "04" && modei != "04" {
			return true
		}

		return string(items[i].path) < string(items[j].path)
	})

	for _, item := range items {
		ret = append(ret, item.mode...)
		ret = append(ret, ' ')
		ret = append(ret, item.path...)
		ret = append(ret, 0x00)
		sha_bytes := make([]byte, 20)
		fmt.Sscanf(item.sha, "%040x", &sha_bytes)
		ret = append(ret, sha_bytes...)
		fmt.Println("Serialized item:", item.String())
	}
	return ret
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

	y := bytes.IndexByte(raw[x:], 0x00) + x

	path := raw[x+1 : y]

	rawsha := raw[y+1 : y+21]

	sha := fmt.Sprintf("%040x", rawsha)

	leaf := GitTreeLeaf{
		mode: mode,
		sha:  sha,
		path: path,
	}

	//fmt.Println("mode: ", string(mode))
	//fmt.Println("path: ", string(path))
	//fmt.Printf("sha: %v \n", sha)
	//fmt.Println("Parsed item:", leaf.String())

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

package GitObjLib

import (
	"bytes"
	"fmt"
	"sort"
)

type GitTreeLeaf struct {
	GitObjectData
	// format []byte
	Mode []byte
	Sha  string
	Path []byte
}

func (leaf GitTreeLeaf) String() string {
	mode := string(leaf.Mode)
	sha := string(leaf.Sha)
	return fmt.Sprintf("%s %s %s", mode, leaf.Path, sha)
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

	// to sort tree items so that directories come first, then files, both in alphabetical order

	sort.Slice(items, func(i, j int) bool {

		modei := string(items[i].Mode[0:2])
		modej := string(items[j].Mode[0:2])

		if modej == "04" && modei != "04" {
			return true
		}

		return string(items[i].Path) < string(items[j].Path)
	})

	for _, item := range items {
		ret = append(ret, item.Mode...)
		ret = append(ret, ' ')
		ret = append(ret, item.Path...)
		ret = append(ret, 0x00)
		sha_bytes := make([]byte, 20)
		fmt.Sscanf(item.Sha, "%040x", &sha_bytes)
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
		Mode: mode,
		Sha:  sha,
		Path: path,
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

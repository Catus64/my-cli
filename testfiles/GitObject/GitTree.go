package GitObjLib

import (
	gitpath "gocmd/testfiles/Gitrepostruct"
)

type GitTree struct {
	GitObjectData
	format []byte
	items  []GitTreeLeaf
}

func (tree GitTree) Get_Format() string {
	return string(tree.format)
}

func (tree GitTree) Serialize() *[]byte {
	return nil
}

func (tree GitTree) Deserialize() []byte {
	return tree.data
}

func (tree *GitTree) DeserializeData(data []byte) error {
	tree.items = Tree_Parse(data)
	return nil
}

func Ls_Tree(repo gitpath.GitRepository, ref string, recursive bool, prefix string) {

}

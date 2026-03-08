package gitCheckout

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
	"path/filepath"
)

func PrepareCheckoutDir(repo gitpath.GitRepository, path string) (string, error) {

	work_path := repo.WorkTree
	output_path := filepath.Join(work_path, "..", path)

	info, err := os.Stat(output_path)

	// Directory does not exist
	if os.IsNotExist(err) {
		return output_path, os.MkdirAll(output_path, 0755)
	}

	if err != nil {
		return "", err
	}

	// If something exists but isn't a directory
	if !info.IsDir() {
		return "", fmt.Errorf("path exists but is not a directory: %s", output_path)
	}

	// Ask user for overwrite
	fmt.Printf("Directory %s already exists. Overwrite? (y/n): ", output_path)

	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		return "", fmt.Errorf("checkout aborted by user")
	}

	// Delete directory contents
	err = os.RemoveAll(output_path)
	if err != nil {
		return "", err
	}

	// Recreate directory
	return output_path, os.MkdirAll(output_path, 0755)
}

func Load(arg string, name string, repo gitpath.GitRepository) (string, error) {

	//upcasting to commit object to access its fields
	commit := githashread.Object_Read(repo, arg)
	concreteCommit, ok := commit.(*gitobj.GitCommit)
	if !ok {
		panic("not a commit object")
	}
	concreteCommit.Deserialize()
	commit_tree := concreteCommit.KvlmDict.Dict["tree"]

	tree := githashread.Object_Read(repo, string(commit_tree))
	tree.Deserialize()

	//upcasting to tree object to access its fields
	concreteTree, ok := tree.(gitobj.GitTree)
	if !ok {
		panic("not a tree object")
	}

	// Deserializing the tree to get its items
	concreteTree.DeserializeData(tree.Deserialize())
	fmt.Println(string(concreteTree.Items[6].Path))

	//checkout tree to sister directory
	target_path, err := PrepareCheckoutDir(repo, name)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(target_path, 0o755)
	err = TreeCheckout(repo, concreteTree, target_path)
	if err != nil {
		return "", err
	}
	return target_path, nil
}

func TreeCheckout(repo gitpath.GitRepository, tree gitobj.GitTree, path string) error {

	for _, item := range tree.Items {

		obj := githashread.Object_Read(repo, item.Sha)

		dest := filepath.Join(path, string(item.Path))

		if obj.Get_Format() == "tree" {

			err := os.Mkdir(dest, 0o755)
			if err != nil {
				return err
			}

			subTree := obj.(gitobj.GitTree)

			err = TreeCheckout(repo, subTree, dest)
			if err != nil {
				return err
			}

		} else if obj.Get_Format() == "blob" {

			blob := obj.(gitobj.GitBlob)

			err := os.WriteFile(dest, blob.Deserialize(), 0o644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

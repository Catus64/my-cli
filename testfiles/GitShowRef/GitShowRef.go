package gitshowref

import (
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
	"path/filepath"
)

func Ref_list(repo gitpath.GitRepository, path string, temp_path string) map[string]string {

	if path == "" {
		ref_path, err := gitpath.Repo_Dir(repo, false, "refs")
		if err != nil {
			panic(err)
		}
		path = ref_path
		temp_path = "refs"
	}

	refs := map[string]string{}

	//fmt.Println("PATH:", path)

	//get entries in directory
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		entry_path := filepath.Join(path, entry.Name())
		info, _ := os.Stat(entry_path)

		//fmt.Println(entry_path)

		//if directory recurse else add ref
		if info.IsDir() {
			temp_path := filepath.Join(temp_path, entry.Name())
			entries_recursive := Ref_list(repo, entry_path, temp_path)
			//after recursion manually add entries to refs map
			for k, v := range entries_recursive {
				refs[k] = v
			}
		} else {
			temp_path := filepath.Join(temp_path, entry.Name())
			refs[temp_path] = *gitobj.Ref_Resolve(repo, entry_path)
		}
	}

	return refs
}

func Show_Ref(repo gitpath.GitRepository, ref string) {

}

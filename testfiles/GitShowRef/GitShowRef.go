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
			return map[string]string{} // no refs dir at all — not fatal
		}
		path = ref_path
		temp_path = "refs"
	}
	refs := map[string]string{}
	entries, err := os.ReadDir(path)
	if err != nil {
		return refs //refs/remotes doesn't exist, ignore
	}
	for _, entry := range entries {
		entry_path := filepath.Join(path, entry.Name())
		info, err := os.Stat(entry_path)
		if err != nil {
			continue // broken symlink maybe
		}
		if info.IsDir() {
			sub_temp_path := filepath.Join(temp_path, entry.Name())
			for k, v := range Ref_list(repo, entry_path, sub_temp_path) {
				refs[k] = v
			}
		} else {
			ref_name := filepath.Join(temp_path, entry.Name())
			ret, err := gitobj.Ref_Resolve(repo, ref_name)
			if err != nil {
				continue // unresolvable ref, no panic yet
			}
			refs[ref_name] = *ret
		}
	}
	return refs
}

func Show_Ref(repo gitpath.GitRepository, ref string) {

}

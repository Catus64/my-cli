package gitcurrent

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	"os"
	"path/filepath"
	"strings"
)

// Check Active Branch from .git/HEAD
func Get_Active_Branch(repo gitpath.GitRepository) (string, error) {
	headPath := gitpath.Repo_Path(repo, "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}
	logger.L().Debug("HEAD content read", "content", string(data))

	head := strings.TrimSpace(string(data))
	if strings.HasPrefix(head, "ref: refs/heads/") {
		return strings.TrimPrefix(head, "ref: refs/heads/"), nil
	}
	return "", nil // detached HEAD, handle separately
}

func Get_Tree_SHA(repo gitpath.GitRepository, ref string) (string, error) {
	commitSHA, _ := gitobj.Ref_Resolve(repo, ref)
	logger.L().Debug("Resolved ref to SHA", "ref", ref, "sha", *commitSHA)

	Commit_Object, err := githashread.Object_Read(repo, ref)
	if err != nil {
		return "", err
	}
	Commit_Object.Deserialize()
	Concrete_Commit, ok := Commit_Object.(*gitobj.GitCommit)
	if !ok {
		panic("not a commit object")
	}
	treeSHA := string(Concrete_Commit.Dict["tree"])
	logger.L().Debug("Commit object deserialized", "TREESHA", treeSHA)

	return treeSHA, nil
}

func TreeToMap(repo gitpath.GitRepository, ref string, prefix string, result map[string]string) error {
	// Get tree SHA from ref (use Get_Tree_SHA only at top level, then read tree directly)
	// ref here can be either a commit ref ("HEAD") or a raw tree SHA from recursion

	// Object_Read your tree object here, should return a *GitTree or similar
	treeObj, err := githashread.Object_Read(repo, ref)
	if err != nil {
		return err
	}

	leafs := gitobj.Tree_Parse(treeObj.Deserialize())

	for _, leaf := range leafs {

		fullPath := filepath.Join(prefix, string(leaf.Path))

		// mode starting with "04" means it's a subtree (directory)
		if strings.HasPrefix(string(leaf.Mode), "04") {
			err := TreeToMap(repo, leaf.Sha, fullPath, result)
			if err != nil {
				return err
			}
		} else {
			result[fullPath] = leaf.Sha
		}
	}
	return nil
}

func StatusHeadIndex(repo gitpath.GitRepository, index gitobj.GitIndex) error {
	fmt.Println("Changes to be committed:")

	// Build flat map of HEAD tree: path -> sha
	head := make(map[string]string)

	treeSHA, err := Get_Tree_SHA(repo, "HEAD")
	if err != nil {
		return err
	}

	err = TreeToMap(repo, treeSHA, "", head)
	if err != nil {
		// No commits yet, everything in index is new
		fmt.Println("No commits yet")
		// for _, entry := range index.Entries {
		// 	fmt.Println("  added:    ", entry.Name)
		// }
		return nil
	}

	logger.L().Debug("HEAD tree mapped to path->sha", "head", head)

	// Compare index entries against HEAD
	for _, entry := range index.Entries {
		if headSHA, exists := head[entry.Name]; exists {
			logger.L().Debug("Comparing index entry to HEAD", "entry", entry.Name, "indexSHA", entry.SHA, "headSHA", headSHA)
			if headSHA != entry.SHA {
				fmt.Println("  modified: ", entry.Name)
			}
			delete(head, entry.Name) // mark as seen
		} else {
			fmt.Println("  added:    ", entry.Name)
		}
	}

	// Anything left in head was not in the index — deleted
	for path := range head {
		fmt.Println("  deleted:  ", path)
	}

	return nil
}

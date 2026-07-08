package prettyprint

import (
	"fmt"
	gitsave "gocmd/testfiles/GitSave"
	gitshowref "gocmd/testfiles/GitShowRef"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"strings"
)

func ResolveCommitLabel(repo gitpath.GitRepository, sha string, branch string) string {
	// Priority 1 - Tag name via set name
	refs := gitshowref.Ref_list(repo, "", "")
	for refName, refSHA := range refs {
		if refSHA == sha && strings.HasPrefix(refName, "refs/tags/") {
			// extract just the tag name
			tagName := strings.TrimPrefix(refName, "refs/tags/")
			return fmt.Sprintf("tag: %s", tagName)
		}
	}

	// Priority 2 - ezgit automated tagging
	if branch != "" {
		entry, err := gitsave.ReadVersionRef(repo, branch, sha)
		if err == nil {
			return fmt.Sprintf("v%d · %s | %s", entry.Number, entry.Name, branch)
		}
	}

	// Priority 3 - nothing(not tagged and non ezgit made version/commit)
	return ""
}

//8ba6c00223f5cc9d20ac5bec0f47aa9a678f36f3
//f4b260b6d5d81ca9111ce994181f46654e412bbc
//434e36349e2b18bd9aa1d28e31f47173e82f857c
//bd9654ba920809a88b1ab1d31744977032040341

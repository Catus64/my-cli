package remove

import (
	"fmt"
	gitaddremove "gocmd/testfiles/GitAddRemove"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	prettyprint "gocmd/testfiles/PrettyPrint"

	"github.com/spf13/cobra"
)

func remove(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return fmt.Errorf("specify files to remove")
	}

	// snapshot index before removal so we know each requested file's SHA
	beforeIndex, err := gitobj.Index_Read2(*repo)
	if err != nil {
		return err
	}
	beforeMap := make(map[string]string, len(beforeIndex.Entries))
	for _, e := range beforeIndex.Entries {
		beforeMap[e.Name] = e.SHA
	}

	afterIndex, err := gitaddremove.Remove(repo, args, gitaddremove.RemoveOptions{
		Delete:          false,
		SkipMissingFile: true,
	})
	if err != nil {
		return err
	}
	afterSet := make(map[string]struct{}, len(afterIndex.Entries))
	for _, e := range afterIndex.Entries {
		afterSet[e.Name] = struct{}{}
	}

	var rows []prettyprint.RemoveRow
	var skipped []string
	for _, path := range args {
		sha, wasTracked := beforeMap[path]
		if !wasTracked {
			skipped = append(skipped, path)
			continue
		}
		if _, stillPresent := afterSet[path]; stillPresent {
			skipped = append(skipped, path)
			continue
		}
		rows = append(rows, prettyprint.RemoveRow{File: path, SHA: sha})
	}

	prettyprint.PrintRemoveTable(rows, skipped)
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [files]",
		Aliases: []string{"rm"},
		Short:   "Remove files from the entry list/index",
		RunE:    remove,
	}

	return cmd
}

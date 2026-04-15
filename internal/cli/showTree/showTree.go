package showTree

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func LsTree(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}
	obj, err := githashread.Object_Read(*repo, args[0])
	if err != nil {
		panic(err)
	}
	leafs := gitobj.Tree_Parse(obj.Deserialize())

	//formatting required later

	// for _, leaf := range leafs {
	// 	if string(leaf.Mode) == "40000" {
	// 		println(leaf.Mode, " tree ", leaf.Sha, "\t", string(leaf.Path))
	// 		continue
	// 	}
	// 	println(leaf.Mode, " blob ", leaf.Sha, "\t", string(leaf.Path))
	// }
	PrintTreeWithDots(args[0], leafs)
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-tree [tree-hash]",
		Short: "Show contents of Tree object",
		Args:  cobra.MaximumNArgs(1),
		RunE:  LsTree,
	}

	return cmd
}

func PrintTreeWithDots(sha string, leafs []gitobj.GitTreeLeaf) {
	const width = 90

	repeat := func(s string, n int) string {
		out := ""
		for i := 0; i < n; i++ {
			out += s
		}
		return out
	}

	row := func(text string) {
		fmt.Printf("│ %-*s │\n", width-2, text)
	}

	dotRow := func() {
		fmt.Printf("│%s│\n", repeat(".", width))
	}

	// Top of box
	fmt.Printf("┌%s┐\n", repeat("─", width))
	row("Tree: " + sha)
	fmt.Printf("├%s┤\n", repeat("─", width))
	row("Mode   Type   SHA                   Path")

	// Print each leaf with dotted separator
	for i, leaf := range leafs {
		objType := "blob"
		if string(leaf.Mode) == "040000" {
			objType = "tree"
		}
		row(fmt.Sprintf("%s %-6s %s  %s", leaf.Mode, objType, leaf.Sha, leaf.Path))

		if i != len(leafs)-1 { // no dot after last
			dotRow()
		}
	}

	// Bottom of box
	fmt.Printf("└%s┘\n", repeat("─", width))
}

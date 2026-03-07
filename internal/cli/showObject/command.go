package showObject

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func show(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}
	obj := githashread.Object_Read(*repo, args[0])

	// fmt.Println("Object format:", obj.Get_Format())

	if obj.Get_Format() == "tree" || obj.Get_Format() == "leaf" {
		println("This is a tree object. please use the 'ls-tree' command to view its contents.")
		return nil
	}
	//fmt.Println(string(obj.Deserialize()))
	content := string(obj.Deserialize())
	PrintObjectContent(args[0], []byte(content))
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-object [object-hash]",
		Short: "Print a friendly greeting",
		Args:  cobra.MaximumNArgs(1),
		RunE:  show,
	}

	return cmd
}

func PrintObjectContent(sha string, content []byte) {
	const width = 65 // adjust if object content is long

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

	lines := func(text string) []string {
		out := []string{}
		current := ""
		for _, r := range text {
			if r == '\n' {
				out = append(out, current)
				current = ""
				continue
			}
			current += string(r)
		}
		if current != "" {
			out = append(out, current)
		}
		return out
	}

	fmt.Printf("┌%s┐\n", repeat("─", width))
	row("Object: " + sha)
	fmt.Printf("├%s┤\n", repeat("─", width))
	row("")

	for _, line := range lines(string(content)) {
		row(line)
	}

	row("")
	fmt.Printf("└%s┘\n", repeat("─", width))
}

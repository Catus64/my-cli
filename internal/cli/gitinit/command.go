package gitinit

import (
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func Git_init(cmd *cobra.Command, args []string) error {
	// fmt.Println("making repo")
	rootpath := gitpath.Get_Os_Dir()
	gitrepo := gitpath.Repo_create(rootpath)
	PrintRepoInitialized(gitrepo.GitDir)
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-repo",
		Short: "Create a new git repository",
		RunE:  Git_init,
	}
	return cmd
}

func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}

func center(s string, width int) string {
	if len(s) >= width {
		return s
	}
	padding := (width - len(s)) / 2
	return repeat(" ", padding) + s
}

func PrintRepoInitialized(path string) {
	const width = 69

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

	fmt.Printf("┌%s┐\n", repeat("─", width))
	row(center("Repository Initialized Successfully", width-2))
	fmt.Printf("├%s┤\n", repeat("─", width))

	row("Path : " + path)
	row("")
	row("An empty repository has been created and")
	row("is ready for use.")

	fmt.Printf("└%s┘\n", repeat("─", width))
}

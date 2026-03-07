package hashObject

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func hash(cmd *cobra.Command, args []string) error {
	required := true
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		panic(err)
	}
	sha, _ := githashread.Hash_Object(args[0], "blob", *repo)
	// fmt.Printf("file: {%s} has been compressed and stored in the objects folder\n", args[0])
	// fmt.Printf("SHA: %x \n", sha)
	// fmt.Println("To search for the file in the repository please use this hash to find it later")
	PrintObjectStored("blob", args[0], fmt.Sprintf("%x", sha))
	return nil
}

// Ai slop will be rewritten
func PrintObjectStored(objectType, fileName, sha string) {
	const width = 55

	line := func() {
		fmt.Printf("┌%s┐\n", repeat("─", width))
	}
	sep := func() {
		fmt.Printf("├%s┤\n", repeat("─", width))
	}
	end := func() {
		fmt.Printf("└%s┘\n", repeat("─", width))
	}

	row := func(text string) {
		fmt.Printf("│ %-*s │\n", width-2, text)
	}

	line()
	row(center("Object Stored Successfully", width-2))
	sep()
	row(fmt.Sprintf("Type : %s", objectType))
	row(fmt.Sprintf("File : %s", fileName))
	row("")
	row(fmt.Sprintf("SHA  : %s", sha))
	row("")
	row("The object has been compressed and written")
	row("to the repository object database.")
	sep()
	row("Use the SHA above to reference this object.")
	end()
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

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hash-file [object-hash]",
		Short: "Print a friendly greeting",
		Args:  cobra.MaximumNArgs(1),
		RunE:  hash,
	}

	return cmd
}

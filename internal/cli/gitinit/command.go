package gitinit

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Git_init(cmd *cobra.Command, args []string) error {
	fmt.Println("making repo")
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a new git repository",
		RunE:  Git_init,
	}
	return cmd
}

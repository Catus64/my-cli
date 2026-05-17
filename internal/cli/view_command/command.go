package viewcommand

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view",
		Aliases: []string{"vw"},
		Short:   "View repository objects, trees, refs and versions",
	}

	cmd.AddCommand(newResolveCommand())
	cmd.AddCommand(newTreeCommand())
	cmd.AddCommand(newObjectCommand())
	cmd.AddCommand(newRefsCommand())
	cmd.AddCommand(newVersionCommand())

	return cmd
}

package greet

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "greet [name]",
		Short: "Print a friendly greeting",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "World"

			if len(args) == 1 {
				name = args[0]

				if args[0] == "all" {
					name = "Everyone"
				}
			}

			fmt.Printf("Greetings, %s!\n", name)
			return nil
		},
	}
	return cmd
}

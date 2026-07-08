package greet

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func Greet(cmd *cobra.Command, args []string) error {
	name := "World"

	if len(args) == 1 {
		name = args[0]

		if args[0] == "all" {
			name = "Everyone"
		}
	}

	shout, _ := cmd.Flags().GetBool("shout")
	if shout {
		name = strings.ToUpper(name)
	}

	fmt.Printf("Greetings, %s!\n", name)
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "greet [name]",
		Short: "Print a friendly greeting",
		Args:  cobra.MaximumNArgs(1),
		RunE:  Greet,
	}

	cmd.Flags().BoolP("shout", "s", false, "Shout the greeting in uppercase")

	return cmd
}

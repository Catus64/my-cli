package extract

import (
	"fmt"

	"github.com/spf13/cobra"

	extractor "gocmd/testfiles/GitPacketExtractor" // adjust to your actual extractor import path
	gitpath "gocmd/testfiles/Gitrepostruct"
)

func extract(cmd *cobra.Command, args []string) error {
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), true)
	if err != nil {
		return err
	}

	if err := extractor.ForceExtract(*repo); err != nil {
		return fmt.Errorf("failed to extract packfiles: %w", err)
	}

	fmt.Println("Packfiles re-extracted successfully.")
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract",
		Short: "extract git's packfiles(try this command if you have issues finding objects)",
		RunE:  extract,
	}
	return cmd
}

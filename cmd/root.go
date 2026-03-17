/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gocmd/internal/cli/gitLog"
	gittag "gocmd/internal/cli/gitTag"
	"gocmd/internal/cli/gitinit"
	"gocmd/internal/cli/hashObject"
	"gocmd/internal/cli/load"
	showref "gocmd/internal/cli/show-ref"
	showsavelist "gocmd/internal/cli/show-savelist"
	"gocmd/internal/cli/showObject"
	"gocmd/internal/cli/showTree"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ezgit",
	Short: "Lightweight Versioning Engine",
	Long: `This is a lightweight versioning engine similar to git.
	It has a smaller call set and is intended for Single User use.

	You can use it to create repositories, add files, commit changes, view logs and tags, and checkout previous versions of your files.
	This tool is designed to be simple and easy to use, making it ideal for personal projects and small teams.
	The command line interface is intuitive, allowing users to quickly learn and utilize its features effectively
	while also being transparent about what happens under the hood.
	
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gocmd.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(gitinit.NewCommand())
	rootCmd.AddCommand(showObject.NewCommand())
	rootCmd.AddCommand(hashObject.NewCommand())
	rootCmd.AddCommand(gitLog.NewCommand())
	rootCmd.AddCommand(showTree.NewCommand())
	rootCmd.AddCommand(load.NewCommand())
	rootCmd.AddCommand(showref.NewCommand())
	rootCmd.AddCommand(gittag.NewCommand())
	rootCmd.AddCommand(showsavelist.NewCommand())
}

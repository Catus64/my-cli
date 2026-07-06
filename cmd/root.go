/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gocmd/internal/cli/add"
	altversion "gocmd/internal/cli/alt-version"
	checkignore "gocmd/internal/cli/check-ignore"
	"gocmd/internal/cli/combine"
	getcurrent "gocmd/internal/cli/get-current"
	gitsave "gocmd/internal/cli/git-save"
	"gocmd/internal/cli/gitLog"
	gittag "gocmd/internal/cli/gitTag"
	gitignoreeditor "gocmd/internal/cli/gitignore-editor"
	"gocmd/internal/cli/gitinit"
	"gocmd/internal/cli/hashObject"
	"gocmd/internal/cli/load"
	"gocmd/internal/cli/quicksave"
	"gocmd/internal/cli/remove"
	setconfig "gocmd/internal/cli/set-config"
	prettyprint "gocmd/testfiles/PrettyPrint"

	// showref "gocmd/internal/cli/show-ref"
	showsavelist "gocmd/internal/cli/show-savelist"
	// "gocmd/internal/cli/showTree"
	switchver "gocmd/internal/cli/switch-ver"
	view "gocmd/internal/cli/view_command"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ezg",
	Short: "Lightweight Versioning Engine",
	Long: `This is a lightweight versioning engine similar to git.
	It has a smaller call set and is intended for Single User use.

	You can use it to create repositories, add files, commit changes, view logs and tags, and checkout previous versions of your files.
	This tool is designed to be simple and easy to use, making it ideal for personal projects and small teams.
	The command line interface is intuitive, allowing users to quickly learn and utilize its features effectively
	while also being transparent about what happens under the hood.

	If you are unsure what to do just type 'ezg quicksave' and 'ezg history'
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		printQuickHelp()
	},
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

const ezgLogo = `
__/\\\\\\\\\\\\\\\__/\\\\\\\\\\\\\\\_____/\\\\\\\\\\\\__/\\\\\\\\\\\__/\\\\\\\\\\\\\\\_        
 _\/\\\///////////__\////////////\\\____/\\\//////////__\/////\\\///__\///////\\\/////__       
  _\/\\\_______________________/\\\/____/\\\_________________\/\\\___________\/\\\_______      
   _\/\\\\\\\\\\\_____________/\\\/_____\/\\\____/\\\\\\\_____\/\\\___________\/\\\_______     
    _\/\\\///////____________/\\\/_______\/\\\___\/////\\\_____\/\\\___________\/\\\_______    
     _\/\\\_________________/\\\/_________\/\\\_______\/\\\_____\/\\\___________\/\\\_______   
      _\/\\\_______________/\\\/___________\/\\\_______\/\\\_____\/\\\___________\/\\\_______  
       _\/\\\\\\\\\\\\\\\__/\\\\\\\\\\\\\\\_\//\\\\\\\\\\\\/___/\\\\\\\\\\\_______\/\\\_______ 
        _\///////////////__\///////////////___\////////////____\///////////________\///________
`

func printQuickHelp() {
	fmt.Print(ezgLogo)
	const width = 100

	prettyprint.Top(width)
	prettyprint.Row(prettyprint.Center("EzGit: Simplified Versioning Engine", width), width)
	prettyprint.Mid(width)
	prettyprint.Row("Start with:", width)
	prettyprint.EmptyRow(width)
	prettyprint.Row("  ezg create-repo       create a new repository", width)
	prettyprint.Row("  ezg quicksave         quickly save your current changes", width)
	prettyprint.Row("  ezg history           view past versions", width)
	prettyprint.Mid(width)
	prettyprint.Row("Common commands:", width)
	prettyprint.EmptyRow(width)
	prettyprint.Row("  ezg add <file>        stage a file(use -a for all files)", width)
	prettyprint.Row("  ezg save              save a version with a message", width)
	prettyprint.Row("  ezg name/tag          name/tag a version", width)
	prettyprint.Row("  ezg view refs         view all named/tagged versions", width)
	prettyprint.Row("  ezg list              view files that are about to be saved (after ezg add)", width)
	prettyprint.Mid(width)
	prettyprint.Row(prettyprint.Center("Run 'ezg -h' to see all available commands", width), width)
	prettyprint.Bottom(width)
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddGroup(
		&cobra.Group{ID: "setup", Title: "Setup and Config Commands:"},
		&cobra.Group{ID: "save", Title: "Saving related commands:"},
		&cobra.Group{ID: "branch", Title: "Branching and Merging Commands:"},
		&cobra.Group{ID: "view", Title: "View objects and history:"},
		&cobra.Group{ID: "hash", Title: "Low-Level commands"},
	)

	addToGroup("setup",
		gitinit.NewCommand(),
		setconfig.NewCommand(),
		gitignoreeditor.NewCommand(),
		checkignore.NewCommand(),
	)

	addToGroup("save",
		add.NewCommand(),
		remove.NewCommand(),
		gitsave.NewCommand(),
		quicksave.NewCommand(),
		gittag.NewCommand(),
		load.NewCommand(),
	)

	addToGroup("branch",
		altversion.NewCommand(),
		switchver.NewCommand(),
		combine.NewCommand(),
		getcurrent.NewCommand(),
	)
	addToGroup("view",
		gitLog.NewCommand(),
		showsavelist.NewCommand(),
		view.NewCommand(),
	)

	addToGroup("hash",
		hashObject.NewCommand(),
	)

	// Commands commented .
	// rootCmd.AddCommand(showObject.NewCommand())
	// rootCmd.AddCommand(showTree.NewCommand())
	// rootCmd.AddCommand(showref.NewCommand())
}

// addToGroup assigns a GroupID to each command and registers it on rootCmd.
func addToGroup(groupID string, cmds ...*cobra.Command) {
	for _, c := range cmds {
		c.GroupID = groupID
		rootCmd.AddCommand(c)
	}
}

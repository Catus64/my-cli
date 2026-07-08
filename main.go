/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"gocmd/cmd"
	extractor "gocmd/testfiles/GitPacketExtractor"
	gitpath "gocmd/testfiles/Gitrepostruct"
	helper "gocmd/testfiles/Helper"
	"log/slog"
	"os"
	"path/filepath"
)

func main() {

	required := false
	repo, err := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)
	if err != nil {
		slog.Warn("no Repository found")
	}

	if repo != nil {

		logpath := gitpath.Repo_Path(*repo, "ezgit.log")
		err = helper.Init(logpath, slog.LevelDebug)
		if err != nil {
			panic(err)
		}

		if err := extractor.Extract(*repo); err != nil {
			slog.Error("failed to extract packfiles", "error", err)
		}
	}

	if repo == nil {
		var path string
		configDir, err := os.UserConfigDir()
		if err != nil {
			// fallback if os.UserConfigDir() fails
			home, err := os.UserHomeDir()
			if err != nil {

			}
			path = filepath.Join(home, ".ezgit", "ezgit.log")
		} else {
			path = filepath.Join(configDir, "ezgit", "config")

		}
		err = helper.Init(path, slog.LevelDebug)
		if err != nil {
			panic(err)
		}
	}

	cmd.Execute()
}

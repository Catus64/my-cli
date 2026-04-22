/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"gocmd/cmd"
	extractor "gocmd/testfiles/GitPacketExtractor"
	helper "gocmd/testfiles/Helper"
	"log/slog"
)

func main() {
	err := helper.Init("ezgit.log", slog.LevelDebug)
	if err != nil {
		panic(err)
	}

	if err := extractor.Extract(); err != nil {
		slog.Error("failed to extract packfiles", "error", err)
	}

	cmd.Execute()
}

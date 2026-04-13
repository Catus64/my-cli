/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"gocmd/cmd"
	helper "gocmd/testfiles/Helper"
	"log/slog"
)

func main() {
	err := helper.Init("ezgit.log", slog.LevelDebug)
	if err != nil {
		panic(err)
	}
	//testcomment
	//testmorethings
	cmd.Execute()
}

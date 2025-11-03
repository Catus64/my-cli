package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	path := filepath.Join("home", "user", "document", "file.txt")

	fmt.Println(path)

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("dir: ", dir)

	//pkgtest.Send_pkg()
}

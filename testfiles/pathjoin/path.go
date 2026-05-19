package main

import (
	"fmt"
	"strings"
)

func main() {
	pattern := "/.venv"
	// path := ".venv/something"

	// state := strings.Contains(path, pattern[1:])
	state := strings.HasPrefix(pattern, "/")
	fmt.Println(state)
}

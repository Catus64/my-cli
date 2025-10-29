package main

import (
	"fmt"
	"log"

	"gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load("settings.ini")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	fmt.Println(cfg.Section("test").Key("name").String())
}

package main

import (
	"fmt"
	"log"

	"gopkg.in/ini.v1"
)

func repo_default_config() (*ini.File, error) {
	cfg := ini.Empty()
	cfg.Section("core").Key("repositoryformatversion").SetValue("0")
	cfg.Section("core").Key("filemode").SetValue("false")
	cfg.Section("core").Key("bare").SetValue("false")
	return cfg, nil
}

func main() {
	cfg, err := ini.Load("settings")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	fmt.Println(cfg.Section("test").Key("name").String())

	cfg2, err := repo_default_config()
	if err != nil {
		log.Fatalf("screw GO")
	}
	cfg2.SaveTo("settings2")
}

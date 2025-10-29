package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type GitRepository struct {
	WorkTree string
	GitDir   string
	cfg      *ini.File
}

func MakeRepo(path string) GitRepository {
	repo := GitRepository{
		WorkTree: path,
		GitDir:   filepath.Join(path, ".git")}
	Load_ini(&repo)

	return repo
}

func Load_ini(repo *GitRepository) {
	tempcfg, err := ini.Load("settings.ini")
	if err != nil {
		panic(err)
	}
	repo.cfg = tempcfg
}

func Repo_path(repo GitRepository, path ...string) string {
	fmt.Println(path)
	paths := filepath.Join(path...)
	return filepath.Join(repo.GitDir, paths)
}

func Repo_Dir(repo GitRepository, paths ...string) (string, error) {

	path := Repo_path(repo, paths...)

	fmt.Println(path)

	return "", nil
}

func main() {

	path := filepath.Join("home", "user", "document", "file.txt")

	fmt.Println(path)

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	//fmt.Println("dir: ", dir)

	repo := MakeRepo(dir)

	fmt.Println(repo.GitDir)

	//res := Repo_path(repo, "one", "two", "three")
	//fmt.Println(res)
}

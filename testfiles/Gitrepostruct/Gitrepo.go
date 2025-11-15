package RepoPath

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/ini.v1"
)

// struct for initializing info about git stuff
type GitRepository struct {
	WorkTree string
	GitDir   string
	cfg      *ini.File
}

func MakeRepo(path string, force bool) *GitRepository {
	repo := GitRepository{
		WorkTree: path,
		GitDir:   filepath.Join(path, ".git")}

	repo.Read_Conf(force)

	if !force {
		vers := repo.cfg.Section("core").Key("repositoryformatversion").String()
		versnum, err := strconv.Atoi(vers)
		if err != nil || versnum != 0 {
			panic("Unsurported repositoryformatversion")
		}
	}
	return &repo
}

func (repo *GitRepository) Read_Conf(force bool) {
	read, err := ini.Load(filepath.Join(repo.GitDir, "config"))
	if err != nil {
		if !force {
			panic("Configuration file missing")
		} else {
			repo.cfg = nil
			return
		}
	}
	repo.cfg = read

}

func Repo_Path(repo GitRepository, path ...string) string {
	paths := filepath.Join(path...)
	return filepath.Join(repo.GitDir, paths)
}

func Repo_File(repo GitRepository, mkdir bool, paths ...string) string {
	//last file is an object therefore removed from dir path
	filepaths := paths[:len(paths)-1]
	//fmt.Println("Current path:", filepaths)

	_, err := Repo_Dir(repo, mkdir, filepaths...)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	return Repo_Path(repo, paths...)
}

func Repo_Dir(repo GitRepository, mkdir bool, paths ...string) (string, error) {

	path := Repo_Path(repo, paths...)
	//fmt.Println(path)

	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			//fmt.Println("dir exist")
			return path, nil
		}
		//panic("Not a Directory")
		return "", fmt.Errorf("not a directory/file already exist: %v", path)
	}

	if os.IsNotExist(err) {
		//fmt.Println("dir dont exist")
		// Mkdir here
		if mkdir {
			//fmt.Println("making dir")
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return "", fmt.Errorf("error creating directory %v", err)
			}
			fmt.Println("Directory has been created:", path)
			return path, nil
		} else {
			return "", fmt.Errorf("error while making files %v", path)
		}
	}

	return "", nil
}

func mustRepo_Dir(repo GitRepository, mkdir bool, paths ...string) {
	str, err := Repo_Dir(repo, mkdir, paths...)
	if err != nil {
		panic(err)
	} else if str == "" {
		panic("mkdir function fails")
	}
}

func Repo_create(path string) *GitRepository {
	//create a repository

	force := true
	repo := MakeRepo(path, force)

	//make sure path doesn't exist
	infoTree, err := os.Stat(repo.WorkTree)
	if err == nil && !infoTree.IsDir() {
		panic(fmt.Sprintf("%v is not a directory", path))
	}

	info, err := os.Stat(repo.GitDir)
	if err == nil {
		if info.IsDir() {
			panic(fmt.Sprintf("%v is not empty", path))
		}
	}

	if infoTree == nil {
		os.MkdirAll(repo.WorkTree, 0755)
	}

	mkdir := true

	mustRepo_Dir(*repo, mkdir, "branches")
	mustRepo_Dir(*repo, mkdir, "objects")
	mustRepo_Dir(*repo, mkdir, "refs", "tags")
	mustRepo_Dir(*repo, mkdir, "refs", "head")

	head_file_path := Repo_Path(*repo, "HEAD")
	content := "ref: refs/heads/master\n"

	err = os.WriteFile(head_file_path, []byte(content), 0644)
	if err != nil {
		panic("writing file fail")
	}

	config_file_path := Repo_Path(*repo, " config")
	config := Repo_default_config()
	config.SaveTo(config_file_path)

	return repo

}

func Repo_default_config() *ini.File {
	cfg := ini.Empty()
	cfg.Section("core").Key("repositoryformatversion").SetValue("0")
	cfg.Section("core").Key("filemode").SetValue("false")
	cfg.Section("core").Key("bare").SetValue("false")

	return cfg
}

func Repo_find(path string, required bool) (*GitRepository, error) {

	//assumes "." to get current location

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	//fmt.Println("Checking for git repo in", absPath)

	//## Checking if file exist ##//
	info, err := os.Stat(filepath.Join(absPath, ".git"))
	if err == nil {
		if info.IsDir() {
			//fmt.Println("Git repo exist")
			force := false
			repo := MakeRepo(absPath, force)
			return repo, nil
		}
	}

	parent_path := filepath.Join(absPath, "..")
	if parent_path == absPath {
		if required {
			panic("cannot find a repository")
		}
		return nil, fmt.Errorf("there is no git repo in sight ")
	}

	return Repo_find(parent_path, required)
}

func Get_Os_Dir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

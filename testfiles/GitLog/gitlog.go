package gitLog

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"strconv"
	"strings"
	"time"
)

// func Read_Master(repo gitpath.GitRepository) string {
// 	head := gitpath.Repo_Path(repo, "refs", "heads", "master")
// 	data, err := os.ReadFile(head)
// 	if err != nil {
// 		panic(err)
// 	}
// 	data = bytes.ReplaceAll(data, []byte(" "), []byte(""))
// 	data = bytes.ReplaceAll(data, []byte("\n"), []byte(""))

// 	return string(data)
// }

func Format_Date_Author(text string) (string, string) {
	temp_string := strings.SplitN(text, " ", -1)
	temp_time := temp_string[2]
	num, err := strconv.ParseInt(temp_time, 10, 64)
	if err != nil {
		panic(err)
	}
	loc, _ := time.LoadLocation("Local")
	t := time.Unix(num, 0).In(loc)

	author := temp_string[0] + temp_string[1]

	date := t.Format(time.RFC1123Z)

	return date, author
}

func Recurse_Log(repo *gitpath.GitRepository, commit gitobj.GitCommit, sha string) {
	const (
		Yellow = "\033[33m"
		Reset  = "\033[0m"
	)

	//log one obj
	//fmt.Println(Yellow+"Commit: ", sha+Reset)
	date, author := Format_Date_Author(string(commit.Dict["author"]))

	// fmt.Println("Date: ", date)
	// fmt.Println("Author: ", author)
	// fmt.Println("tree: ", string(commit.Dict["tree"]))
	// fmt.Println("parent: ", string(commit.Dict["parent"]))

	// fmt.Println(string(commit.Dict["data"]))

	PrintCommit(
		sha,
		author,
		date,
		string(commit.Dict["tree"]),
		string(commit.Dict["parent"]),
		string(commit.Dict["data"]),
	)

	if commit.Dict["parent"] == nil {
		println("end")
		return
	}

	//recurse
	Parent_Obj := githashread.Object_Read(*repo, string(commit.Dict["parent"]))
	Concrete_Parent_Commit, ok := Parent_Obj.(*gitobj.GitCommit)
	if !ok {
		panic("not a commit object")
	}
	Concrete_Parent_Commit.Deserialize()
	//println("\n")
	Recurse_Log(repo, *Concrete_Parent_Commit, string(commit.Dict["parent"]))
}

func Log(repo gitpath.GitRepository) error {
	// head := Read_Master(repo)
	point, err := gitobj.Ref_Resolve(repo, "HEAD")
	if err != nil {
		return err
	}
	// fmt.Println("HEAD points to:", point)
	head := *point
	Commit_Object := githashread.Object_Read(repo, head)
	Commit_Object.Deserialize()
	Concrete_Commit, ok := Commit_Object.(*gitobj.GitCommit)
	if !ok {
		panic("not a commit object")
	}
	//fmt.Println(string(Concrete_Commit.Dict["data"]))
	Recurse_Log(&repo, *Concrete_Commit, head)
	return nil
}

// AI slop will be rewritten
func PrintCommit(sha, author, date, tree, parent, message string) {
	const width = 55

	line := func(left, mid, right string) {
		fmt.Printf("%s%s%s\n", left, mid, right)
	}

	row := func(text string) {
		fmt.Printf("│ %-*s │\n", width-2, text)
	}

	repeat := func(s string, n int) string {
		out := ""
		for i := 0; i < n; i++ {
			out += s
		}
		return out
	}

	line("┌", repeat("─", width), "┐")
	row("Version")
	line("├", repeat("─", width), "┤")

	row("SHA    : " + sha)
	row("Author : " + author)
	row("Date   : " + date)
	row("")
	row("Tree   : " + tree)
	row("Parent : " + parent)

	line("├", repeat("─", width), "┤")
	row("Message")
	row("")

	for _, line := range splitLines(message) {
		row(line)
	}

	line("└", repeat("─", width), "┘")
}

func splitLines(s string) []string {
	out := []string{}
	current := ""

	for _, r := range s {
		if r == '\n' {
			out = append(out, current)
			current = ""
			continue
		}
		current += string(r)
	}
	if current != "" {
		out = append(out, current)
	}
	return out
}

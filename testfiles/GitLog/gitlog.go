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

	date, author := Format_Date_Author(string(commit.Dict["author"]))
	treeSHA := strings.TrimSpace(string(commit.Dict["tree"]))
	parents := gitobj.GetKvlmValues(commit.Dict, "parent")

	parentDisplay := ""
	if len(parents) > 0 {
		parentDisplay = parents[0]
		if len(parents) > 1 {
			parentDisplay += fmt.Sprintf(" (+%d more)", len(parents)-1)
		}
	}

	PrintCommit(
		sha,
		author,
		date,
		treeSHA,
		parentDisplay,
		string(commit.Dict["data"]),
	)

	if commit.Dict["parent"] == nil {
		println("end")
		return
	}

	//recurse
	if len(parents) == 0 {
		println("end")
		return
	}

	// follow first parent only for linear log
	firstParent := parents[0]

	Parent_Obj, err := githashread.Object_Read(*repo, firstParent)
	if err != nil {
		panic(err)
	}
	Concrete_Parent_Commit, ok := Parent_Obj.(*gitobj.GitCommit)
	if !ok {
		panic("not a commit object")
	}
	Concrete_Parent_Commit.Deserialize()
	Recurse_Log(repo, *Concrete_Parent_Commit, firstParent)
}

func Log(repo gitpath.GitRepository) error {
	// head := Read_Master(repo)
	point, err := gitobj.Ref_Resolve(repo, "HEAD")
	if err != nil {
		return err
	}
	// fmt.Println("HEAD points to:", point)
	head := *point
	Commit_Object, err := githashread.Object_Read(repo, head)
	if err != nil {
		return err
	}
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
	const width = 72        // total box width including borders
	const inner = width - 4 // "│ " + content + " │"

	border := strings.Repeat("─", width-2)

	printLine := func(left, mid, right string) {
		fmt.Printf("%s%s%s\n", left, mid, right)
	}

	row := func(text string) {
		// wrap long text into multiple lines
		for _, wrapped := range wrapText(text, inner) {
			fmt.Printf("│ %-*s │\n", inner, wrapped)
		}
	}

	printLine("┌", border, "┐")
	row("Version")
	printLine("├", border, "┤")
	row("SHA    : " + sha)
	row("Author : " + author)
	row("Date   : " + date)
	row("")
	row("Tree   : " + tree)
	if parent != "" {
		row("Parent : " + parent)
	}
	printLine("├", border, "┤")
	row("Message")
	row("")
	for _, line := range splitLines(message) {
		row(line)
	}
	printLine("└", border, "┘")
}

// wrapText splits text into lines of max `width` chars
func wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}
	var lines []string
	for len(text) > width {
		lines = append(lines, text[:width])
		text = text[width:]
	}
	if text != "" {
		lines = append(lines, text)
	}
	return lines
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

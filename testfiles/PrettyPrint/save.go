package prettyprint

import (
	"fmt"
	"time"
)

type SaveResult struct {
	Branch      string
	CommitSHA   string
	TreeSHA     string
	ParentSHA   string // empty if first commit
	Author      string
	Timestamp   time.Time
	Message     string
	VersionNum  int
	VersionName string
}

// func get_short_hash(sha string) string {
// 	if len(sha) >= 7 {
// 		return sha[:7]
// 	}
// 	return sha
// }

func PrintSaveResult(r SaveResult) {
	const width = 72

	Top(width)
	Row(Center("Version Saved Successfully", width), width)
	Mid(width)

	// savefile + branch
	Row(fmt.Sprintf("Savefile : %s", r.Branch), width)
	EmptyRow(width)

	// version info — most prominent
	Row(fmt.Sprintf("Version  : v%d · %s", r.VersionNum, r.VersionName), width)
	EmptyRow(width)

	// commit details
	Row(fmt.Sprintf("SHA      : %s", (r.CommitSHA)), width)
	Row(fmt.Sprintf("Tree     : %s", (r.TreeSHA)), width)
	if r.ParentSHA != "" {
		Row(fmt.Sprintf("Parent   : %s", (r.ParentSHA)), width)
	} else {
		Row("Parent   : none (first commit)", width)
	}
	EmptyRow(width)

	// author + time
	Row(fmt.Sprintf("Author   : %s", r.Author), width)
	Row(fmt.Sprintf("Time     : %s", r.Timestamp.Format("Mon, 02 Jan 2006 15:04:05 -0700")), width)
	EmptyRow(width)

	// version message
	Mid(width)
	Row("Message", width)
	Mid(width)
	EmptyRow(width)
	for _, line := range SplitLines(r.Message) {
		if line != "" {
			Row("  "+line, width)
		}
	}
	EmptyRow(width)

	// help message
	Mid(width)
	Row(Center("use 'ezgit view version' to browse all latest versions", width), width)
	Bottom(width)
	Mid(width)
	Row(Center("use 'ezgit set-name' to track a specific version", width), width)
	Bottom(width)
}

func PrintMergeResult(r SaveResult) {
	const width = 72
	Top(width)
	Row(Center("Version Merged Successfully", width), width)
	Mid(width)
	Row(fmt.Sprintf("Savefile : %s", r.Branch), width)
	EmptyRow(width)
	Row(fmt.Sprintf("Version  : v%d · %s", r.VersionNum, r.VersionName), width)
	EmptyRow(width)
	Row(fmt.Sprintf("SHA      : %s", r.CommitSHA), width)
	Row(fmt.Sprintf("Tree     : %s", r.TreeSHA), width)
	Row(fmt.Sprintf("Parents  : %s", r.ParentSHA), width)
	EmptyRow(width)
	Row(fmt.Sprintf("Author   : %s", r.Author), width)
	Row(fmt.Sprintf("Time     : %s", r.Timestamp.Format("Mon, 02 Jan 2006 15:04:05 -0700")), width)
	EmptyRow(width)
	Mid(width)
	Row("Merge Message", width)
	Mid(width)
	EmptyRow(width)
	for _, line := range SplitLines(r.Message) {
		if line != "" {
			Row("  "+line, width)
		}
	}
	EmptyRow(width)
	Mid(width)
	Row(Center("use 'ezgit view version' to browse all latest versions", width), width)
	Bottom(width)
	Mid(width)
	Row(Center("use 'ezgit set-name' to track a specific version", width), width)
	Bottom(width)
}

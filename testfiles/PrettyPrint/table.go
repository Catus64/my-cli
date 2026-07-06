package prettyprint

import (
	"fmt"
)

type TableRow struct {
	File   string
	SHA    string
	Status string
}

type RemoveRow struct {
	File string
	SHA  string
}

func PrintAddTable(rows []TableRow, ignored []string, isAll bool) {
	const width = 72
	fileW := 35
	shaW := 10
	statusW := 10

	newCount := 0
	modCount := 0
	for _, r := range rows {
		if r.Status == "new" {
			newCount++
		} else {
			modCount++
		}
	}

	Top(width)
	Row(Center("Files Added to Entry List/Index", width), width)
	Mid(width)
	Row(fmt.Sprintf("Total: %d new, %d modified", newCount, modCount), width)
	Mid(width)

	if len(rows) == 0 {
		Row("  No new or modified files.", width)
	} else {
		// header row
		Row(fmt.Sprintf("%-*s  %-*s  %-*s", fileW, "File", shaW, "SHA", statusW, "Status"), width)
		Mid(width)

		pageSize := 20
		for i, r := range rows {
			shortSHA := r.SHA
			if len(shortSHA) > shaW {
				shortSHA = shortSHA[:shaW]
			}
			fileName := r.File
			if len(fileName) > fileW {
				fileName = "..." + fileName[len(fileName)-fileW+3:]
			}
			Row(fmt.Sprintf("%-*s  %-*s  %-*s", fileW, fileName, shaW, shortSHA, statusW, r.Status), width)

			// paginate — only when --all and more than pageSize
			if isAll && (i+1)%pageSize == 0 && i+1 < len(rows) {
				Mid(width)
				Row(fmt.Sprintf("  showing %d/%d — press ENTER for more...", i+1, len(rows)), width)
				Mid(width)
				fmt.Scanln()
			}
		}
	}

	if len(ignored) > 0 {
		Mid(width)
		Row(fmt.Sprintf("Ignored: %d files (run with --show-ignored to see)", len(ignored)), width)
	}

	Bottom(width)
}

func PrintIgnoredFiles(ignored []string) {
	const width = 72
	Top(width)
	Row(Center("Ignored Files", width), width)
	Mid(width)
	for _, f := range ignored {
		Row("  "+f, width)
	}
	Mid(width)
	Row(fmt.Sprintf("Total: %d ignored", len(ignored)), width)
	Bottom(width)
}

func PrintRemoveTable(rows []RemoveRow, skipped []string) {
	const width = 72
	fileW := 45
	shaW := 15

	Top(width)
	Row(Center("Files Removed from Entry List/Index", width), width)
	Mid(width)
	Row(fmt.Sprintf("Total: %d removed, %d skipped", len(rows), len(skipped)), width)
	Mid(width)

	if len(rows) == 0 {
		Row("  No files removed.", width)
	} else {
		Row(fmt.Sprintf("%-*s  %-*s", fileW, "File", shaW, "SHA"), width)
		Mid(width)
		for _, r := range rows {
			shortSHA := r.SHA
			if len(shortSHA) > shaW {
				shortSHA = shortSHA[:shaW]
			}
			fileName := r.File
			if len(fileName) > fileW {
				fileName = "..." + fileName[len(fileName)-fileW+3:]
			}
			Row(fmt.Sprintf("%-*s  %-*s", fileW, fileName, shaW, shortSHA), width)
		}
	}

	if len(skipped) > 0 {
		Mid(width)
		Row(fmt.Sprintf("Skipped: %d files (not tracked or missing)", len(skipped)), width)
		for _, s := range skipped {
			Row("  "+s, width)
		}
	}
	Bottom(width)
}

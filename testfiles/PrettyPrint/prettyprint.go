package prettyprint

import (
	"fmt"
	"strings"
)

// default value *can mess around with this
const DefaultWidth = 65

// main building blocks
// -------------------------------------------------------------
func border(left, fill, right string, width int) string {
	return fmt.Sprintf("%s%s%s", left, strings.Repeat(fill, width), right)
}

func Top(width int) {
	fmt.Println(border("┌", "─", "┐", width))
}
func Mid(width int) {
	fmt.Println(border("├", "─", "┤", width))
}
func Bottom(width int) {
	fmt.Println(border("└", "─", "┘", width))
}

func Row(text string, width int) {
	inner := width - 4
	for _, line := range wrapText(text, inner) {
		fmt.Printf("│ %-*s   │\n", inner, line)
	}
}

func EmptyRow(width int) {
	Row("", width)
}

func Header(title string, width int) {
	Row(title, width)
	Mid(width)
}

//-------------------------------------------------------------

// Helper function to support main logic

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

func Center(text string, width int) string {
	inner := width - 4 // account for "│ " and " │"
	if len(text) >= inner {
		return text
	}
	padding := (inner - len(text)) / 2
	return strings.Repeat(" ", padding) + text
}

func SplitLines(text string) []string {
	var out []string
	current := ""
	for _, r := range text {
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

//-------------------------------------------------------------

// for printing objects content

func PrintObjectContent(sha string, content []byte) {
	const width = DefaultWidth
	Top(width)
	Row("Object: "+sha, width)
	Mid(width)
	EmptyRow(width)
	for _, line := range SplitLines(string(content)) {
		Row(line, width)
	}
	EmptyRow(width)
	Bottom(width)
}

//print commit info

func PrintCommit(sha, author, date, tree, parent, message string) {
	const width = 72 //slightly longer
	Top(width)
	Header("Version", width)
	Row("SHA    : "+sha, width)
	Row("Author : "+author, width)
	Row("Date   : "+date, width)
	EmptyRow(width)
	Row("Tree   : "+tree, width)
	if parent != "" {
		Row("Parent : "+parent, width)
	}
	Mid(width)
	Header("Message", width)
	EmptyRow(width)
	for _, line := range SplitLines(message) {
		Row(line, width)
	}
	Bottom(width)
}

//after storing object

func PrintObjectStored(objectType, fileName, sha string) {
	const width = 55
	Top(width)
	Row(Center("Object Stored Successfully", width), width)
	Mid(width)
	Row(fmt.Sprintf("Type : %s", objectType), width)
	Row(fmt.Sprintf("File : %s", fileName), width)
	EmptyRow(width)
	Row(fmt.Sprintf("SHA  : %s", sha), width)
	EmptyRow(width)
	Row("The object has been compressed and written", width)
	Row("to the repository object database.", width)
	Mid(width)
	Row("Use the SHA above to reference this object.", width)
	Bottom(width)
}

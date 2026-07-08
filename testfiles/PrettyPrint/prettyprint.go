package prettyprint

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// default value *can mess around with this
const DefaultWidth = 65

// rawMode flag for rawmode for better interactions
var rawMode bool

// rawMode related funcs so that format will not be messed up
func SetRawMode(enabled bool) {
	rawMode = enabled
}

var ansiEscapeRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func visibleLength(s string) int {
	stripped := ansiEscapeRegex.ReplaceAllString(s, "")
	return utf8.RuneCountInString(stripped)
}

func println(s string) {
	if rawMode {
		fmt.Print(s + "\r\n")
	} else {
		fmt.Println(s)
	}
}

func printf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	if rawMode {
		// replace any \n with \r\n
		s = strings.ReplaceAll(s, "\n", "\r\n")
		fmt.Print(s)
	} else {
		fmt.Print(s)
	}
}

// main building blocks
// -------------------------------------------------------------
func border(left, fill, right string, width int) string {
	return fmt.Sprintf("%s%s%s", left, strings.Repeat(fill, width), right)
}

func Top(width int)    { println(border("┌", "─", "┐", width)) }
func Mid(width int)    { println(border("├", "─", "┤", width)) }
func Bottom(width int) { println(border("└", "─", "┘", width)) }

func Row(s string, width int) {
	padding := width - visibleLength(s) - 2
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("│ %s%s │\r\n", s, strings.Repeat(" ", padding))
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
	for len(text) > width { //fit text to width
		lines = append(lines, text[:width])
		text = text[width:]
	}
	if text != "" {
		lines = append(lines, text)
	}
	return lines
}

func Center(text string, width int) string {
	inner := width - 4 // reserve for border
	if len(text) >= inner {
		return text
	}
	padding := (inner - len(text)) / 2
	return strings.Repeat(" ", padding) + text
}

// split strings with newlines into arrays
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
	const width = DefaultWidth + 5
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

func PrintObjectContentNoBorder(sha string, content []byte) {
	const width = 72
	separator := strings.Repeat("─", width)

	println("  Object: " + sha)
	println(separator)
	println("")

	for _, line := range SplitLines(string(content)) {
		println(line)
	}

	println("")
	println(separator)
}

//print commit info

func PrintCommit(sha, author, date, tree, parent, header, message string) {
	const width = 72 //slightly longer
	Top(width)
	Header("Version "+header, width)
	Row("SHA    : "+sha, width)
	Row("Author : "+author, width)
	Row("Date   : "+date, width)
	EmptyRow(width)
	Row("Tree SHA   : "+tree, width)
	if parent != "" {
		Row("Parent SHA : "+parent, width)
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

func PrintMessage(header string, path string, message string) {
	const width = 69

	Top(width)
	Row(Center(header, width), width)
	Mid(width)

	if path != "" {
		Row("Path : "+path, width)
		EmptyRow(width)
	}

	for _, line := range SplitLines(message) {
		Row(line, width)
	}

	Bottom(width)
}

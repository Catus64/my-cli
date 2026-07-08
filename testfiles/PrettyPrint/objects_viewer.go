package prettyprint

import (
	"bufio"
	"fmt"
	diffengine "gocmd/testfiles/diffEngine"
	"os"
	"strings"

	"golang.org/x/term"
)

type FileChangeStatus string

const (
	StatusAdded    FileChangeStatus = "added"
	StatusModified FileChangeStatus = "modified"
	StatusDeleted  FileChangeStatus = "deleted"
)

const (
	colorRed   = "\x1b[31m"
	colorGreen = "\x1b[32m"
	colorReset = "\x1b[0m"
)

const separatorLine = "────────────────────────────────────────────────────────────────"

type FileChange struct {
	Path   string
	Status FileChangeStatus
	OldSHA string // empty if added
	NewSHA string // empty if deleted
}

type objectsViewerState struct {
	changes []FileChange
	cursor  int
	config  ViewerConfig
	header  string
}

func printRawLine(s string) {
	fmt.Print(s + "\r\n")
}

func drawObjectsTable(state *objectsViewerState) {
	clearScreen()
	const width = 72

	if state.header != "" {
		Top(width)
		Row(Center(state.header, width), width)
		Mid(width)
	} else {
		Top(width)
	}
	Row(Center("Changed Files", width), width)
	Mid(width)

	if len(state.changes) == 0 {
		Row("  Nothing to show.", width)
		Bottom(width)
		return
	}

	fileW := 45
	statusW := 12
	Row(fmt.Sprintf("%-*s  %-*s", fileW, "File", statusW, "Status"), width)
	Mid(width)

	for i, c := range state.changes {
		pointer := "  "
		if i == state.cursor {
			pointer = "> "
		}
		fileName := c.Path
		if len(fileName) > fileW {
			fileName = "..." + fileName[len(fileName)-fileW+3:]
		}
		statusLabel := strings.ToUpper(string(c.Status))
		Row(fmt.Sprintf("%s%-*s  %-*s", pointer, fileW, fileName, statusW, statusLabel), width)
	}

	Mid(width)
	Row("  j/k or arrows: move   Enter: view   q: back", width)
	Bottom(width)
}

func RunObjectsViewer(
	changes []FileChange,
	fetchBlob func(sha string) ([]byte, error),
	config ViewerConfig,
	header string,
) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to enter raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	SetRawMode(true)
	defer SetRawMode(false)

	state := &objectsViewerState{
		changes: changes,
		cursor:  0,
		config:  config,
		header:  header,
	}

	drawObjectsTable(state)
	reader := bufio.NewReader(os.Stdin)

	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}

		switch {
		case b == 'q' || b == 3:
			term.Restore(int(os.Stdin.Fd()), oldState)
			clearScreen()
			return nil

		case b == '\x1b':
			b2, _ := reader.ReadByte()
			b3, _ := reader.ReadByte()
			if b2 != '[' {
				break
			}
			switch b3 {
			case 'A': // up
				if state.cursor > 0 {
					state.cursor--
				}
				drawObjectsTable(state)
			case 'B': // down
				if state.cursor < len(state.changes)-1 {
					state.cursor++
				}
				drawObjectsTable(state)
			}

		case b == 'k':
			if state.cursor > 0 {
				state.cursor--
			}
			drawObjectsTable(state)

		case b == 'j':
			if state.cursor < len(state.changes)-1 {
				state.cursor++
			}
			drawObjectsTable(state)
		case b == '\r' || b == '\n':
			selected := state.changes[state.cursor]
			clearScreen()
			if err := viewFileChangeContent(reader, selected, fetchBlob); err != nil {
				return err
			}
			drawObjectsTable(state)
		}
	}
	return nil
}

// viewFileChangeContent shows blob content for a single FileChange.
// For added/deleted files, there's only one side — no toggling.
// For modified files, left/right arrows toggle between old and new.
func viewFileChangeContent(reader *bufio.Reader, change FileChange, fetchBlob func(string) ([]byte, error)) error {
	showingNew := true // default to showing "new" side when both exist

	var oldLines, newLines []string
	var oldKept, newKept []bool
	diffComputed := false

	// Only modified files need the diff computed (added/deleted have nothing to compare against)
	ensureDiff := func() error {
		if diffComputed || change.Status != StatusModified {
			return nil
		}
		oldContent, err := fetchBlob(change.OldSHA)
		if err != nil {
			return err
		}
		newContent, err := fetchBlob(change.NewSHA)
		if err != nil {
			return err
		}
		oldLines = strings.Split(string(oldContent), "\n")
		newLines = strings.Split(string(newContent), "\n")
		oldKept, newKept = diffengine.ComputeLineDiff(oldLines, newLines)
		diffComputed = true
		return nil
	}

	render := func() error {
		clearScreen()

		switch change.Status {
		case StatusAdded:
			content, err := fetchBlob(change.NewSHA)
			if err != nil {
				fmt.Print("  error reading blob: " + err.Error() + "\r\n")
				return nil
			}
			fmt.Print("  " + change.Path + " — AFTER (added)\r\n")
			fmt.Print("  " + separatorLine + "\r\n")
			for _, line := range strings.Split(string(content), "\n") {
				fmt.Print(colorGreen + line + colorReset + "\r\n")
			}
			fmt.Print("  " + separatorLine + "\r\n")
			fmt.Print("\r\n  press any key to go back\r\n")

		case StatusDeleted:
			content, err := fetchBlob(change.OldSHA)
			if err != nil {
				fmt.Print("  error reading blob: " + err.Error() + "\r\n")
				return nil
			}
			fmt.Print("  " + change.Path + " — BEFORE (deleted)\r\n")
			fmt.Print("  " + separatorLine + "\r\n")
			for _, line := range strings.Split(string(content), "\n") {
				fmt.Print(colorRed + line + colorReset + "\r\n")
			}
			fmt.Print("  " + separatorLine + "\r\n")
			fmt.Print("\r\n  press any key to go back\r\n")

		case StatusModified:
			if err := ensureDiff(); err != nil {
				fmt.Print("  error computing diff: " + err.Error() + "\r\n")
				return nil
			}
			if showingNew {
				fmt.Print("  " + change.Path + " — AFTER\r\n")
				fmt.Print("  " + separatorLine + "\r\n\r\n")
				renderDiffLines(newLines, newKept, false)
			} else {
				fmt.Print("  " + change.Path + " — BEFORE\r\n")
				fmt.Print("  " + separatorLine + "\r\n\r\n")
				renderDiffLines(oldLines, oldKept, true)
			}
			fmt.Print("\r\n  " + separatorLine + "\r\n")
			fmt.Print("  (end of file)\r\n")
			fmt.Print("\r\n  ← / → to toggle before/after, q to go back\r\n")
		}
		return nil
	}

	if err := render(); err != nil {
		return err
	}

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil
		}
		if b == 'q' || b == 3 {
			return nil
		}
		if change.Status != StatusModified {
			return nil // any key goes back for added/deleted
		}
		if b == '\x1b' {
			b2, _ := reader.ReadByte()
			b3, _ := reader.ReadByte()
			if b2 == '[' {
				switch b3 {
				case 'C', 'D': // left or right — toggle, direction doesn't matter with only 2 states
					showingNew = !showingNew
					if err := render(); err != nil {
						return err
					}
				}
			}
		}
	}
}

func renderDiffLines(lines []string, kept []bool, isOld bool) {
	expanded := make([]string, len(lines))
	for i, l := range lines {
		expanded[i] = expandTabs(l)
	}

	width := maxLineWidth(expanded) + 5

	Top(width)
	for i, line := range expanded {
		color := colorGreen
		if isOld {
			color = colorRed
		}
		if kept[i] {
			printWrappedRow(line, "", width)
		} else {
			printWrappedRow(line, color, width)
		}
	}
	Bottom(width)
}

// printWrappedRow prints a line inside the box border, wrapping it across
// multiple rows if it's longer than the box's content width, and applying
// color without corrupting the width/padding calculation.
func printWrappedRow(line string, color string, width int) {
	innerWidth := width - 2 // account for the box's own left/right border chars, adjust if your Row already does this

	if line == "" {
		EmptyRow(width)
		return
	}

	for len(line) > innerWidth {
		chunk := line[:innerWidth]
		line = line[innerWidth:]
		if color != "" {
			Row(color+chunk+colorReset, width)
		} else {
			Row(chunk, width)
		}
	}
	if color != "" {
		Row(color+line+colorReset, width)
	} else {
		Row(line, width)
	}
}

func maxLineWidth(lines []string) int {
	max := DefaultWidth
	for _, line := range lines {
		l := len(expandTabs(line))
		if l > max {
			max = l
		}
	}
	return max
}

func expandTabs(line string) string {
	const tabWidth = 4 // pick whatever your code style uses
	return strings.ReplaceAll(line, "\t", strings.Repeat(" ", tabWidth))
}

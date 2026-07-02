package prettyprint

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/term"
)

type SavelistEntry struct {
	File string
	Type string // "blob", "tree", "exec", "symlink"
	SHA  string
}

type savelistState struct {
	entries     []SavelistEntry
	currentPage int
	totalPages  int
	pageSize    int
	config      ViewerConfig
}

func ModeToType(modeType uint16, modePerms uint16) string {
	switch modeType {
	case 0b1000: // regular file
		if modePerms&0111 != 0 {
			return "exec"
		}
		return "blob"
	case 0b1010:
		return "symlink"
	case 0b1110:
		return "gitlink"
	case 0b0100:
		return "tree"
	default:
		return "blob"
	}
}

func drawSavelist(state *savelistState) {
	if state.config.ClearOnRedraw {
		clearScreen()
	}

	const width = 72
	fileW := 38
	typeW := 8
	shaW := 10

	start := state.currentPage * state.pageSize
	end := start + state.pageSize
	if end > len(state.entries) {
		end = len(state.entries)
	}
	pageEntries := state.entries[start:end]

	Top(width)
	Row(Center("Entry List", width), width)
	Mid(width)
	Row(fmt.Sprintf("Page %d / %d   (%d files total)",
		state.currentPage+1, state.totalPages, len(state.entries)), width)
	Mid(width)
	Row(fmt.Sprintf("%-*s  %-*s  %-*s", fileW, "File", typeW, "Type", shaW, "SHA"), width)
	Mid(width)

	for _, e := range pageEntries {
		fileName := e.File
		if len(fileName) > fileW {
			fileName = "..." + fileName[len(fileName)-fileW+3:]
		}
		shortSHA := e.SHA
		if len(shortSHA) > shaW {
			shortSHA = shortSHA[:shaW]
		}
		typeStr := e.Type
		if len(typeStr) > typeW {
			typeStr = typeStr[:typeW]
		}
		Row(fmt.Sprintf("%-*s  %-*s  %-*s", fileW, fileName, typeW, typeStr, shaW, shortSHA), width)
	}

	// pad empty rows if last page has fewer entries
	for i := len(pageEntries); i < state.pageSize; i++ {
		EmptyRow(width)
	}

	Mid(width)
	Row("← prev page   → next page   q quit", width)
	Bottom(width)
}

func RunSavelistViewer(entries []SavelistEntry, config ViewerConfig) error {
	if len(entries) == 0 {
		const width = 72
		Top(width)
		Row(Center("Entry List is empty", width), width)
		Row("No files have been added yet. Run ezgit add to stage files.", width)
		Bottom(width)
		return nil
	}

	const pageSize = 20
	totalPages := (len(entries) + pageSize - 1) / pageSize

	state := &savelistState{
		entries:     entries,
		currentPage: 0,
		totalPages:  totalPages,
		pageSize:    pageSize,
		config:      config,
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to enter raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	SetRawMode(true)
	defer SetRawMode(false)

	drawSavelist(state)

	reader := bufio.NewReader(os.Stdin)
	_ = strconv.Atoi // keep import

	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}

		switch {
		case b == 'q' || b == 3:
			term.Restore(int(os.Stdin.Fd()), oldState)
			SetRawMode(false)
			clearScreen()
			return nil

		case b == '\x1b':
			b2, _ := reader.ReadByte()
			b3, _ := reader.ReadByte()
			if b2 != '[' {
				break
			}
			switch b3 {
			case 'D': // left arrow — prev page
				if state.currentPage > 0 {
					state.currentPage--
					drawSavelist(state)
				}
			case 'C': // right arrow — next page
				if state.currentPage < state.totalPages-1 {
					state.currentPage++
					drawSavelist(state)
				}
			}
		}
	}
	return nil
}

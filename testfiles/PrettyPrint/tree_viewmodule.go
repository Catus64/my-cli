package prettyprint

import (
	"bufio"
	"fmt"
	logger "gocmd/testfiles/Helper"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// if user prefer not redraw
type ViewerConfig struct {
	ClearOnRedraw bool
}

// Default for now
var DefaultViewerConfig = ViewerConfig{
	ClearOnRedraw: true,
}

type TreeItem struct {
	Name   string
	SHA    string
	IsTree bool
}

// context switch to view trees info
// breadcump means what directory you're at
type treeViewerState struct {
	items      []TreeItem
	cursor     int
	breadcrumb []string
	config     ViewerConfig
	header     string // optional when called with version SHA
}

type stackEntry struct {
	items      []TreeItem
	cursor     int
	breadcrumb []string
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func drawTree(state *treeViewerState) {
	if state.config.ClearOnRedraw {
		clearScreen()
	}

	const width = 72
	fileW := 45
	shaW := 10

	Top(width)

	// optional version header
	if state.header != "" {
		Row(state.header, width)
		Mid(width)
	}

	// breadcrumb
	crumb := "root"
	if len(state.breadcrumb) > 0 {
		crumb = "root > " + strings.Join(state.breadcrumb, " > ")
	}
	Row(crumb, width)
	Mid(width)
	Row(fmt.Sprintf("  %-4s  %-*s  %-*s", "#", fileW, "Name", shaW, "SHA"), width)
	Mid(width)

	for i, item := range state.items {
		cursor := " "
		if i == state.cursor {
			cursor = ">"
		}
		name := item.Name
		if item.IsTree {
			name = "[" + name + "]"
		}
		if len(name) > fileW {
			name = "..." + name[len(name)-fileW+3:]
		}
		shortSHA := item.SHA
		if len(shortSHA) > shaW {
			shortSHA = shortSHA[:shaW]
		}
		Row(fmt.Sprintf("%s %-4d  %-*s  %-*s", cursor, i+1, fileW, name, shaW, shortSHA), width)
	}

	Mid(width)
	Row("arrows/j/k: navigate   number+enter: jump   enter: select   b: back   q: quit", width)
	Bottom(width)
}

// RunTreeViewer starts the interactive tree viewer
// fetchChildren: given a tree SHA, return its items
// fetchBlob: given a blob SHA, return its content
func RunTreeViewer(
	rootItems []TreeItem,
	fetchChildren func(sha string) ([]TreeItem, error),
	fetchBlob func(sha string) ([]byte, error),
	config ViewerConfig,
	header string, // empty string = no header
) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to enter raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	SetRawMode(true)
	defer SetRawMode(false)

	state := &treeViewerState{
		items:      rootItems,
		cursor:     0,
		breadcrumb: []string{},
		config:     config,
		header:     header,
	}

	// stack for navigation history
	var stack []stackEntry

	drawTree(state)

	// number input buffer
	numBuf := ""
	reader := bufio.NewReader(os.Stdin)

	for { //input loop
		b, err := reader.ReadByte()
		if err != nil {
			break
		}

		switch {
		// quit
		case b == 'q' || b == 3: // 3 is for ctrl-c input
			term.Restore(int(os.Stdin.Fd()), oldState)
			clearScreen()
			return nil

		// arrow keys — escape sequence \x1b[A / \x1b[B
		case b == '\x1b':
			b2, _ := reader.ReadByte()
			b3, _ := reader.ReadByte()
			if b2 == '[' {
				switch b3 {
				case 'A': // up
					numBuf = ""
					if state.cursor > 0 {
						state.cursor--
					}
					drawTree(state)
				case 'B': // down
					numBuf = ""
					if state.cursor < len(state.items)-1 {
						state.cursor++
					}
					drawTree(state)
				}
			}

		// vim keys
		case b == 'k': // up
			numBuf = ""
			if state.cursor > 0 {
				state.cursor--
			}
			drawTree(state)

		case b == 'j': // down
			numBuf = ""
			if state.cursor < len(state.items)-1 {
				state.cursor++
			}
			drawTree(state)

		// number input
		case b >= '0' && b <= '9':
			numBuf += string(b)
			// show the number being typed in the box
			drawTree(state)
			// reprint bottom with number hint
			fmt.Printf("\r  jumping to: %s", numBuf)

		// enter, select or jump to number
		case b == '\r' || b == '\n':
			if numBuf != "" {
				// jump to number
				num, err := strconv.Atoi(numBuf)
				numBuf = ""
				if err == nil && num >= 1 && num <= len(state.items) {
					state.cursor = num - 1
				}
				drawTree(state)
				break
			}

			// select current item
			selected := state.items[state.cursor]

			if selected.IsTree {
				// push current state to stack and navigate in
				stack = append(stack, stackEntry{
					items:      state.items,
					cursor:     state.cursor,
					breadcrumb: state.breadcrumb,
				})
				children, err := fetchChildren(selected.SHA)
				if err != nil {
					// show error then redraw
					term.Restore(int(os.Stdin.Fd()), oldState)
					fmt.Printf("error reading tree: %v\n", err)
					oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
					drawTree(state)
					break
				}
				state.items = children
				state.cursor = 0
				state.breadcrumb = append(append([]string{}, state.breadcrumb...), selected.Name)
				drawTree(state)

			} else {
				// blob show content
				term.Restore(int(os.Stdin.Fd()), oldState)
				clearScreen()

				content, err := fetchBlob(selected.SHA)
				if err != nil {
					fmt.Printf("error reading blob: %v\n", err)
				} else {
					PrintObjectContentNoBorder(selected.SHA, content)
				}

				fmt.Println("\npress q or ctrl+c to quit, any other key to go back...")
				// wait for keypress
				oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
				key, _ := reader.ReadByte()
				if key == 'q' || key == 3 {
					term.Restore(int(os.Stdin.Fd()), oldState)
					clearScreen()
					return nil
				}
				drawTree(state)
			}

		// back
		case b == 'b':
			numBuf = ""
			if len(stack) > 0 {
				prev := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				state.items = prev.items
				state.cursor = prev.cursor
				state.breadcrumb = prev.breadcrumb
				drawTree(state)
			}
		}
	}

	return nil
}

func RunRefsViewer(
	items []TreeItem,
	fetchChildren func(sha string) ([]TreeItem, error),
	onSelect func(sha string) error,
	config ViewerConfig,
	header string,
) error {
	logger.L().Debug("running refs viewer")
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to enter raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	SetRawMode(true)
	defer SetRawMode(false)

	state := &treeViewerState{
		items:      items,
		cursor:     0,
		breadcrumb: []string{},
		config:     config,
		header:     header,
	}

	drawTree(state)

	numBuf := ""
	reader := bufio.NewReader(os.Stdin)

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
			if b2 == '[' {
				switch b3 {
				case 'A':
					numBuf = ""
					if state.cursor > 0 {
						state.cursor--
					}
					drawTree(state)
				case 'B':
					numBuf = ""
					if state.cursor < len(state.items)-1 {
						state.cursor++
					}
					drawTree(state)
				}
			}

		case b == 'k':
			numBuf = ""
			if state.cursor > 0 {
				state.cursor--
			}
			drawTree(state)

		case b == 'j':
			numBuf = ""
			if state.cursor < len(state.items)-1 {
				state.cursor++
			}
			drawTree(state)

		case b >= '0' && b <= '9':
			numBuf += string(b)
			drawTree(state)
			fmt.Printf("\r  jumping to: %s", numBuf)

		case b == '\r' || b == '\n':
			if numBuf != "" {
				num, err := strconv.Atoi(numBuf)
				numBuf = ""
				if err == nil && num >= 1 && num <= len(state.items) {
					state.cursor = num - 1
				}
				drawTree(state)
				break
			}

			selected := state.items[state.cursor]

			// restore terminal before handing off to commit view
			term.Restore(int(os.Stdin.Fd()), oldState)
			SetRawMode(false)
			clearScreen()

			if err := onSelect(selected.SHA); err != nil {
				fmt.Printf("error: %v\n", err)
			}

			// re-enter raw mode and redraw after returning
			oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
			SetRawMode(true)
			drawTree(state)
		}
	}
	return nil
}

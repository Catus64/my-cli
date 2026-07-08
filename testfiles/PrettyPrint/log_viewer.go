package prettyprint

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

type LogCommit struct {
	SHA         string
	ShortSHA    string
	Author      string
	Date        string
	Message     string
	TreeSHA     string
	Parents     []string
	VersionNum  int
	VersionName string
	HasVersion  bool
}

type logViewerState struct {
	current      LogCommit
	history      []LogCommit // stack of visited commits (for right arrow)
	parentCursor int         // for merge commits — which parent is selected
	config       ViewerConfig
}

func drawLogCommit(state *logViewerState) {
	if state.config.ClearOnRedraw {
		// clearScreen()
	}

	const width = 72
	commit := state.current

	Top(width)

	// version header if available
	if commit.HasVersion {
		Row(Center(fmt.Sprintf("v%d · %s", commit.VersionNum, commit.VersionName), width), width)
		Mid(width)
	}

	// commit info
	Row(fmt.Sprintf("SHA      : %s", commit.ShortSHA), width)
	Row(fmt.Sprintf("Author   : %s", commit.Author), width)
	Row(fmt.Sprintf("Date     : %s", commit.Date), width)
	EmptyRow(width)
	Row(fmt.Sprintf("Tree   SHA: %s", shortSHA(commit.TreeSHA)), width)

	// parents section
	if len(commit.Parents) == 0 {
		Row("Parent   : none (first commit)", width)
	} else if len(commit.Parents) == 1 {
		Row(fmt.Sprintf("Parent SHA: %s", shortSHA(commit.Parents[0])), width)
	} else {
		// merge commit — show both with cursor
		Mid(width)
		Row("Merge parents - use up/down to select, left to follow:", width)
		Mid(width)
		for i, p := range commit.Parents {
			cursor := "  "
			if i == state.parentCursor {
				cursor = "> "
			}
			Row(fmt.Sprintf("%s[%d] %s", cursor, i+1, shortSHA(p)), width)
		}
	}

	EmptyRow(width)

	// message
	Mid(width)
	Row("Message", width)
	Mid(width)
	EmptyRow(width)
	for _, line := range SplitLines(commit.Message) {
		if line != "" {
			printWrappedMessageRow("  "+line, width)
		}
	}
	EmptyRow(width)

	// nav footer
	Mid(width)
	if len(commit.Parents) > 1 {
		Row("← older (selected parent)   → newer   v view tree   q quit", width)
	} else if len(commit.Parents) == 1 {
		Row("← older   → newer   v view tree   q quit   c show changes", width)
	} else {
		Row("→ newer   v view tree   q quit", width)
	}
	Bottom(width)
}

func printWrappedMessageRow(s string, width int) {
	const borderOverhead = 4 // adjust to match however much Row()'s own "│ " + " │" eats up

	innerWidth := width - borderOverhead
	if innerWidth <= 0 {
		Row(s, width) // fallback, box too narrow to wrap meaningfully
		return
	}

	runes := []rune(s)
	for len(runes) > innerWidth {
		chunk := string(runes[:innerWidth])
		Row(chunk, width)
		runes = runes[innerWidth:]
	}
	Row(string(runes), width)
}

func drawEndOfHistory(config ViewerConfig) {
	if config.ClearOnRedraw {
		// clearScreen()
	}
	const width = 72
	Top(width)
	Row(Center("End of History", width), width)
	Mid(width)
	EmptyRow(width)
	Row(Center("No more commits to show.", width), width)
	EmptyRow(width)
	Mid(width)
	Row("→ newer   q quit", width)
	Bottom(width)
}

func shortSHA(sha string) string {
	if len(sha) >= 7 {
		return sha[:7]
	}
	return sha
}

// RunLogViewer starts the interactive log viewer
// fetchCommit: given a SHA returns a LogCommit
// fetchBlob/fetchChildren: passed through to tree viewer
func RunLogViewer(
	initial LogCommit,
	fetchCommit func(sha string) (LogCommit, error),
	fetchTreeItems func(sha string) ([]TreeItem, error),
	fetchBlob func(sha string) ([]byte, error),
	fetchChanges func(currentTreeSHA, parentTreeSHA string) ([]FileChange, error),
	config ViewerConfig,
) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to enter raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	SetRawMode(true)
	defer SetRawMode(false)

	state := &logViewerState{
		current:      initial,
		history:      []LogCommit{},
		parentCursor: 0,
		config:       config,
	}

	atEnd := false // true when we've gone past first commit
	drawLogCommit(state)

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
			// clearScreen()
			return nil

		case b == '\x1b':
			b2, _ := reader.ReadByte()
			b3, _ := reader.ReadByte()
			if b2 != '[' {
				break
			}
			switch b3 {
			case 'D': // left arrow — go older
				if atEnd {
					break
				}
				parents := state.current.Parents
				if len(parents) == 0 {
					// first commit next page will show end of history
					atEnd = true
					// still add the first commit to stack
					state.history = append(state.history, state.current)
					drawEndOfHistory(state.config)
					break
				}
				// pick parent based on cursor
				targetSHA := parents[0]
				if len(parents) > 1 {
					targetSHA = parents[state.parentCursor]
				}
				next, err := fetchCommit(targetSHA)
				if err != nil {
					break
				}
				state.history = append(state.history, state.current)
				state.current = next
				state.parentCursor = 0
				atEnd = false
				drawLogCommit(state)

			case 'C': // right arrow — go newer
				if len(state.history) == 0 {
					break // already at newest
				}
				prev := state.history[len(state.history)-1]
				state.history = state.history[:len(state.history)-1]
				state.current = prev
				state.parentCursor = 0
				atEnd = false
				drawLogCommit(state)

			case 'A': // up arrow — select previous parent (merge commits)
				if len(state.current.Parents) > 1 {
					if state.parentCursor > 0 {
						state.parentCursor--
					}
					drawLogCommit(state)
				}

			case 'B': // down arrow — select next parent (merge commits)
				if len(state.current.Parents) > 1 {
					if state.parentCursor < len(state.current.Parents)-1 {
						state.parentCursor++
					}
					drawLogCommit(state)
				}
			}

		case b == 'v': // view tree
			treeSHA := state.current.TreeSHA
			header := ""
			if state.current.HasVersion {
				header = fmt.Sprintf("v%d · %s", state.current.VersionNum, state.current.VersionName)
			}

			term.Restore(int(os.Stdin.Fd()), oldState)
			SetRawMode(false)
			clearScreen()

			items, err := fetchTreeItems(treeSHA)
			if err != nil {
				fmt.Printf("error reading tree: %v\n", err)
			} else {
				RunTreeViewer(items, fetchTreeItems, fetchBlob, config, header)
			}

			// re-enter raw mode after tree viewer exits
			oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
			SetRawMode(true)
			if atEnd {
				drawEndOfHistory(state.config)
			} else {
				drawLogCommit(state)
			}
		case b == 'c': // view changed files vs parent
			var parentTreeSHA string
			noParent := len(state.current.Parents) == 0

			if !noParent {
				targetParentSHA := state.current.Parents[0]
				if len(state.current.Parents) > 1 {
					targetParentSHA = state.current.Parents[state.parentCursor]
				}
				parentCommit, err := fetchCommit(targetParentSHA)
				if err == nil {
					parentTreeSHA = parentCommit.TreeSHA
				} else {
					noParent = true // treat resolution failure as "nothing to compare"
				}
			}

			term.Restore(int(os.Stdin.Fd()), oldState)
			SetRawMode(false)
			clearScreen()

			if noParent {
				fmt.Println("\n  Nothing to show. (root commit has no parent)")
				fmt.Println("\n  press any key to go back...")
				reader.ReadByte()
			} else {
				changes, err := fetchChanges(state.current.TreeSHA, parentTreeSHA)
				if err != nil {
					fmt.Printf("error computing changes: %v\n", err)
					fmt.Println("\npress any key to go back...")
					reader.ReadByte()
				} else {
					header := ""
					if state.current.HasVersion {
						header = fmt.Sprintf("v%d · %s", state.current.VersionNum, state.current.VersionName)
					}
					RunObjectsViewer(changes, fetchBlob, config, header)
				}
			}

			oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
			SetRawMode(true)
			if atEnd {
				drawEndOfHistory(state.config)
			} else {
				drawLogCommit(state)
			}
		}
	}
	return nil
}

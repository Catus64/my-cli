package gitcheckignore

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"

	setconfig "gocmd/testfiles/Gitrepostruct"
)

// ---------- Types ----------

type mode int

const (
	modeNormal mode = iota
	modeAction      // after Enter on a rule: showing edit/delete options
	modeEdit        // typing a new/edited rule
)

const ctrlC = 0x03
const boxWidth = 70

var errCancelled = errors.New("cancelled")

type editorState struct {
	rules    []string
	cursor   int
	mode     mode
	input    string // buffer for add/edit typing
	editing  bool   // true if input buffer is editing existing rule vs adding new
	repoPath string
}

// ---------- Entry point ----------

// EditGitignore is the entry point called from the cobra command.
// repoRoot should be the repository's working tree root (e.g. repo.WorkTree).
func EditGitignore(repoRoot string) error {
	path := filepath.Join(repoRoot, ".gitignore")

	rules, err := loadOrCreate(path)
	if err != nil {
		return err
	}

	state := &editorState{
		rules:    rules,
		cursor:   0,
		mode:     modeNormal,
		repoPath: path,
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to enter raw terminal mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	err = runLoop(state)
	fmt.Print("\x1b[2J\x1b[H") // clear screen on exit either way

	if errors.Is(err, errCancelled) {
		fmt.Println("Cancelled — no changes saved.")
		return nil
	}
	if err != nil {
		return err
	}

	return saveRules(path, state.rules)
}

// ---------- Loading / Saving ----------

func loadOrCreate(path string) ([]string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaults, err := setconfig.LoadIgnoreDefaults()
		if err != nil {
			return nil, fmt.Errorf("failed to load default ignore template: %w", err)
		}
		if err := saveRules(path, defaults); err != nil {
			return nil, fmt.Errorf("failed to create .gitignore: %w", err)
		}
		return defaults, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to check .gitignore: %w", err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var rules []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		rules = append(rules, line)
	}
	return rules, scanner.Err()
}

func saveRules(path string, rules []string) error {
	var sb strings.Builder
	sb.WriteString("# Managed by ezgit - edit with 'ezgit ignore'\n")
	for _, r := range rules {
		sb.WriteString(r)
		sb.WriteString("\n")
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}

// ---------- Main loop ----------

func runLoop(state *editorState) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		render(state)

		b, err := reader.ReadByte()
		if err != nil {
			return err
		}

		if b == ctrlC {
			return errCancelled
		}

		switch state.mode {
		case modeNormal:
			if handleNormalKey(state, b, reader) {
				return nil // user quit & save
			}
		case modeAction:
			handleActionKey(state, b)
		case modeEdit:
			if handleEditKey(state, b) {
				commitInput(state)
			}
		}
	}
}

// handleNormalKey returns true if the user requested quit (save & exit).
func handleNormalKey(state *editorState, b byte, reader *bufio.Reader) bool {
	switch b {
	case 'q':
		return true
	case 'a':
		state.mode = modeEdit
		state.editing = false
		state.input = ""
	case 'j':
		moveCursor(state, 1)
	case 'k':
		moveCursor(state, -1)
	case '\r', '\n':
		if len(state.rules) > 0 {
			state.mode = modeAction
		}
	case 0x1b: // ESC — could be an arrow key escape sequence
		next1, err := reader.ReadByte()
		if err != nil || next1 != '[' {
			return false
		}
		next2, err := reader.ReadByte()
		if err != nil {
			return false
		}
		switch next2 {
		case 'A': // up
			moveCursor(state, -1)
		case 'B': // down
			moveCursor(state, 1)
		}
	}
	return false
}

func moveCursor(state *editorState, delta int) {
	if len(state.rules) == 0 {
		return
	}
	state.cursor += delta
	if state.cursor < 0 {
		state.cursor = 0
	}
	if state.cursor >= len(state.rules) {
		state.cursor = len(state.rules) - 1
	}
}

func handleActionKey(state *editorState, b byte) {
	switch b {
	case 'e':
		state.mode = modeEdit
		state.editing = true
		state.input = state.rules[state.cursor]
	case 'd':
		state.rules = append(state.rules[:state.cursor], state.rules[state.cursor+1:]...)
		if state.cursor >= len(state.rules) && state.cursor > 0 {
			state.cursor--
		}
		state.mode = modeNormal
	case 0x1b: // ESC cancels back to normal mode
		state.mode = modeNormal
	}
}

// handleEditKey returns true when Enter is pressed (input complete).
func handleEditKey(state *editorState, b byte) bool {
	switch b {
	case '\r', '\n':
		return true
	case 0x7f, 0x08: // backspace
		if len(state.input) > 0 {
			state.input = state.input[:len(state.input)-1]
		}
	case 0x1b: // ESC cancels
		state.mode = modeNormal
		state.input = ""
	default:
		if b >= 32 && b < 127 { // printable ASCII only
			state.input += string(b)
		}
	}
	return false
}

func commitInput(state *editorState) {
	trimmed := strings.TrimSpace(state.input)
	if trimmed != "" {
		if state.editing {
			state.rules[state.cursor] = trimmed
		} else {
			state.rules = append(state.rules, trimmed)
			state.cursor = len(state.rules) - 1
		}
	}
	state.input = ""
	state.mode = modeNormal
}

// ---------- Rendering ----------

func render(state *editorState) {
	fmt.Print("\x1b[2J\x1b[H")

	top(boxWidth)
	textRow(centerText(".gitignore Editor", boxWidth-1), boxWidth)
	mid(boxWidth)

	legend := []string{
		"dir/        -> ignore an entire directory",
		"*.ext       -> ignore all files with this extension",
		"file.txt    -> ignore one specific file",
		"!keep.txt   -> negate/un-ignore a previously ignored file",
		"**/name     -> match 'name' at any depth",
	}
	for _, l := range legend {
		textRow("  "+l, boxWidth)
	}
	mid(boxWidth)

	textRow("  j/k or arrows: move   a: add   Enter: select", boxWidth)
	textRow("  q: save & quit   ctrl-c: cancel", boxWidth)
	mid(boxWidth)

	if len(state.rules) == 0 {
		textRow("  (no rules yet — press 'a' to add one)", boxWidth)
	}
	for i, rule := range state.rules {
		line := "  " + rule
		if i == state.cursor && state.mode != modeEdit {
			line = "> " + rule
		}
		textRow(line, boxWidth)
	}

	switch state.mode {
	case modeAction:
		mid(boxWidth)
		textRow("  [e] edit    [d] delete    [esc] cancel", boxWidth)
	case modeEdit:
		label := "Add rule"
		if state.editing {
			label = "Edit rule"
		}
		mid(boxWidth)
		textRow(fmt.Sprintf("  %s: %s_", label, state.input), boxWidth)
	}

	bottom(boxWidth)
}

func top(w int) {
	fmt.Print("  ┌" + strings.Repeat("─", w) + "┐\r\n")
}

func mid(w int) {
	fmt.Print("  ├" + strings.Repeat("─", w) + "┤\r\n")
}

func bottom(w int) {
	fmt.Print("  └" + strings.Repeat("─", w) + "┘\r\n")
}

func textRow(s string, w int) {
	if len(s) > w-1 {
		s = s[:w-4] + "..."
	}
	padding := w - len(s)
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("  │%s%s│\r\n", s, strings.Repeat(" ", padding))
}

func centerText(s string, w int) string {
	if len(s) >= w {
		return s
	}
	left := (w - len(s)) / 2
	right := w - len(s) - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

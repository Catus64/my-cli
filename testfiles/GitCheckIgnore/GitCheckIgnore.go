package gitcheckignore

import "strings"

func Gitignore_Parse1(raw string) (pattern string, include bool, ok bool) {

	// include: whether the pattern is an include pattern (starts with '!')
	// ok: whether the line is a valid pattern (not empty and not a comment)

	raw = strings.TrimSpace(raw)

	if raw == "" || raw[0] == '#' {
		return "", false, false
	}

	if raw[0] == '!' {
		return raw[1:], false, true
	}

	if raw[0] == '\\' {
		return raw[1:], true, true
	}

	return raw, true, true
}

func Gitignore_Parse(lines []string) []string {
	patterns := make([]string, 0)

	for _, line := range lines {
		pattern, include, ok := Gitignore_Parse1(line)
		if ok && include {
			patterns = append(patterns, pattern)
		}
	}
	return patterns

}

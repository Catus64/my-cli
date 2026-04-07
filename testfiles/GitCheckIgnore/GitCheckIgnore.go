package gitcheckignore

import (
	"errors"
	logger "gocmd/testfiles/Helper"
	"path/filepath"
	"strings"
)

func Gitignore_Parse1(raw string) (pattern string, include bool, ok bool) {

	// include: whether the pattern is an include pattern (starts with '!')
	// (true means ignore this file, false means ignore this rule)
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

// parse multiple string lines into patterns, ignore invalid lines
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

func Check_Ignore_1(patterns []string, path string) *bool {
	var result *bool

	for _, p := range patterns {
		isNegation := false
		pattern := p

		if strings.HasPrefix(p, "!") {
			isNegation = true
			pattern = p[1:]
		}

		match, err := filepath.Match(pattern, path)
		if err != nil {
			continue // ignore malformed patterns
		}

		if strings.HasPrefix(pattern, "*") {
			if strings.HasSuffix(path, pattern[1:]) {
				match = true
			}
		}

		if strings.HasSuffix(pattern, "/") {
			if strings.Contains(path, pattern) {
				match = true
			}
		}

		logger.L().Debug("Checking pattern",
			"pattern", pattern,
			"path", path,
			"match", match,
			"isNegation", isNegation,
		)
		if match {
			val := !isNegation
			result = &val // last match wins
		}

	}

	return result
}

func Check_Ignore_Scoped(rules map[string][]string, path string) *bool {
	parent := filepath.Dir(path)

	for {
		if ruleset, ok := rules[parent]; ok {
			result := Check_Ignore_1(ruleset, path)
			if result != nil {
				return result
			}
		}

		if parent == "." || parent == "" {
			break
		}

		parent = filepath.Dir(parent)
	}

	return nil
}

func Check_Ignore_Absolute(patterns []string, path string) bool {
	result := Check_Ignore_1(patterns, path)
	if result != nil {
		return *result
	}
	return false
}

func CheckIgnore(rules *GitIgnore, path string) (bool, error) {
	if filepath.IsAbs(path) {
		return false, errors.New("path must be relative to repository root")
	}

	// Scoped Rules
	result := Check_Ignore_Scoped(rules.Scoped, path)
	if result != nil {
		return *result, nil
	}

	// Absolute Rules
	return Check_Ignore_Absolute(rules.Absolute, path), nil
}

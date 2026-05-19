package gitcheckignore

import (
	"errors"
	"fmt"
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

func Check_Ignore_Scoped(rules map[string][]string, path string) *bool {
	path = filepath.ToSlash(path)
	path = filepath.ToSlash(strings.TrimSuffix(path, "/"))
	parent := filepath.ToSlash(filepath.Dir(path))

	for {
		if ruleset, ok := rules[parent]; ok {
			relativePath := path
			if parent != "." && parent != "" {
				relativePath = strings.TrimPrefix(path, parent+"/")
			}

			result := Check_Ignore_1(ruleset, relativePath)
			if result != nil {
				return result
			}
		} else {
		}

		if parent == "." || parent == "" {
			break
		}
		parent = filepath.ToSlash(filepath.Dir(parent))
	}
	return nil
}

// helpers to make prints readable
func scopeKeys(rules map[string][]string) []string {
	keys := make([]string, 0, len(rules))
	for k := range rules {
		keys = append(keys, k)
	}
	return keys
}

func ptrVal(b *bool) string {
	if b == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", *b)
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

func Check_Ignore_1(patterns []string, path string) *bool {
	var result *bool
	// Normalize to forward slashes for consistent matching
	path = filepath.ToSlash(path)

	for _, p := range patterns {
		isNegation := strings.HasPrefix(p, "!")
		pattern := p
		if isNegation {
			pattern = p[1:]
		}
		pattern = filepath.ToSlash(pattern)

		matched := matchGitignorePattern(pattern, path)
		if matched {
			val := !isNegation
			result = &val
		}
	}
	return result
}

func matchGitignorePattern(pattern, path string) bool {
	// Leading "/" means anchored to root — strip it and do prefix match
	anchored := false
	if strings.HasPrefix(pattern, "/") {
		anchored = true
		pattern = pattern[1:]
	}

	if strings.HasSuffix(pattern, "/") {
		dirPattern := strings.TrimSuffix(pattern, "/")
		if anchored {
			// Must match as the first directory component only
			return strings.HasPrefix(path, dirPattern+"/") || path == dirPattern
		}
		return containsDirComponent(path, dirPattern)
	}

	if anchored {
		// Match from root only
		matched, err := filepath.Match(pattern, path)
		if err != nil {
			return false
		}
		// Also match everything inside if it's a directory name
		return matched || strings.HasPrefix(path, pattern+"/")
	}

	// Non-anchored, no slash — match basename or any component
	parts := strings.Split(path, "/")
	base := parts[len(parts)-1]
	if matched, err := filepath.Match(pattern, base); err == nil && matched {
		return true
	}
	for _, part := range parts[:len(parts)-1] {
		if matched, err := filepath.Match(pattern, part); err == nil && matched {
			return true
		}
	}
	return false
}

func containsDirComponent(path, dirPattern string) bool {
	parts := strings.Split(path, "/")
	// Check each directory component in the path
	for i, part := range parts[:len(parts)-1] { // exclude filename
		if matched, err := filepath.Match(dirPattern, part); err == nil && matched {
			_ = i
			return true
		}
	}
	return false
}

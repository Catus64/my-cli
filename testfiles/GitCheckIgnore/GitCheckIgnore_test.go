package gitcheckignore

import (
	"reflect"
	"testing"
)

func TestGitignoreParse1(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedRaw     string
		expectedInclude bool
		expectedOk      bool
	}{
		{
			name:            "comment line",
			input:           "#asdasdas",
			expectedRaw:     "",
			expectedInclude: false,
			expectedOk:      false,
		},
		{
			name:            "empty line",
			input:           "",
			expectedRaw:     "",
			expectedInclude: false,
			expectedOk:      false,
		},
		{
			name:            "negation rule",
			input:           "!test.txt",
			expectedRaw:     "test.txt",
			expectedInclude: false, // depends on your logic
			expectedOk:      true,
		},
		{
			name:            "normal pattern",
			input:           "node_modules/",
			expectedRaw:     "node_modules/",
			expectedInclude: true, // depends on your logic
			expectedOk:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, include, ok := Gitignore_Parse1(tt.input)

			if raw != tt.expectedRaw {
				t.Errorf("raw mismatch: got %q, want %q", raw, tt.expectedRaw)
			}

			if include != tt.expectedInclude {
				t.Errorf("include mismatch: got %v, want %v", include, tt.expectedInclude)
			}

			if ok != tt.expectedOk {
				t.Errorf("ok mismatch: got %v, want %v", ok, tt.expectedOk)
			}
		})
	}
}

func TestCollectPatterns(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected []string
	}{
		{
			name: "filters comments and empty lines",
			lines: []string{
				"#comment",
				"",
				"node_modules/",
				"!important.txt",
			},
			expected: []string{
				"node_modules/", // assuming include=true
			},
		},
		{
			name: "multiple valid patterns",
			lines: []string{
				"build/",
				"dist/",
				"#ignore this",
				"!keep.txt",
			},
			expected: []string{
				"build/",
				"dist/",
			},
		},
		{
			name:     "all ignored",
			lines:    []string{"#a", "", "#b"},
			expected: []string{},
		},
		{
			name:     "empty input",
			lines:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Gitignore_Parse(tt.lines)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

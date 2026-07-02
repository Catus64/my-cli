package GitHashRead

import (
	"bytes"
	"testing"

	gitobj "gocmd/testfiles/GitObject" // adjust this import path
)

// helper to create a GitBlob easily
func newBlob(format, data string) gitobj.GitBlob {
	return gitobj.GitBlob{
		GitObjectData: gitobj.GitObjectData{Data: []byte(data)},
		Format:        []byte(format),
	}
}

func TestBuildGitObjectToWrite_BasicBlob(t *testing.T) {
	blob := newBlob("blob", "hello world")
	result := BuildGitObjectToWrite(blob)

	// expected: "blob 11\x00hello world"
	expected := []byte("blob 11\x00hello world")

	if !bytes.Equal(result, expected) {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestBuildGitObjectToWrite_EmptyData(t *testing.T) {
	blob := newBlob("blob", "")
	result := BuildGitObjectToWrite(blob)

	// expected: "blob 0\x00"
	expected := []byte("blob 0\x00")

	if !bytes.Equal(result, expected) {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestBuildGitObjectToWrite_CorrectNullByteSeparator(t *testing.T) {
	blob := newBlob("blob", "abc")
	result := BuildGitObjectToWrite(blob)

	// find the null byte — must exist
	nullIdx := bytes.IndexByte(result, 0x00)
	if nullIdx == -1 {
		t.Fatal("expected null byte separator, none found")
	}

	// everything after null byte should be the raw data
	afterNull := result[nullIdx+1:]
	if !bytes.Equal(afterNull, []byte("abc")) {
		t.Errorf("data after null byte = %q, want %q", afterNull, "abc")
	}
}

func TestBuildGitObjectToWrite_HeaderFormat(t *testing.T) {
	blob := newBlob("blob", "hello world")
	result := BuildGitObjectToWrite(blob)

	// header is everything before the null byte
	nullIdx := bytes.IndexByte(result, 0x00)
	header := string(result[:nullIdx])

	if header != "blob 11" {
		t.Errorf("header = %q, want %q", header, "blob 11")
	}
}

func TestBuildGitObjectToWrite_DifferentFormat(t *testing.T) {
	// if you ever support tree/commit objects, this checks format flexibility
	blob := newBlob("tree", "some tree data")
	result := BuildGitObjectToWrite(blob)

	if !bytes.HasPrefix(result, []byte("tree ")) {
		t.Errorf("expected result to start with 'tree ', got %q", result)
	}
}

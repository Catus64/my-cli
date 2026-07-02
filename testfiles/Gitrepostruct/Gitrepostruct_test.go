package RepoPath

import (
	helper "gocmd/testfiles/Helper"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// create a temp dir for the log file
	tmpDir, err := os.MkdirTemp("", "test-logs-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir) // cleanup after all tests finish

	logPath := filepath.Join(tmpDir, "test.log")
	if err := helper.Init(logPath, slog.LevelDebug); err != nil {
		panic(err)
	}

	// run all tests
	os.Exit(m.Run())
}

func newTestRepo(t *testing.T) GitRepository {
	t.Helper()
	tmpDir := t.TempDir() // auto-cleaned after each test
	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0755)
	return GitRepository{
		WorkTree: tmpDir,
		GitDir:   gitDir,
		cfg:      nil,
	}
}

func TestRepoPath_SingleSegment(t *testing.T) {
	repo := newTestRepo(t)
	result := Repo_Path(repo, "objects")
	expected := filepath.Join(repo.GitDir, "objects")

	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestRepoPath_MultipleSegments(t *testing.T) {
	repo := newTestRepo(t)
	result := Repo_Path(repo, "objects", "pack")
	expected := filepath.Join(repo.GitDir, "objects", "pack")

	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestRepoPath_NoSegments(t *testing.T) {
	repo := newTestRepo(t)
	result := Repo_Path(repo)
	// filepath.Join("") returns "" so result should just be GitDir
	expected := repo.GitDir

	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestRepoDir_ExistingDir(t *testing.T) {
	repo := newTestRepo(t)
	// objects dir already exists (we made it in newTestRepo via .git)
	os.MkdirAll(filepath.Join(repo.GitDir, "objects"), 0755)

	path, err := Repo_Dir(repo, false, "objects")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != filepath.Join(repo.GitDir, "objects") {
		t.Errorf("got %q, want %q", path, filepath.Join(repo.GitDir, "objects"))
	}
}

func TestRepoDir_MkdirTrue_CreatesDir(t *testing.T) {
	repo := newTestRepo(t)

	path, err := Repo_Dir(repo, true, "objects", "pack")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// check the dir was actually created on disk
	info, statErr := os.Stat(path)
	if os.IsNotExist(statErr) {
		t.Fatal("expected directory to be created, but it doesn't exist")
	}
	if !info.IsDir() {
		t.Errorf("expected a directory, got a file")
	}
}

func TestRepoDir_MkdirFalse_DirNotExist_ReturnsError(t *testing.T) {
	repo := newTestRepo(t)

	_, err := Repo_Dir(repo, false, "nonexistent", "path")
	if err == nil {
		t.Error("expected an error when dir doesn't exist and mkdir=false, got nil")
	}
}

func TestRepoDir_PathIsFile_ReturnsError(t *testing.T) {
	repo := newTestRepo(t)

	// create a FILE where a dir is expected
	filePath := filepath.Join(repo.GitDir, "iamfile")
	os.WriteFile(filePath, []byte("data"), 0644)

	_, err := Repo_Dir(repo, false, "iamfile")
	if err == nil {
		t.Error("expected error when path is a file, not a directory")
	}
}

func TestRepoFile_ReturnsFinalPath(t *testing.T) {
	repo := newTestRepo(t)
	// pre-create the parent dir so Repo_Dir doesn't fail
	os.MkdirAll(filepath.Join(repo.GitDir, "objects"), 0755)

	result := Repo_File(repo, false, "objects", "somefile")
	expected := filepath.Join(repo.GitDir, "objects", "somefile")

	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestRepoFile_MkdirTrue_CreatesParentDir(t *testing.T) {
	repo := newTestRepo(t)

	result := Repo_File(repo, true, "objects", "pack", "somefile")
	expected := filepath.Join(repo.GitDir, "objects", "pack", "somefile")

	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}

	// parent dir should exist, the file itself should NOT (Repo_File doesn't create it)
	parentDir := filepath.Join(repo.GitDir, "objects", "pack")
	info, err := os.Stat(parentDir)
	if os.IsNotExist(err) {
		t.Fatal("expected parent dir to be created")
	}
	if !info.IsDir() {
		t.Error("expected parent path to be a directory")
	}
}

func TestRepoFile_PanicsWhenDirMissing(t *testing.T) {
	repo := newTestRepo(t)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when parent dir missing and mkdir=false, but did not panic")
		}
	}()

	// "nonexistent" dir doesn't exist and mkdir=false → should panic
	Repo_File(repo, false, "nonexistent", "somefile")
}

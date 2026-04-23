package gitpacketextractor

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

const sentinelFile = ".git/objects/pack/.unpacked"
const backupDir = "/tmp/pack_backup"

//First run - moves packgiles out of .git/objects/pack,
// unpacks them, and writes a sentinel file

// Subsequent runs - checks for sentinel file,
// if exists, skip extraction

// This is a one-time extraction process to unpack Git packfiles into loose objects.

// AlreadyExtracted checks if the sentinel file exists,
// meaning extraction has already been done before.
func AlreadyExtracted() bool {
	_, err := os.Stat(sentinelFile)
	return err == nil
}

// Extract moves all packfiles out, unpacks them into loose objects,
// and writes a sentinel file so it won't run again next time.
func Extract() error {
	packDir := ".git/objects/pack"

	if AlreadyExtracted() {
		slog.Info("packfile extraction already done, skipping")
		return nil
	}

	slog.Info("starting packfile extraction")

	// 1. Create backup dir
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup dir: %w", err)
	}

	// 2. Move all pack files out
	entries, err := os.ReadDir(packDir)
	if err != nil {
		// No pack dir means nothing to extract, write sentinel and return
		slog.Info("no pack directory found, nothing to extract")
		return writeSentinel()
	}

	var packFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		src := filepath.Join(packDir, entry.Name())
		dst := filepath.Join(backupDir, entry.Name())

		if err := moveFile(src, dst); err != nil {
			return fmt.Errorf("failed to move %s: %w", src, err)
		}
		slog.Debug("moved packfile", "src", src, "dst", dst)

		if filepath.Ext(entry.Name()) == ".pack" {
			packFiles = append(packFiles, dst)
		}
	}

	if len(packFiles) == 0 {
		slog.Info("no .pack files found, nothing to unpack")
		return writeSentinel()
	}

	// 3. Unpack each .pack file
	for _, pack := range packFiles {
		slog.Info("unpacking", "file", pack)

		if err := unpackFile(pack); err != nil {
			return fmt.Errorf("failed to unpack %s: %w", pack, err)
		}

		slog.Info("successfully unpacked", "file", pack)
	}

	// 4. Write sentinel so we skip next time
	if err := writeSentinel(); err != nil {
		return fmt.Errorf("failed to write sentinel: %w", err)
	}

	slog.Info("packfile extraction complete")
	return nil
}

// unpackFile pipes a single .pack file into git unpack-objects
func unpackFile(packPath string) error {
	f, err := os.Open(packPath)
	if err != nil {
		return fmt.Errorf("failed to open pack file: %w", err)
	}
	defer f.Close()

	cmd := exec.Command("git", "unpack-objects")
	cmd.Stdin = f
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// writeSentinel creates the marker file that signals extraction is done
func writeSentinel() error {
	// Ensure pack dir still exists (it may be empty now)
	if err := os.MkdirAll(".git/objects/pack", 0755); err != nil {
		return err
	}
	f, err := os.Create(sentinelFile)
	if err != nil {
		return err
	}
	f.Close()
	slog.Debug("sentinel file written", "path", sentinelFile)
	return nil
}

// moveFile handles cross-device moves by falling back to copy+delete
func moveFile(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return os.Remove(src)
}

package packextractor

import (
	"fmt"
	logger "gocmd/testfiles/Helper"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	gitpath "gocmd/testfiles/Gitrepostruct"
)

const backupDir = "/tmp/pack_backup"

func sentinelPath(repo gitpath.GitRepository) string {
	return gitpath.Repo_Path(repo, "objects", "pack", ".unpacked")
}

func packDirPath(repo gitpath.GitRepository) string {
	return gitpath.Repo_Path(repo, "objects", "pack")
}

func AlreadyExtracted(repo gitpath.GitRepository) bool {
	_, err := os.Stat(sentinelPath(repo))
	return err == nil
}

func Extract(repo gitpath.GitRepository) error {
	if AlreadyExtracted(repo) {
		logger.L().Debug("extract: packfile extraction already done, skipping")
		return nil
	}

	logger.L().Info("extract: starting packfile extraction")

	packDir := packDirPath(repo)

	// 1. Create backup dir
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup dir: %w", err)
	}

	// 2. Move all pack files out
	entries, err := os.ReadDir(packDir)
	if err != nil {
		logger.L().Info("extract: no pack directory found, nothing to extract")
		return writeSentinel(repo)
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
		logger.L().Debug("extract: moved packfile", "src", src, "dst", dst)

		if filepath.Ext(entry.Name()) == ".pack" {
			packFiles = append(packFiles, dst)
		}
	}

	if len(packFiles) == 0 {
		logger.L().Info("extract: no .pack files found, nothing to unpack")
		return writeSentinel(repo)
	}

	// 3. Unpack each .pack file
	for _, pack := range packFiles {
		logger.L().Info("extract: unpacking", "file", pack)
		if err := unpackFile(pack); err != nil {
			return fmt.Errorf("failed to unpack %s: %w", pack, err)
		}
		logger.L().Info("extract: successfully unpacked", "file", pack)
	}

	// 4. Write sentinel
	if err := writeSentinel(repo); err != nil {
		return fmt.Errorf("failed to write sentinel: %w", err)
	}

	logger.L().Info("extract: packfile extraction complete")
	return nil
}

func ForceExtract(repo gitpath.GitRepository) error {
	sentinelPath := filepath.Join(gitpath.Repo_Path(repo, ""), "ezgit_extracted") // adjust to match your actual sentinel filename/path
	if err := os.Remove(sentinelPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clear sentinel: %w", err)
	}
	return Extract(repo)
}

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

func writeSentinel(repo gitpath.GitRepository) error {
	packDir := packDirPath(repo)
	if err := os.MkdirAll(packDir, 0755); err != nil {
		return err
	}
	f, err := os.Create(sentinelPath(repo))
	if err != nil {
		return err
	}
	f.Close()
	logger.L().Debug("extract: sentinel file written", "path", sentinelPath(repo))
	return nil
}

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

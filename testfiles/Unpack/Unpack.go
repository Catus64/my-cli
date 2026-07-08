package unpack

import (
	"bytes"
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func UnpackPackfiles(repo gitpath.GitRepository) error {
	packDir := filepath.Join(repo.GitDir, "objects", "pack")

	entries, err := os.ReadDir(packDir)
	if err != nil {
		return nil // no pack directory, nothing to do
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".pack") {
			continue
		}
		packPath := filepath.Join(packDir, entry.Name())
		logger.L().Debug("Unpacking packfile", "pack", packPath)

		cmd := exec.Command("git", "unpack-objects", "-r")

		packData, err := os.ReadFile(packPath)
		if err != nil {
			return err
		}
		cmd.Stdin = bytes.NewReader(packData)
		cmd.Dir = repo.WorkTree

		out, err := cmd.CombinedOutput()
		logger.L().Debug("unpack-objects output", "out", string(out))
		if err != nil {
			return fmt.Errorf("unpack-objects failed: %w", err)
		}
	}
	return nil
}

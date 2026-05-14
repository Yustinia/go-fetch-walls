package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var WgetFlags = []string{
	"--show-progress",
	"-P",
}

func setDlPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("HOME not found: %w", err)
	}
	dlDir := "Downloads"
	outPath := filepath.Join(homeDir, dlDir)

	return outPath, nil
}

func Downloader(wall string) error {
	dlPath, err := setDlPath()
	if err != nil {
		return err
	}

	args := append(WgetFlags, dlPath, wall)
	cmd := exec.Command("wget", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to download %v: %w\n", wall, err)
	}

	return nil
}

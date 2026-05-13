package cmd

import (
	"fmt"
	"go-fetch-walls/api"
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

func Downloader(result *api.Response) error {
	dlPath, err := setDlPath()
	if err != nil {
		return err
	}

	for _, wall := range result.Data {
		args := append(WgetFlags, dlPath, wall.Path)
		cmd := exec.Command("wget", args...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to download %v: %w\n", wall.Path, err)
		}
	}

	return nil
}

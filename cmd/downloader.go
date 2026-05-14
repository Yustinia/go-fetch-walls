package cmd

import (
	"fmt"
	"go-fetch-walls/internal"
	"io"
	"net/http"
	"os"
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

func WallDownloader(wall internal.Wallpaper) error {
	dlPath, err := setDlPath()
	if err != nil {
		return err
	}

	resp, err := http.Get(wall.Path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath.Join(dlPath, filepath.Base(wall.Path)))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

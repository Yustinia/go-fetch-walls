package main

import (
	"fmt"
	"go-fetch-walls/api"
	"go-fetch-walls/internal"
	"go-fetch-walls/tui"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

func printWallData(result *api.Response) {
	for index, wall := range result.Data {
		fmt.Printf("[%d]: %v\n", index, wall.Path)
		fmt.Printf("	Category: %v\n", wall.Category)
		fmt.Printf("   	Purity: %v\n", wall.Purity)
		fmt.Printf("   	Resolution: %v\n", wall.Resolution)
		fmt.Println()
	}
}

func devModeSettings(mode bool) (string, error) {
	if mode {
		return "configs/dev_settings.json", nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("HOME not found: %w", err)
	}

	fullPath := filepath.Join(homeDir, ".config", "go-fetch-walls", "settings.json")
	return fullPath, nil
}

func main() {
	const DevMode = true

	configPath, err := devModeSettings(DevMode)
	if err != nil {
		panic(err)
	}

	baseURL := "https://wallhaven.cc/api/v1/search?"
	var settings internal.Settings

	err = internal.LoadSettings(configPath, &settings)
	if err != nil {
		panic(err)
	}
	err = internal.ValidateSettings(&settings)
	if err != nil {
		panic(err)
	}

	params := api.BuildParams(&settings)
	fullURL := baseURL + params.Encode()

	result, err := api.GetResponse(fullURL)
	if err != nil {
		panic(err)
	}

	p := tea.NewProgram(tui.WallsModel(result))
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

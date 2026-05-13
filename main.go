package main

import (
	"fmt"
	"go-fetch-walls/api"
	"go-fetch-walls/internal"
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

type model struct {
	walls  []internal.Wallpaper
	cursor int
}

func initialModel(result api.Response) model {
	return model{walls: result.Data}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.walls)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	s := ""
	for index, wall := range m.walls {
		cursor := "  "
		if m.cursor == index {
			cursor = "> "
		}
		s += fmt.Sprintf("%s %s\n", cursor, wall.Path)
	}

	return tea.NewView(s)
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

	p := tea.NewProgram(initialModel(result))
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	// printWallData(&result)

	// err = cmd.Downloader(&result)
	// if err != nil {
	// 	panic(err)
	// }
}

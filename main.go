package main

import (
	"fmt"
	"go-fetch-walls/api"
	"go-fetch-walls/internal"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

type state int

const (
	stateList state = iota
	stateDownload
)

type model struct {
	walls   []internal.Wallpaper
	cursor  int
	current state
}

func wallsModel(result api.Response) model {
	return model{
		walls:   result.Data,
		cursor:  0,
		current: stateList,
	}
}

var (
	leftCol  = lipgloss.NewStyle().Width(70)
	rightCol = lipgloss.NewStyle().Width(80)
)

func renderMeta(m model) string {
	wall := m.walls[m.cursor]
	return fmt.Sprintf(
		"Path: %s\nCategory: %s\nPurity: %s\nResolution: %s\n",
		wall.Path,
		wall.Category,
		wall.Purity,
		wall.Resolution,
	)
}

func renderList(m model) string {
	s := ""
	for index, wall := range m.walls {
		cursor := " "

		if m.cursor == index {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, wall.Path)
	}

	return s
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
		case "up", "h":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "l":
			if m.cursor < len(m.walls)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	switch m.current {
	case stateList:
		left := leftCol.Render(renderList(m))
		right := rightCol.Render(renderMeta(m))
		return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top, left, right))
	}

	return tea.NewView("")
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

	p := tea.NewProgram(wallsModel(result))
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

package tui

import (
	"fmt"
	"go-fetch-walls/api"
	"go-fetch-walls/cmd"
	"go-fetch-walls/internal"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type state int
type dlStatus int
type downloadDoneMsg struct{}
type newPageMsg struct {
	walls []internal.Wallpaper
}

const (
	dlIdle dlStatus = iota
	dlInProgress
	dlDone
)

const (
	stateList state = iota
	stateDownload
)

type model struct {
	walls    []internal.Wallpaper
	cursor   int
	current  state
	download dlStatus
	page     uint
	baseURL  string
}

func WallsModel(result api.Response, baseURL string) model {
	return model{
		walls:   result.Data,
		cursor:  0,
		current: stateList,
		page:    1,
		baseURL: baseURL,
	}
}

var (
	leftCol  = lipgloss.NewStyle().Width(70)
	rightCol = lipgloss.NewStyle().Width(80)

	metaStyle   = lipgloss.NewStyle().Height(20)
	statusStyle = lipgloss.NewStyle().Height(5)
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

func renderStatus(m model) string {
	wall := m.walls[m.cursor]
	switch m.download {
	case dlIdle:
		return fmt.Sprint("Enter to Download")

	case dlInProgress:
		return fmt.Sprint("Downloading...")

	case dlDone:
		return fmt.Sprintf("Finished %v", wall.Path)
	}

	return ""
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
				m.download = dlIdle
			}
		case "down", "l":
			if m.cursor < len(m.walls)-1 {
				m.cursor++
				m.download = dlIdle
			}
		case "enter":
			wall := m.walls[m.cursor]
			return m, func() tea.Msg {
				err := cmd.WallDownloader(wall)
				if err != nil {
					return err
				}
				return downloadDoneMsg{}
			}
		case "n":
			m.page++
			return m, func() tea.Msg {
				result, err := api.GetResponse(fmt.Sprintf("%s&page=%d", m.baseURL, m.page))
				if err != nil {
					return err
				}
				return newPageMsg{walls: result.Data}
			}
		}

	case downloadDoneMsg:
		m.download = dlDone
		return m, nil

	case newPageMsg:
		m.walls = msg.walls
		m.cursor = 0
		m.download = dlIdle
		return m, nil
	}

	return m, nil
}

func (m model) View() tea.View {
	switch m.current {
	case stateList:
		left := leftCol.Render(renderList(m))

		meta := metaStyle.Render(renderMeta(m))
		status := statusStyle.Render(renderStatus(m))

		right := rightCol.Render(lipgloss.JoinVertical(lipgloss.Top, meta, status))

		return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top, left, right))
	}

	return tea.NewView("")
}

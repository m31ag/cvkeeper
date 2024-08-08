package model

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/m31ag/cvkeeper/repo"
	"strings"
)

const (
	emptyCursor  = "  "
	filledCursor = "->"
	menuFormat   = " %s %s\n"
)

var (
	style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#CE6797"))
)

type Model struct {
	repo    repo.Repository
	files   []repo.File
	cursor  int
	checked int
}

func InitModel(r repo.Repository) Model {
	files := r.GetFilesByParentId(-1)
	checked := 0
	if len(files) > 0 {
		checked = files[0].Id
	}
	return Model{
		repo:    r,
		files:   files,
		checked: checked,
	}
}
func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}
		case "b":
			if m.checked != 0 {
				return m.Back(), nil
			}
		case "enter", " ":
			return m.Forward(), nil
		}
	}

	return m, nil
}
func (m Model) View() string {
	s := fmt.Sprintf("%s\n\n", lipgloss.NewStyle().Width(20).MarginLeft(20).AlignHorizontal(lipgloss.Center).Bold(true).Background(lipgloss.Color("63")).Render("Working tool"))

	for i, item := range m.files {

		cursor := emptyCursor
		colored := false
		if m.cursor == i {
			cursor = filledCursor
			colored = true
		}

		// Render the row

		s += showItem(fmt.Sprintf(menuFormat, cursor, item.Filename), colored)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func showItem(txt string, colored bool) string {
	if colored {
		return strings.TrimSpace(style.Render(txt))
	} else {
		return txt
	}

}
func (m Model) Back() Model {
	files := m.repo.GetFilesByParentId(m.files[0].ParentId)
	m.files = files
	m.cursor = 0
	return m
}
func (m Model) Forward() Model {
	files := m.repo.GetFilesByParentId(m.files[0].Id)
	m.files = files
	m.cursor = 0
	return m
}

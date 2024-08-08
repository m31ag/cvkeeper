package model

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	items    []Item
	cursor   int
	selected map[int]struct{}
	level    int
}

func InitModel() Model {
	c, l := GetChoicesByLevel(0)
	return Model{
		items:    c,
		level:    l,
		selected: make(map[int]struct{}),
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
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "b":
			if m.level > 0 {
				return m.ChangeLevel(m.level - 1), nil
			}
		case "enter", " ":
			return m.ChangeLevel(m.level + 1), nil
		}
	}

	return m, nil
}
func (m Model) View() string {
	s := fmt.Sprintf("%s\n\n", lipgloss.NewStyle().Width(20).MarginLeft(20).AlignHorizontal(lipgloss.Center).Bold(true).Background(lipgloss.Color("63")).Render("Working tool"))

	for i, item := range m.items {

		cursor := emptyCursor
		colored := false
		if m.cursor == i {
			cursor = filledCursor
			colored = true
		}

		// Render the row

		s += showItem(fmt.Sprintf(menuFormat, cursor, item.Name), colored)
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
func (m Model) ChangeLevel(level int) Model {
	c, l := GetChoicesByLevel(level)
	m.items = c
	m.level = l
	m.cursor = 0
	return m
}

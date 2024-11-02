package model

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
	"github.com/m31ag/cvkeeper/repo"
	"strings"
)

type CreatingState int

const (
	emptyCursor   = "  "
	filledCursor  = "->"
	menuFormat    = " %s %s\n"
	historyFormat = "\n%s\n\n"

	UsualState           CreatingState = 0
	WaitFilenameState    CreatingState = 1
	WaitFileContentState CreatingState = 2
)

var (
	style        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#CE6797"))
	historyStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
)

type Input struct {
	input textinput.Model
	area  textarea.Model
}
type Model struct {
	repo            repo.Repository
	files           []repo.File
	cursor          int
	checked         int
	path            []int
	history         []string
	input           Input
	CreatingStateId CreatingState
}

func (m Model) NextState() {
	if m.CreatingStateId == UsualState {
		m.CreatingStateId = WaitFilenameState
	} else if m.CreatingStateId == WaitFilenameState {
		m.CreatingStateId = WaitFileContentState
	} else {
		m.CreatingStateId = UsualState
	}
}
func (m Model) GetCurrentFolderId() int {
	return m.path[len(m.path)-1]
}
func InitModel(r repo.Repository) Model {
	files := r.GetRoot()

	checked := 0
	if len(files) > 0 {
		checked = files[0].Id
	}
	path := make([]int, 0)
	path = append(path, -1)
	history := make([]string, 0)
	history = append(history, "/root")

	return Model{
		repo:    r,
		files:   files,
		checked: checked,
		path:    path,
		history: history,
	}
}
func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.CreatingStateId == UsualState {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
					return m, cmd
				}

			case "down", "j":
				if m.cursor < len(m.files)-1 {
					m.cursor++
					return m, cmd
				}
			case "b", "left", "h":
				if m.checked != 0 {
					return m.Back(), cmd
				}
			case "n":
				m.CreatingStateId = WaitFilenameState
				return m.SetInput("filename"), cmd
			case "enter", " ", "right", "l":
				return m.Forward(), cmd
			}
		} else if m.CreatingStateId == WaitFilenameState {

			switch msg.String() {

			case "enter":
				if err := m.repo.Save(m.input.input.Value(), "test", m.GetCurrentFolderId()); err != nil {
					return m, tea.Quit
				}
				m.CreatingStateId = UsualState
				m.files = m.repo.GetFilesByParentId(m.GetCurrentFolderId())
				return m, cmd
			default:
				m.input.input, cmd = m.input.input.Update(msg)
				return m, cmd
			}

		}

	}
	return m, cmd
}
func (m Model) View() string {
	s := fmt.Sprintf(
		"%s\n\n",
		lipgloss.
			NewStyle().
			MarginLeft(20).
			PaddingLeft(10).
			PaddingRight(10).
			AlignHorizontal(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#5f5fff")).
			Render("CVKeeper"))
	if m.CreatingStateId == UsualState {
		s += render(strings.Join(m.history, "/"), historyFormat, historyStyle)
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
	} else {
		s += m.input.input.View()
	}
	s += "\nPress n to add one-string, N to add multiple-string\nPress f to create folder"
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
func render(txt, format string, style lipgloss.Style) string {
	return strings.TrimSpace(style.Render(fmt.Sprintf(format, txt)))
}
func (m Model) SetInput(placeholder string) Model {
	t := textinput.New()
	t.Placeholder = placeholder
	t.Focus()
	m.input.input = t
	return m
}
func (m Model) Back() Model {
	if len(m.path) > 1 {
		files := m.repo.GetFilesByParentId(m.path[len(m.path)-2])
		m.files = files
		m.path = m.path[:len(m.path)-1]
		m.history = m.history[:len(m.history)-1]
		m.cursor = 0
	}
	return m
}
func (m Model) Forward() Model {
	if len(m.files) != 0 {
		files := m.repo.GetFilesByParentId(m.files[m.cursor].Id)

		m.path = append(m.path, m.files[m.cursor].Id)
		m.history = append(m.history, m.files[m.cursor].Filename)
		m.files = files
		m.cursor = 0
	}
	return m
}

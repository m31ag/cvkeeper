package model

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/m31ag/cvkeeper/repo"
	"strings"
)

type ViewState int

const (
	defaultRootId     int = 0
	defaultFirstDirId     = -2

	emptyCursor   = "  "
	filledCursor  = "->"
	menuFormat    = "%s %s %s\n"
	historyFormat = "\n%s\n\n"

	StandardState        ViewState = 0
	WaitFilenameState    ViewState = 1
	WaitDirnameState     ViewState = 2
	WaitFileContentState ViewState = 3
	ShowFileContentState ViewState = 4
	DeleteState          ViewState = 5
)

var (
	style        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#CE6797"))
	historyStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
)

type Input struct {
	input textinput.Model
	value string
}

type Model struct {
	repo        repo.Repository
	files       []repo.File
	cursor      int
	order       []repo.File
	history     []string
	input       Input
	StateId     ViewState
	fileContent string
}

func (m Model) GetChecked() repo.File {
	if len(m.files) > 0 {
		return m.files[m.cursor]
	}
	return repo.File{
		Id:       0,
		ParentId: 0,
		Filename: "",
		IsFolder: true,
	}
}

func (m Model) GetCurrentOrderId() int {
	return m.order[len(m.order)-1].Id
}
func InitModel(r repo.Repository) Model {
	files := r.GetRoot()

	order := make([]repo.File, 0)
	root := r.GetFilesByParentId(defaultRootId)
	order = append(order, root[0])
	history := make([]string, 0)
	history = append(history, "/root")

	return Model{
		repo:    r,
		files:   files,
		order:   order,
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

		if m.StateId == StandardState {
			return m.OnStandard(msg)
		} else if m.StateId == WaitFilenameState {
			return m.OnWaitFilename(msg)
		} else if m.StateId == WaitFileContentState {
			return m.OnWaitFileContent(msg)
		} else if m.StateId == WaitDirnameState {
			return m.OnWaitDirName(msg)
		} else if m.StateId == ShowFileContentState {
			return m.OnShowFileContent(msg)
		}

	}
	return m, cmd
}
func (m Model) View() string {
	//header
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

	s += render(strings.Join(m.history, "/"), historyFormat, historyStyle)
	//content
	if m.StateId == StandardState {
		for i, item := range m.files {

			cursor := emptyCursor
			colored := false
			if m.cursor == i {
				cursor = filledCursor
				colored = true
			}

			// Render the row
			suffix := "\U0001F4C4"
			if item.IsFolder {
				suffix = "\U0001F4C1"
			}
			s += showItem(fmt.Sprintf(menuFormat, suffix, cursor, item.Filename), colored)

		}
	} else if m.StateId == ShowFileContentState {
		s += m.fileContent
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
	if len(m.order) > 0 {
		var files []repo.File
		if len(m.fileContent) > 0 {
			files = m.repo.GetFilesByParentId(m.order[len(m.order)-1].Id)
			m.fileContent = ""
			m.StateId = StandardState
		} else {
			files = m.repo.GetFilesByParentId(m.order[len(m.order)-1].ParentId)
			m.order = m.order[:len(m.order)-1]
			m.history = m.history[:len(m.history)-1]
		}

		m.files = files
		if len(files) > 0 && files[0].Id == defaultFirstDirId {
			m.cursor = 1
		} else {
			m.cursor = 0
		}
	}
	return m
}
func (m Model) Forward() Model {

	if m.GetChecked().Id != defaultFirstDirId {
		files := m.repo.GetFilesByParentId(m.files[m.cursor].Id)

		m.order = append(m.order, m.files[m.cursor])
		m.history = append(m.history, m.files[m.cursor].Filename)
		m.files = files
		if len(files) > 1 && files[0].Id == defaultFirstDirId {
			m.cursor = 1
		} else {
			m.cursor = 0
		}
	} else {
		return m.Back()
	}
	return m

}

package model

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
	"github.com/m31ag/cvkeeper/repo"
	"strings"
)

type InputState int

const (
	defaultRootId int = 0

	emptyCursor   = "  "
	filledCursor  = "->"
	menuFormat    = "%s %s %s\n"
	historyFormat = "\n%s\n\n"

	StandardState        InputState = 0
	WaitFilenameState    InputState = 1
	WaitDirnameState                = 2
	WaitFileContentState InputState = 3
)

var (
	style        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#CE6797"))
	historyStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
)

type Input struct {
	input        textinput.Model
	inputValue   string
	InputStateId InputState
}

type Model struct {
	repo        repo.Repository
	files       []repo.File
	cursor      int
	order       []repo.File
	history     []string
	input       Input
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
func (i *Input) NextState() {
	if i.InputStateId == StandardState {
		i.InputStateId = WaitFilenameState
	} else if i.InputStateId == WaitFilenameState {
		i.InputStateId = WaitFileContentState
	} else if i.InputStateId == WaitDirnameState || i.InputStateId == WaitFileContentState {
		i.InputStateId = StandardState
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

		if m.input.InputStateId == StandardState {
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
				if len(m.order) > 1 {
					return m.Back(), cmd
				}
			case "f":
				if m.order[len(m.order)-1].IsFolder {
					m.input.InputStateId = WaitDirnameState
					return m.SetInput("dirname"), cmd
				}
			case "n":
				if m.order[len(m.order)-1].IsFolder {
					m.input.InputStateId = WaitFilenameState
					return m.SetInput("filename"), cmd
				}
			case "enter", " ", "right", "l":
				if m.GetChecked().IsFolder {
					return m.Forward(), cmd
				} else {
					c, err := m.repo.GetFileContentByFileId(m.GetChecked().Id)
					if err != nil {
						println(err.Error())
						return m, tea.Quit
					}
					m.fileContent = c
					return m, cmd
				}
			}
		} else if m.input.InputStateId == WaitFilenameState {

			switch msg.String() {

			case "enter":

				m.input.NextState()
				m.input.inputValue = m.input.input.Value()
				return m.SetInput("content"), cmd
			default:
				m.input.input, cmd = m.input.input.Update(msg)
				return m, cmd
			}

		} else if m.input.InputStateId == WaitFileContentState {

			switch msg.String() {
			case "enter":
				if err := m.repo.SaveFileWithContent(m.input.inputValue, m.input.input.Value(), m.GetCurrentOrderId()); err != nil {
					return m, tea.Quit
				}
				m.input.NextState()
				m.files = m.repo.GetFilesByParentId(m.GetCurrentOrderId())
				m.input.inputValue = ""
				return m, cmd
			default:
				m.input.input, cmd = m.input.input.Update(msg)
				return m, cmd
			}

		} else if m.input.InputStateId == WaitDirnameState {
			switch msg.String() {
			case "enter":
				if err := m.repo.SaveDir(m.input.input.Value(), m.GetCurrentOrderId()); err != nil {
					println(err.Error())
					return m, tea.Quit
				}
				m.input.InputStateId = StandardState
				m.files = m.repo.GetFilesByParentId(m.GetCurrentOrderId())
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
	if m.input.InputStateId == StandardState && len(m.fileContent) == 0 {
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
	} else if len(m.fileContent) > 0 {
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

		} else {
			files = m.repo.GetFilesByParentId(m.order[len(m.order)-1].ParentId)
			m.order = m.order[:len(m.order)-1]
			m.history = m.history[:len(m.history)-1]
		}

		m.files = files
		m.cursor = 0
	}
	return m
}
func (m Model) Forward() Model {
	files := m.repo.GetFilesByParentId(m.files[m.cursor].Id)
	m.order = append(m.order, m.files[m.cursor])
	m.history = append(m.history, m.files[m.cursor].Filename)
	m.files = files
	m.cursor = 0
	return m
}

package model

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/m31ag/cvkeeper/repo"
)

type ViewState int

const (
	purpleColor = "#CE6797"
	whiteColor  = "#FFFFFF"

	defaultRootId     int = 0
	defaultFirstDirId     = -2

	emptyCursor   = "  "
	filledCursor  = "->"
	menuFormat    = "%s %s %s\n"
	historyFormat = "\n%s\n\n"

	StandardState                ViewState = 0
	WaitFilenameState            ViewState = 1
	WaitDirnameState             ViewState = 2
	WaitFileContentState         ViewState = 3
	ShowFileContentState         ViewState = 4
	DeleteState                  ViewState = 5
	WaitFilenameMultiStringState ViewState = 6
	WaitMultipleFileContentState ViewState = 7
)

var (
	style        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(purpleColor))
	historyStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(whiteColor))
)

type Input struct {
	input textinput.Model
	value string
}
type Area struct {
	area  textarea.Model
	value string
}
type Model struct {
	repo        repo.Repository
	files       []repo.File
	cursor      int
	order       []repo.File
	history     []string
	input       Input
	area        Area
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
	return m.GetCurrentOrder().Id
}
func (m Model) GetCurrentOrder() repo.File {
	return m.order[len(m.order)-1]
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
			return m.OnStandardUpdate(msg)
		} else if m.StateId == WaitFilenameState {
			return m.OnWaitFilenameUpdate(msg)
		} else if m.StateId == WaitFileContentState {
			return m.OnWaitFileContentUpdate(msg)
		} else if m.StateId == WaitFilenameMultiStringState {
			return m.OnWaitFilenameMultiStringUpdate(msg)
		} else if m.StateId == WaitMultipleFileContentState {
			return m.OnWaitMultipleFileContentUpdate(msg)
		} else if m.StateId == WaitDirnameState {
			return m.OnWaitDirNameUpdate(msg)
		} else if m.StateId == ShowFileContentState {
			return m.OnShowFileContentUpdate(msg)
		} else if m.StateId == DeleteState {
			return m.OnDeleteUpdate(msg)
		}

	}
	return m, cmd
}
func (m Model) View() string {
	switch m.StateId {
	case StandardState:
		return m.OnStandardView()
	case ShowFileContentState:
		return m.OnShowFileContentView()
	case DeleteState:
		return m.OnDeleteView()
	case WaitFilenameState, WaitDirnameState, WaitFilenameMultiStringState, WaitFileContentState:
		return m.DefaultInputView()
	case WaitMultipleFileContentState:
		return m.DefaultAreaView()
	default:
		return m.OnStandardView()
	}

}

func (m Model) SetInput(placeholder string) Model {
	t := textinput.New()
	t.Placeholder = placeholder
	t.Focus()
	m.input.input = t
	return m
}
func (m Model) SetArea(placeholder string) Model {
	t := textarea.New()
	t.Placeholder = placeholder
	t.Focus()
	m.area.area = t
	return m
}
func (m Model) SetDefaultCursor() Model {
	if len(m.files) > 1 && m.files[0].Id == defaultFirstDirId {
		m.cursor = 1
	} else {
		m.cursor = 0
	}
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
		return m.SetDefaultCursor()
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

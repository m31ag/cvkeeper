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
)

var (
	style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(purpleColor))
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
	vars        Vars
}

func (m Model) GetChecked() repo.File {
	if len(m.files) > 0 {
		return m.files[m.cursor]
	}
	return repo.File{IsFolder: true}
}

func (m Model) GetCurrentOrderId() int {
	return m.GetCurrentOrder().Id
}
func (m Model) GetCurrentOrder() repo.File {
	return m.order[len(m.order)-1]
}
func InitModel(r repo.Repository, v Vars) Model {
	files := r.GetRoot()
	root := r.GetFilesByParentId(defaultRootId)

	return Model{
		repo:    r,
		files:   files,
		order:   []repo.File{root[0]},
		history: []string{"/root"},
		vars:    v,
	}
}
func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if fn, ok := updateMap[m.StateId]; ok {
			return fn(m, msg)
		}
	}
	return m, cmd
}
func (m Model) View() string {
	if fn, ok := viewMap[m.StateId]; ok {
		return fn(m)
	}
	return m.OnStandardView()
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
	m.area.area.SetWidth(50)
	return m
}

// clampCursor resets cursor to first real file/folder,
// skipping "/.." back folder if any
func (m Model) clampCursor() Model {
	if len(m.files) > 0 && m.files[0].Id == defaultFirstDirId {
		m.cursor = 1
	} else {
		m.cursor = 0
	}
	return m
}

// MoveCursor move cursor on delta positions without going beyond the boundaries
func (m Model) MoveCursor(delta int) Model {
	m.cursor = clamp(m.cursor+delta, 0, len(m.files)-1)
	return m
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Back handles the "back" press.
// If the content of the file is open, closes it, remaining in the same folder.
// Otherwise - goes higher.
func (m Model) Back() Model {
	// for root folder
	if len(m.order) == 1 && m.order[0].Id == -1 {
		if len(m.fileContent) > 0 {
			m.fileContent = ""
			m.StateId = StandardState
		}
		return m
	}

	//any nested folder
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
		return m.clampCursor()
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
func LoadVars() {

}

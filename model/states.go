package model

import tea "github.com/charmbracelet/bubbletea"

const (
	StandardState                ViewState = 0
	WaitFilenameState            ViewState = 1
	WaitDirnameState             ViewState = 2
	WaitFileContentState         ViewState = 3
	ShowFileContentState         ViewState = 4
	DeleteState                  ViewState = 5
	WaitFilenameMultiStringState ViewState = 6
	WaitMultipleFileContentState ViewState = 7
)

// updates
type UpdateFunc func(m Model, msg tea.Msg) (tea.Model, tea.Cmd)

var updateMap = map[ViewState]UpdateFunc{
	ShowFileContentState:         Model.OnShowFileContentUpdate,
	DeleteState:                  Model.OnDeleteUpdate,
	WaitFilenameState:            Model.OnWaitFilenameUpdate,
	WaitDirnameState:             Model.OnWaitDirNameUpdate,
	WaitFileContentState:         Model.OnWaitFileContentUpdate,
	WaitFilenameMultiStringState: Model.OnWaitFilenameMultiStringUpdate,
	WaitMultipleFileContentState: Model.OnWaitMultipleFileContentUpdate,
	StandardState:                Model.OnStandardUpdate,
}

// views
type ViewFunc func(m Model) string

var viewMap = map[ViewState]ViewFunc{
	StandardState:                Model.OnStandardView,
	ShowFileContentState:         Model.OnShowFileContentView,
	DeleteState:                  Model.OnDeleteView,
	WaitFilenameState:            Model.DefaultInputView,
	WaitDirnameState:             Model.DefaultInputView,
	WaitFilenameMultiStringState: Model.DefaultInputView,
	WaitFileContentState:         Model.DefaultInputView,
	WaitMultipleFileContentState: Model.DefaultAreaView,
}

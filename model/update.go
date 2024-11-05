package model

import tea "github.com/charmbracelet/bubbletea"

func (m Model) OnStandard(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "delete":
			m.StateId = DeleteState
			return m, cmd
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
				m.StateId = WaitDirnameState
				return m.SetInput("dirname"), cmd
			}
		case "n":
			if m.order[len(m.order)-1].IsFolder {
				m.StateId = WaitFilenameState
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
				m.StateId = ShowFileContentState
				m.fileContent = c
				return m, cmd
			}
		}
	}
	return m, cmd
}
func (m Model) OnWaitFilename(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "enter":

			m.StateId = WaitFileContentState
			m.input.value = m.input.input.Value()
			return m.SetInput("content"), cmd
		default:
			m.input.input, cmd = m.input.input.Update(msg)
			return m, cmd
		}
	}
	return m, cmd
}

func (m Model) OnWaitFileContent(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if err := m.repo.SaveFileWithContent(m.input.value, m.input.input.Value(), m.GetCurrentOrderId()); err != nil {
				return m, tea.Quit
			}
			m.StateId = StandardState
			m.files = m.repo.GetFilesByParentId(m.GetCurrentOrderId())
			m.input.value = ""
			return m, cmd
		default:
			m.input.input, cmd = m.input.input.Update(msg)
			return m, cmd
		}
	}
	return m, cmd
}
func (m Model) OnWaitDirName(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if err := m.repo.SaveDir(m.input.input.Value(), m.GetCurrentOrderId()); err != nil {
				println(err.Error())
				return m, tea.Quit
			}
			m.StateId = StandardState
			m.files = m.repo.GetFilesByParentId(m.GetCurrentOrderId())
			return m, cmd
		default:
			m.input.input, cmd = m.input.input.Update(msg)
			return m, cmd

		}
	}
	return m, cmd
}
func (m Model) OnShowFileContent(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "b", "left", "h":
			if len(m.order) > 1 {
				return m.Back(), cmd
			}
		}
	}
	return m, cmd
}
func (m Model) OnDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

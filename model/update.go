package model

import (
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) OnStandardUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "delete":
			if len(m.files) > 0 && m.GetChecked().Id != defaultFirstDirId {
				m.StateId = DeleteState
			}
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
			if m.GetCurrentOrder().IsFolder {
				m.StateId = WaitDirnameState
				return m.SetInput("dirname"), cmd
			}
		case "n":
			if m.GetCurrentOrder().IsFolder {
				m.StateId = WaitFilenameState
				return m.SetInput("filename"), cmd
			}
		case "N":
			if m.GetCurrentOrder().IsFolder {
				m.StateId = WaitFilenameMultiStringState
				return m.SetInput("filename"), cmd
			}
		case "c":
			if !m.GetChecked().IsFolder {
				s, err := m.repo.GetFileContentByFileId(m.GetChecked().Id)
				if err != nil {
					println(err.Error())
				}
				_ = clipboard.WriteAll(s.FileContent)

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
				m.fileContent = c.Viewed()
				return m, cmd
			}
		}
	}
	return m, cmd
}
func (m Model) OnWaitFilenameUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			m.StateId = WaitFileContentState
			m.input.value = m.input.input.Value()
			return m.SetInput("content"), cmd
		case "ctrl+c":
			m.input.value = ""
			m.StateId = StandardState
			return m, cmd
		default:
			m.input.input, cmd = m.input.input.Update(msg)
			return m, cmd
		}
	}
	return m, cmd
}
func (m Model) OnWaitFilenameMultiStringUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "enter":

			m.StateId = WaitMultipleFileContentState
			m.input.value = m.input.input.Value()
			return m.SetArea("content"), cmd
		case "ctrl+c":
			m.input.value = ""
			m.StateId = StandardState
			return m, cmd
		default:
			m.input.input, cmd = m.input.input.Update(msg)
			return m, cmd
		}
	}
	return m, cmd
}
func (m Model) OnWaitFileContentUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "ctrl+c":
			m.input.value = ""
			m.StateId = StandardState
			return m, cmd
		default:
			m.input.input, cmd = m.input.input.Update(msg)
			return m, cmd
		}
	}
	return m, cmd
}
func (m Model) OnWaitMultipleFileContentUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+]":
			if err := m.repo.SaveFileWithContent(m.input.value, m.area.area.Value(), m.GetCurrentOrderId()); err != nil {
				return m, tea.Quit
			}
			m.StateId = StandardState
			m.files = m.repo.GetFilesByParentId(m.GetCurrentOrderId())
			m.input.value = ""
			m.area.value = ""
			return m, cmd
		case "ctrl+c":
			m.StateId = StandardState
			m.input.value = ""
			m.area.value = ""
			return m, cmd
		default:
			if !m.area.area.Focused() {
				cmd = m.area.area.Focus()
			}
			cmds = append(cmds, cmd)
			m.area.area, cmd = m.area.area.Update(msg)
			return m, cmd
		}
	}
	return m, cmd
}
func (m Model) OnWaitDirNameUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.input.value = ""
			return m, cmd
		case "ctrl+c":
			m.input.value = ""
			m.StateId = StandardState
			return m, cmd
		default:
			m.input.input, cmd = m.input.input.Update(msg)
			return m, cmd

		}
	}
	return m, cmd
}
func (m Model) OnShowFileContentUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "b", "left", "h":
			if len(m.order) > 1 {
				return m.Back(), cmd
			}
		//TODO(fix duplicates)
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}
func (m Model) OnDeleteUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "y":
			m.repo.DeleteFolders(m.GetChecked().Id)
			m.files = m.repo.GetFilesByParentId(m.GetCurrentOrderId())
			m.StateId = StandardState
			return m.SetDefaultCursor(), cmd
		case "n":
			m.StateId = StandardState
			return m, cmd
		}
	}
	return m, cmd
}

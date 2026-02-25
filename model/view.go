package model

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) OnStandardView() string {
	list := make([]string, 0)
	for i, item := range m.files {

		cursor := emptyCursor
		colored := false
		if m.cursor == i {
			cursor = filledCursor
			colored = true
		}

		// Render the row
		suffix := m.vars.Icons.File
		if item.IsFolder {
			suffix = m.vars.Icons.Folder
		}
		t := fmt.Sprintf(menuFormat, suffix, cursor, item.Filename)

		if colored {
			list = append(list, Bold().Foreground(m.vars.Colors.Selected).Render(t))
		} else {
			list = append(list, t)
		}

	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.defaultHeader(),
		lipgloss.JoinVertical(lipgloss.Left, list...),
		"",
		m.defaultFooter(true),
	)
}

func (m Model) OnShowFileContentView() string {
	title := Bold().
		Foreground(m.vars.Colors.ViewFileTitle).
		Render(m.fileContent.Filename)

	content := lipgloss.NewStyle().
		Foreground(m.vars.Colors.DefaultTextColor).
		PaddingTop(1).
		PaddingBottom(1).
		Render(m.fileContent.FileContent)

	c := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		m.Separator(50),
		lipgloss.JoinHorizontal(lipgloss.Center, m.vars.Icons.CipherData, " ", content),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.defaultHeader(),
		c,
		m.defaultFooter(false),
	)
}
func (m Model) OnDeleteView() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.defaultHeader(),
		fmt.Sprintf(
			"%s - delete %s, %s - cancel deleting",
			m.DefaultHint("'y'"),
			m.files[m.cursor].Filename,
			m.DefaultHint("'n'")),
		"",
		m.defaultFooter(false),
	)
}
func (m Model) OnRegisterMasterKeyView() string {
	key1 := m.input.input.Value()
	key2 := m.confirmInput.input.Value()

	matchIndicator := ""
	if len(key1) > 0 && len(key2) > 0 {
		if key1 == key2 {
			matchIndicator = lipgloss.NewStyle().
				Foreground(m.vars.Colors.SuccessTextColor).
				Render(KeyMatchMessage)
		} else {
			matchIndicator = lipgloss.NewStyle().
				Foreground(m.vars.Colors.ErrorTextColor).
				Render(KeyDontMatchMessage)
		}
	}

	boxContent := lipgloss.JoinVertical(
		lipgloss.Left,
		Bold().Render(m.vars.Icons.MasterKey, RegisterMasterKey),
		"",
		lipgloss.NewStyle().Foreground(m.vars.Colors.DefaultTextColor).Render(EnterMasterKeyLabel),
		m.input.input.View(),
		"",
		lipgloss.NewStyle().Foreground(m.vars.Colors.DefaultTextColor).Render(ConfirmMasterKeyLable),
		m.confirmInput.input.View(),
		"",
		matchIndicator,
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			fmt.Sprintf("%s to switch", m.DefaultHint("'Tab'")),
			" • ",
			fmt.Sprintf("%s to confirm", m.DefaultHint("'Enter'")),
		),
	)

	boxStyle := RoundedBorder().
		BorderForeground(m.vars.Colors.BoxBorderColor).
		Padding(1, 2).
		Width(50)

	box := boxStyle.Render(boxContent)

	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}
func (m Model) OnWaitMasterKeyView() string {
	errorMsg := ""
	if m.keyError != "" {
		errorMsg = Bold().
			Foreground(m.vars.Colors.ErrorTextColor).
			Render("✗ " + m.keyError)
	}
	boxContent := lipgloss.JoinVertical(
		lipgloss.Left,
		Bold().Render(m.vars.Icons.MasterKey, MasterKeyLabel),
		"",
		m.input.input.View(),
		"",
		errorMsg,
		"",
		fmt.Sprintf("Press %s to unlock", m.DefaultHint("'Enter'")),
	)

	boxStyle := RoundedBorder().
		BorderForeground(m.vars.Colors.BoxBorderColor).
		Padding(1, 2).
		Width(40).
		AlignHorizontal(lipgloss.Center)

	box := boxStyle.Render(boxContent)

	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

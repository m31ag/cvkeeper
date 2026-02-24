package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) defaultHeader() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		Bold().
			MarginLeft(20).
			PaddingLeft(10).
			PaddingRight(10).
			MarginBottom(1).
			AlignHorizontal(lipgloss.Center).
			Background(m.vars.Colors.TitleBackground).
			Render("CVKeeper"),
		Bold().Foreground(m.vars.Colors.History).Render(strings.Join(m.history, "/")),
		lipgloss.NewStyle().
			Foreground(m.vars.Colors.HorizontalSeparator).
			Render(strings.Repeat("─", 70)),
		"",
	)

}
func (m Model) defaultFooter(showHints bool) string {
	hints := []string{}
	hints = append(hints,
		lipgloss.NewStyle().
			Foreground(m.vars.Colors.HorizontalSeparator).
			Render(strings.Repeat("─", 70)),
	)
	if showHints {
		hints = append(hints,

			fmt.Sprintf(
				"Press %s to add single-string, %s to add multiple-string",
				Bold().Foreground(m.vars.Colors.HintKey).Render("'n'"),
				Bold().Foreground(m.vars.Colors.HintKey).Render("'N'"),
			),
			fmt.Sprintf("Press %s to create folder", Bold().Foreground(m.vars.Colors.HintKey).Render("'f'")),
			fmt.Sprintf("Press %s to copy file content", Bold().Foreground(m.vars.Colors.HintKey).Render("'c'")),
			fmt.Sprintf("Press %s to delete file/folder", Bold().Foreground(m.vars.Colors.HintKey).Render("'d'")),
			"",
		)
	}
	hints = append(hints, fmt.Sprintf("Press %s to quit", Bold().Foreground(m.vars.Colors.HintKey).Render("'q'")))
	return lipgloss.JoinVertical(lipgloss.Left, hints...)
}
func (m Model) inputFooter() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		lipgloss.NewStyle().
			Foreground(m.vars.Colors.HorizontalSeparator).
			Render(strings.Repeat("─", 70)),

		fmt.Sprintf("Press %s to cancel",
			Bold().Foreground(m.vars.Colors.HintKey).Render("'ctrl+c'")),
		"",
		fmt.Sprintf("Press %s to quit.", Bold().Foreground(m.vars.Colors.HintKey).Render("'q'")),
	)

}
func (m Model) areaFooter() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		lipgloss.NewStyle().
			Foreground(m.vars.Colors.HorizontalSeparator).
			Render(strings.Repeat("─", 70)),

		fmt.Sprintf("Press %s to cancel",
			Bold().Foreground(m.vars.Colors.HintKey).Render("'ctrl+c'")),
		fmt.Sprintf("Press %s to save",
			Bold().Foreground(m.vars.Colors.HintKey).Render("'ctrl+]'")),
		"",
		fmt.Sprintf("Press %s to quit.", Bold().Foreground(m.vars.Colors.HintKey).Render("'q'")),
	)

}
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
		suffix := "\U0001F4C4"
		if item.IsFolder {
			suffix = "\U0001F4C1"
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

	separator := lipgloss.NewStyle().
		Foreground(m.vars.Colors.HorizontalSeparator).
		Render(strings.Repeat("─", 50))

	content := lipgloss.NewStyle().
		Foreground(m.vars.Colors.DefaultTextColor).
		PaddingTop(1).
		PaddingBottom(1).
		Render(m.fileContent.FileContent)

	c := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		separator,
		lipgloss.JoinHorizontal(lipgloss.Center, "🔑 ", content),
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
			Bold().Foreground(m.vars.Colors.HintKey).Render("'y'"),
			m.files[m.cursor].Filename,
			Bold().Foreground(m.vars.Colors.HintKey).Render("'n'")),
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
				Render("✓ Keys match")
		} else {
			matchIndicator = lipgloss.NewStyle().
				Foreground(m.vars.Colors.ErrorTextColor).
				Render("✗ Keys don't match")
		}
	}

	boxContent := lipgloss.JoinVertical(
		lipgloss.Left,
		Bold().Render("🔐 Register Master Key"),
		"",
		Faint().Render("Enter key:"),
		m.input.input.View(),
		"",
		Faint().Render("Confirm key:"),
		m.confirmInput.input.View(),
		"",
		matchIndicator,
		"",
		lipgloss.NewStyle().
			Foreground(m.vars.Colors.HintKey).
			Render("Tab to switch • Enter to confirm"),
	)

	boxStyle := RoundedBorder().
		BorderForeground(m.vars.Colors.BoxBorderColor).
		Padding(1, 2).
		Width(50)

	box := boxStyle.Render(boxContent)

	width := m.termWidth
	height := m.termHeight
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 24
	}

	return lipgloss.Place(
		width,
		height,
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
		Bold().Render("🔐 Master Key"),
		"",
		m.input.input.View(),
		"",
		errorMsg,
		"",
		lipgloss.NewStyle().
			Foreground(m.vars.Colors.HintKey).
			Render("Press Enter to unlock"),
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
func (m Model) DefaultInputView() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.defaultHeader(),
		m.input.input.View(),
		m.inputFooter(),
	)

}
func (m Model) DefaultAreaView() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.defaultHeader(),
		m.area.area.View(),
		m.areaFooter(),
	)

}

package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) defaultHeader() string {
	return fmt.Sprint(
		fmt.Sprintf(
			"%s\n\n",
			lipgloss.
				NewStyle().
				MarginLeft(20).
				PaddingLeft(10).
				PaddingRight(10).
				AlignHorizontal(lipgloss.Center).
				Bold(true).
				Background(lipgloss.Color("#5f5fff")).
				Render("CVKeeper")),
		//TODO (change lipcloss.NewStyle to var in model or singleton)
		render(strings.Join(m.history, "/"), historyFormat, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(m.vars.Colors.History))),
	)
}
func (m Model) defaultFooter(showHints bool) string {
	s := ""
	if showHints {
		s += fmt.Sprintf("\nPress %s to add single-string, %s to add multiple-string\n"+
			"Press %s to create folder\n"+
			"Press %s to copy file content\n"+
			"Press %s to delete file/folder\n",
			m.styleAndRender("'n'", true, m.vars.Colors.HintKey),
			m.styleAndRender("'N'", true, m.vars.Colors.HintKey),
			m.styleAndRender("'f'", true, m.vars.Colors.HintKey),
			m.styleAndRender("'c'", true, m.vars.Colors.HintKey),
			m.styleAndRender("'d'", true, m.vars.Colors.HintKey),
		)
	}
	return fmt.Sprint(
		s,
		fmt.Sprintf("\nPress %s to quit.\n", m.styleAndRender("'q'", true, m.vars.Colors.HintKey)),
	)
}
func (m Model) inputFooter() string {
	s := fmt.Sprintf("\nPress %s to cancel\n",
		m.styleAndRender("'ctrl+c'", true, m.vars.Colors.HintKey),
	)
	return fmt.Sprint(
		s,
		fmt.Sprintf("\nPress %s to quit.\n", m.styleAndRender("'q'", true, m.vars.Colors.HintKey)),
	)
}
func (m Model) areaFooter() string {
	s := fmt.Sprintf("\nPress %s to cancel\n"+
		"Press %s to save\n",
		m.styleAndRender("'ctrl+c'", true, m.vars.Colors.HintKey),
		m.styleAndRender("'ctrl+]'", true, m.vars.Colors.HintKey),
	)
	return fmt.Sprint(
		s,
		fmt.Sprintf("\nPress %s to quit.\n", m.styleAndRender("'q'", true, m.vars.Colors.HintKey)),
	)
}
func (m Model) OnStandardView() string {
	list := ""
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
		list += m.showItem(fmt.Sprintf(menuFormat, suffix, cursor, item.Filename), colored)
	}
	return fmt.Sprint(m.defaultHeader(), list, m.defaultFooter(true))
}

func (m Model) OnShowFileContentView() string {
	return fmt.Sprint(
		m.defaultHeader(),
		"\n",
		m.styleAndRender(m.fileContent, true, ""),
		"\n",
		m.defaultFooter(false),
	)
}
func (m Model) OnDeleteView() string {
	return fmt.Sprint(
		m.defaultHeader(),
		m.styleAndRender(strings.Repeat("#", 40)+"\n", true, ""),

		fmt.Sprintf(
			"%s - delete %s, %s - cancel deleting\n",
			m.styleAndRender("'y'", true, m.vars.Colors.HintKey),
			m.files[m.cursor].Filename,
			m.styleAndRender("'n'", true, m.vars.Colors.HintKey)),

		m.styleAndRender(strings.Repeat("#", 40)+"\n", true, ""),

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
				Foreground(lipgloss.Color(m.vars.Colors.SuccessTextColor)).
				Render("✓ Keys match")
		} else {
			matchIndicator = lipgloss.NewStyle().
				Foreground(lipgloss.Color(m.vars.Colors.ErrorTextColor)).
				Render("✗ Keys don't match")
		}
	}

	boxContent := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render("🔐 Register Master Key"),
		"",
		lipgloss.NewStyle().Faint(true).Render("Enter key:"),
		m.input.input.View(),
		"",
		lipgloss.NewStyle().Faint(true).Render("Confirm key:"),
		m.confirmInput.input.View(),
		"",
		matchIndicator,
		"",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.vars.Colors.HintKey)).
			Render("Tab to switch • Enter to confirm"),
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.vars.Colors.BoxBorderColor)).
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
		errorMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.vars.Colors.ErrorTextColor)).
			Bold(true).
			Render("✗ " + m.keyError)
	}
	boxContent := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render("🔐 Master Key"),
		"",
		m.input.input.View(),
		"",
		errorMsg,
		"",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.vars.Colors.HintKey)).
			Render("Press Enter to unlock"),
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.vars.Colors.BoxBorderColor)).
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
	return fmt.Sprint(
		m.defaultHeader(),
		m.input.input.View(),
		m.inputFooter(),
	)

}
func (m Model) DefaultAreaView() string {
	return fmt.Sprint(
		m.defaultHeader(),
		m.area.area.View(),
		m.areaFooter(),
	)

}
func (m Model) showItem(txt string, colored bool) string {

	if colored {
		//maybe make const
		return strings.TrimSpace(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(m.vars.Colors.Selected)).Render(txt))
	}

	return txt

}
func (m Model) styleAndRender(t string, bold bool, color lipgloss.Color) string {
	if len(color) == 0 {
		color = m.vars.Colors.DefaultTextColor
	}
	s := lipgloss.NewStyle().Bold(bold).Foreground(lipgloss.Color(color))

	return strings.TrimSpace(s.Render(t))
}
func render(txt, format string, style lipgloss.Style) string {
	return strings.TrimSpace(style.Render(fmt.Sprintf(format, txt)))
}

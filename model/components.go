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
			Render(AppName),
		Bold().Foreground(m.vars.Colors.History).Render(strings.Join(m.history, "/")),
		m.Separator(70),
		"",
	)

}
func (m Model) defaultFooter(showHints bool) string {
	hints := []string{}
	hints = append(hints, m.Separator(70))
	if showHints {
		hints = append(hints,

			fmt.Sprintf(
				"Press %s to add single-string, %s to add multiple-string",
				m.DefaultHint("'n'"),
				m.DefaultHint("'N'"),
			),
			fmt.Sprintf("Press %s to create folder", m.DefaultHint("'f'")),
			fmt.Sprintf("Press %s to copy file content", m.DefaultHint("'c'")),
			fmt.Sprintf("Press %s to delete file/folder", m.DefaultHint("'d'")),
			"",
		)
	}
	hints = append(hints, fmt.Sprintf("Press %s to quit", m.DefaultHint("'q'")))
	return lipgloss.JoinVertical(lipgloss.Left, hints...)
}
func (m Model) inputFooter() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		m.Separator(70),

		fmt.Sprintf("Press %s to cancel",
			m.DefaultHint("'ctrl+c'")),
		"",
		fmt.Sprintf("Press %s to quit.", m.DefaultHint("'q'")),
	)

}
func (m Model) areaFooter() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		m.Separator(70),
		fmt.Sprintf("Press %s to cancel",
			m.DefaultHint("'ctrl+c'")),
		fmt.Sprintf("Press %s to save",
			m.DefaultHint("'ctrl+]'")),
		"",
		fmt.Sprintf("Press %s to quit.", m.DefaultHint("'q'")),
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

func (m Model) DefaultHint(h string) string {
	return Bold().Foreground(m.vars.Colors.HintKey).Render(h)
}
func (m Model) Separator(l int) string {
	return lipgloss.NewStyle().
		Foreground(m.vars.Colors.HorizontalSeparator).
		Render(strings.Repeat("─", l))

}

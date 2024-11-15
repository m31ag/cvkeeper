package model

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
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
		render(strings.Join(m.history, "/"), historyFormat, historyStyle),
	)
}
func (m Model) defaultFooter(showHints bool) string {
	s := ""
	if showHints {
		s += fmt.Sprintf("\nPress %s to add one-string, %s to add multiple-string\nPress %s to create folder",
			styleAndRender("'n'", true, purpleColor),
			styleAndRender("'N'", true, purpleColor),
			styleAndRender("'f'", true, purpleColor),
		)
	}
	return fmt.Sprint(
		s,
		fmt.Sprintf("\nPress %s to quit.\n", styleAndRender("'q'", true, purpleColor)),
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
		list += showItem(fmt.Sprintf(menuFormat, suffix, cursor, item.Filename), colored)
	}
	return fmt.Sprint(m.defaultHeader(), list, m.defaultFooter(true))
}

func (m Model) OnShowFileContentView() string {
	return fmt.Sprint(
		m.defaultHeader(),
		styleAndRender(m.fileContent, true, whiteColor),
		"\n",
		m.defaultFooter(false),
	)
}
func (m Model) OnDeleteView() string {
	return fmt.Sprint(
		m.defaultHeader(),
		styleAndRender(strings.Repeat("#", 40)+"\n", true, whiteColor),

		fmt.Sprintf(
			"%s - delete %s, %s - cancel deleting\n",
			styleAndRender("'y'", true, purpleColor),
			m.files[m.cursor].Filename,
			styleAndRender("'n'", true, purpleColor)),

		styleAndRender(strings.Repeat("#", 40)+"\n", true, whiteColor),

		m.defaultFooter(false),
	)
}
func (m Model) DefaultContentView() string {
	return fmt.Sprint(
		m.defaultHeader(),
		m.input.input.View(),
		m.defaultFooter(true),
	)

}
func showItem(txt string, colored bool) string {

	if colored {
		return strings.TrimSpace(style.Render(txt))
	} else {
		return txt
	}

}
func styleAndRender(t string, bold bool, color string) string {
	if len(color) == 0 {
		color = whiteColor
	}
	s := lipgloss.NewStyle().Bold(bold).Foreground(lipgloss.Color(color))

	return strings.TrimSpace(s.Render(t))
}
func render(txt, format string, style lipgloss.Style) string {
	return strings.TrimSpace(style.Render(fmt.Sprintf(format, txt)))
}

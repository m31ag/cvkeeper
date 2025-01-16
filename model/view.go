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
		s += fmt.Sprintf("\nPress %s to add one-string, %s to add multiple-string\n"+
			"Press %s to create folder\n"+
			"Press %s to copy file content\n"+
			"Press %s to delete file/folder\n",
			styleAndRender("'n'", true, purpleColor),
			styleAndRender("'N'", true, purpleColor),
			styleAndRender("'f'", true, purpleColor),
			styleAndRender("'c'", true, purpleColor),
			styleAndRender("'d'", true, purpleColor),
		)
	}
	return fmt.Sprint(
		s,
		fmt.Sprintf("\nPress %s to quit.\n", styleAndRender("'q'", true, purpleColor)),
	)
}
func (m Model) inputFooter() string {
	s := fmt.Sprintf("\nPress %s to cancel\n",
		styleAndRender("ctrl+c'", true, purpleColor),
	)
	return fmt.Sprint(
		s,
		fmt.Sprintf("\nPress %s to quit.\n", styleAndRender("'q'", true, purpleColor)),
	)
}
func (m Model) areaFooter() string {
	s := fmt.Sprintf("\nPress %s to cancel\n"+
		"Press %s to save\n",
		styleAndRender("ctrl+c'", true, purpleColor),
		styleAndRender("ctrl+]'", true, purpleColor),
	)
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
		"\n",
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

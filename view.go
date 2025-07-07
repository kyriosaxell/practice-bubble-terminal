package main

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

const gap = "\n\n"

var (
	appNameStyle        = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)
	faintStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Faint(true)
	listEnumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
)

// View genera la representación visual del modelo actual.
func (m Model) View() string {
	s := appNameStyle.Render("NOTES APP") + "\n\n"

	if m.state == titleView {
		s += "Título de la nota:" + gap
		s += m.textInput.View() + gap
		s += faintStyle.Render("enter - guardar • ctrl+q - descartar")
	}

	if m.state == bodyView {
		s += "Note:" + gap
		s += m.textArea.View() + gap
		s += faintStyle.Render("ctrl+s - guardar • ctrl+q - descartar")
	}

	if m.state == listView {
		for i, n := range m.notes {
			prefix := " "
			if i == m.listIndex {
				prefix = "> "
			}
			shortBody := strings.ReplaceAll(n.Body, "\n", "")
			if len(shortBody) > 30 {
				shortBody = shortBody[:30]
			}
			s += listEnumeratorStyle.Render(prefix) + n.Title + " | " + faintStyle.Render(shortBody) + gap
		}
		s += faintStyle.Render("n - nueva nota • q - quit")
	}
	return s
}

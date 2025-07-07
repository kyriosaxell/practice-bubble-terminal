package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
)

type editorFinishedMsg struct {
	err error
}
type model struct {
	altScreenActive bool
	err             error
}

func (m model) Init() tea.Cmd {
	return nil
}

func openEditor() tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	c := exec.Command(editor)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err: err}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.altScreenActive = !m.altScreenActive
			cmd := tea.EnterAltScreen
			if !m.altScreenActive {
				cmd = tea.ExitAltScreen
			}
			return m, cmd
		case "e":
			return m, openEditor()
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}
	return "Press 'e' to open your EDITOR.\nPress 'a' to toggle the altScreen\nPress 'q' to quit.\n"
}

func main() {
	m := model{}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

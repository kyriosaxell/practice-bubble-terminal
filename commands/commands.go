package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"net/http"
	"os"
	"time"
)

const url = "https://charm.sh/"

var outputStyle = lipgloss.NewStyle().Background(lipgloss.Color("#008000")).Bold(true).Render

type model struct {
	status int
	err    error
}

type errMsg struct {
	err error
}

type statusMsg int

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}

// checkServer crea un cliente HTTP y hace una petición GET.
func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)

	if err != nil {
		return errMsg{err: err}
	}

	return statusMsg(res.StatusCode)
}

// Para los mensajes que contienen errores, a menudo resulta útil implementar
// también la interfaz de error en el mensaje.
func (e errMsg) Error() string { return e.err.Error() }

// Init es el entry point de la app.
func (m model) Init() tea.Cmd {
	return checkServer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		m.status = int(msg)
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	// Se le informa al usuario lo que estamos haciendo.
	s := fmt.Sprintf("Checking URL: %s ...", url)
	// Cuando el server responsa con un status, se agrega a la linea actual.
	if m.status > 0 {
		s += "\n\n"
		s += outputStyle(fmt.Sprintf("%d %s", m.status, http.StatusText(m.status)))
	}

	return "\n" + s + "\n\n"
}

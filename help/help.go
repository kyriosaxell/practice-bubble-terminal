package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

var outputTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#0589ea")).Bold(true)
var helpViewStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#c3a62e")).Bold(true)

// keyMap define un set de keybindings.
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Help  key.Binding
	Hide  key.Binding
	Quit  key.Binding
}

// ShortHelp regresa keybindings para ser mostrados en la mini vista de ayuda. Es parte
// del key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp regresa keybindings para la vista expandida. Es parte
// del key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // Primera columna
		{k.Hide, k.Quit},                // Segunda columna
	}
}

// keys almacena el keyMap con combinaciones de teclas predefinidas y sus correspondientes descripciones de ayuda.
var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"), // Binding con teclas flecha arriba y 'k' estilo vim.
		key.WithHelp("↑/k", "hacia arriba"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "hacia abajo"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "hacia la izquierda"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "hacia la derecha"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "muestra la ayuda"),
	),
	Hide: key.NewBinding(
		key.WithKeys("o", "O"),
		key.WithHelp("o/O", "oculta la ayuda"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "salir"),
	),
}

type model struct {
	keys             keyMap
	help             help.Model
	expandedHelpView bool
	inputStyle       lipgloss.Style
	helpViewStyle    lipgloss.Style
	lastKey          string
	quitting         bool
}

// newModel regresa un modelo
func newModel() model {

	// Estilos para los textos del menú de ayuda.
	h := help.New()
	h.Styles = help.Styles{
		ShortKey:  helpViewStyle,
		ShortDesc: helpViewStyle,
	}

	return model{
		keys: keys,
		help: h,
		inputStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#063970")).
			Padding(0, 1),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Si establecemos un ancho en el menú de ayuda, este puede truncarse con elegancia
		// según sea necesario.
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.lastKey = "↑"
		case key.Matches(msg, m.keys.Down):
			m.lastKey = "↓"
		case key.Matches(msg, m.keys.Left):
			m.lastKey = "←"
		case key.Matches(msg, m.keys.Right):
			m.lastKey = "→"
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

			m.keys.Help.SetEnabled(false)
			m.keys.Hide.SetEnabled(true)
		case key.Matches(msg, m.keys.Hide):
			m.help.ShowAll = !m.help.ShowAll

			m.keys.Help.SetEnabled(true)
			m.keys.Hide.SetEnabled(false)
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Adiós\n"
	}

	var status string
	if m.lastKey == "" {
		status = "Esperando por teclas..."
	} else {
		status = "Tecleaste: " + m.inputStyle.Render(m.lastKey)
	}

	vp := viewport.New(130, 5)
	vp.SetContent(
		outputTitleStyle.Width(vp.Width).Render("\nBienvenido a la ayuda de Charm! Brace yourself 🥸," +
			" esto es un poco complicado. \n\n"),
	)

	helpView := m.help.View(m.keys)
	styledHelpView := helpViewStyle.Render(helpView)
	height := 4 - strings.Count(status, "\n") - strings.Count(helpView, "\n")
	return vp.View() +
		"\n" + status +
		"\n" + strings.Repeat("\n", height) +
		styledHelpView // Textos del menú de ayuda
}

func main() {
	if os.Getenv("HELP_DEBUG") != "" {
		f, err := tea.LogToFile("debug.log", "help")
		if err != nil {
			fmt.Println("No se pudo abrir el archivo para logging:", err)
			os.Exit(1)
		}

		defer func(f *os.File) {
			_ = f.Close()
		}(f)
	}

	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Printf("No se pudo iniciar el programa: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"log"
)

const (
	listView uint = iota
	titleView
	bodyView
)

type Model struct {
	state     uint
	store     *Store
	notes     []Note
	currNote  Note
	listIndex int
	textArea  textarea.Model
	textInput textinput.Model
}

func NewModel(store *Store) Model {
	notes, err := store.GetNotes()
	if err != nil {
		log.Fatalf("no se pudieron obtener las notas: %s", err)
	}

	return Model{state: listView, store: store, notes: notes, textArea: textarea.New(), textInput: textinput.New()}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.textArea, cmd = m.textArea.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String() //up, down, enter, esc, etc.
		switch m.state {
		// En caso de que el estado sea listView, se muestran las notas.
		case listView:
			switch key {
			case "q":
				return m, tea.Quit
			case "n": // Nueva nota
				m.textInput.SetValue("")
				m.textInput.Placeholder = "Escribe el título de la nota"
				m.textInput.Focus()
				m.textInput.CharLimit = 100
				m.state = titleView
			case tea.KeyUp.String(), "k":
				if m.listIndex > 0 {
					m.listIndex--
				}
			case tea.KeyDown.String(), "j":
				if m.listIndex < len(m.notes)-1 {
					m.listIndex++
				}
			case tea.KeyEnter.String():
				m.currNote = m.notes[m.listIndex]
				m.textArea.SetValue(m.currNote.Body)
				m.textArea.Focus()
				m.textArea.CursorEnd()
				m.state = bodyView
			}
		case titleView:
			switch key {
			case tea.KeyEnter.String():
				title := m.textInput.Value()
				if title != "" {
					m.currNote.Title = title
					m.textArea.SetValue("")
					m.textArea.Focus()
					m.textArea.CursorEnd()

					m.state = bodyView
				}
			case tea.KeyCtrlQ.String(): // Por alguna razón 'esc' no funciona con la terminal de Jetbrains y Alacritty.
				m.state = listView
			}
		case bodyView:
			switch key {
			case tea.KeyCtrlS.String():
				body := m.textArea.Value()
				m.currNote.Body = body

				var err error
				if err = m.store.SaveNote(m.currNote); err != nil {
					return m, tea.Quit
				}

				m.notes, err = m.store.GetNotes()
				if err != nil {
					return m, tea.Quit
				}
				m.currNote = Note{} // Reinicia la nota actual
				m.state = listView
			case tea.KeyCtrlQ.String():
				m.state = listView
			}
		}

	}
	// Los comandos cmds se pasan empaquetados para ser ejecutados simultáneamente
	return m, tea.Batch(cmds...)
}

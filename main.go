package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	// Se crea una instancia del Store DB
	store := &Store{}

	if err := store.Init(); err != nil {
		log.Printf("Error initializing store: %s", err)
	}

	m := NewModel(store)

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

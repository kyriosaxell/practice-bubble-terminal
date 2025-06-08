package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Note struct {
	ID    int64
	Title string
	Body  string
}

type Store struct {
	conn *sql.DB
}

func (s *Store) Init() error {
	var err error
	s.conn, err = sql.Open("sqlite3", "./notes_db.db")
	if err != nil {
		return err
	}

	createTableStmt := `CREATE TABLE IF NOT EXISTS notes (
    id integer not null primary key,
    title text not null,
    body text not null
    );`
	if _, err = s.conn.Exec(createTableStmt); err != nil {
		return err
	}
	return nil
}

func (s *Store) GetNotes() ([]Note, error) {
	rows, err := s.conn.Query("SELECT * FROM notes")
	if err != nil {
		return nil, err
	}
	// Se coloca inmediatamente después de abrir el recurso para evitar fugas de memoria.
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, &note.Title, &note.Body)
		if err != nil {
			return nil, fmt.Errorf("error al escanear la consulta: %w", err)
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (s *Store) SaveNote(note Note) error {
	if note.ID == 0 {
		note.ID = time.Now().UTC().UnixNano()
	}

	insertQuery := `INSERT INTO notes (id, title, body) VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET title=excluded.title, body=excluded.body;`

	if _, err := s.conn.Exec(insertQuery, note.ID, note.Title, note.Body); err != nil {
		return fmt.Errorf("error al insertar o actualizar la nota: %w", err)
	}

	return nil
}

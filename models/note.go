package models

import (
	"encoding/json"
	"net/http"
)

type Note struct {
	Content string `json:"content"`
}

// Implement render.Binder for Note
func (n *Note) Bind(r *http.Request) error {
	type NoteWithFields struct {
		Content string `json:"content"`
	}

	var temp NoteWithFields
	if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
		return err
	}

	n.Content = temp.Content

	return nil
}

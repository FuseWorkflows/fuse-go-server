package models

import (
	"encoding/json"
	"net/http"
)

type AISuggestions struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Chapters    []string `json:"chapters"`
	Thumbnail   string   `json:"thumbnail"`
	// prompt
	// seed
}

// Implement render.Binder for AISuggestions
func (a *AISuggestions) Bind(r *http.Request) error {
	// Decode the JSON body into a temporary struct that includes all fields
	type AISuggestionsWithFields struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Keywords    []string `json:"keywords"`
		Chapters    []string `json:"chapters"`
		Thumbnail   string   `json:"thumbnail"`
	}

	var temp AISuggestionsWithFields
	if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
		return err
	}

	// Copy the data from the temporary struct to the AISuggestions struct
	a.Title = temp.Title
	a.Description = temp.Description
	a.Keywords = temp.Keywords
	a.Chapters = temp.Chapters
	a.Thumbnail = temp.Thumbnail

	return nil
}

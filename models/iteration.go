package models

import (
	"encoding/json"
	"net/http"
)

type Iteration struct {
	ID        string          `json:"id"`
	Video     Video           `json:"video"`
	URL       string          `json:"url"`
	Length    string          `json:"length"`
	Status    IterationStatus `json:"status"`
	Notes     string          `json:"notes"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
	//createdby  Editor
}

// Implement render.Binder for Iteration
func (i *Iteration) Bind(r *http.Request) error {
	// Decode the JSON body into a temporary struct that includes Video
	type IterationWithVideo struct {
		ID        string          `json:"id"`
		Video     Video           `json:"video"`
		URL       string          `json:"url"`
		Length    string          `json:"length"`
		Status    IterationStatus `json:"status"`
		Notes     string          `json:"notes"`
		CreatedAt string          `json:"createdAt"`
		UpdatedAt string          `json:"updatedAt"`
	}

	var temp IterationWithVideo
	if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
		return err
	}

	// Copy the data from the temporary struct to the Iteration struct
	i.ID = temp.ID
	i.Video = temp.Video
	i.URL = temp.URL
	i.Length = temp.Length
	i.Status = temp.Status
	i.Notes = temp.Notes
	i.CreatedAt = temp.CreatedAt
	i.UpdatedAt = temp.UpdatedAt

	return nil
}

type IterationStatus string

const (
	Processing IterationStatus = "processing"
	Completed  IterationStatus = "completed"
	Failed     IterationStatus = "failed"
)

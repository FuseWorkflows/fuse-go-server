package models

import (
	"encoding/json"
	"net/http"
)

type Editor struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Tier      Tier   `json:"tier"`
	Trial     bool   `json:"trial"`
}

// Implement render.Binder for Editor
func (e *Editor) Bind(r *http.Request) error {
	// Decode the JSON body into a temporary struct that includes CreatedAt and UpdatedAt
	type EditorWithTimestamps struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
		Tier      Tier   `json:"tier"`
		Trial     bool   `json:"trial"`
	}

	var temp EditorWithTimestamps
	if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
		return err
	}

	// Copy the data from the temporary struct to the Editor struct
	e.ID = temp.ID
	e.Username = temp.Username
	e.Email = temp.Email
	e.Password = temp.Password
	e.CreatedAt = temp.CreatedAt
	e.UpdatedAt = temp.UpdatedAt
	e.Tier = temp.Tier
	e.Trial = temp.Trial

	return nil
}

type Tier string

const (
	Premium Tier = "premium"
	Basic   Tier = "basic"
	Free    Tier = "free"
)

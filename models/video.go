package models

import (
	"encoding/json"
	"net/http"
)

type Video struct {
	ID            string      `json:"id"`
	Status        Status      `json:"status"`
	Iterations    []Iteration `json:"iterations"`
	Resources     string      `json:"resources"`
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Keywords      []string    `json:"keywords"`
	Category      string      `json:"category"`
	PrivacyStatus bool        `json:"privacyStatus"`
	Channel       Channel     `json:"channel"`
	Editors       []Editor    `json:"editors"`
	CreatedAt     string      `json:"createdAt"`
	UpdatedAt     string      `json:"updatedAt"`
}

func (v *Video) MarshalJSON() ([]byte, error) {
	type VideoAlias Video
	return json.Marshal(&struct {
		*VideoAlias
		Channel *Channel `json:"channel"`
		Editors []Editor `json:"editors"`
	}{
		VideoAlias: (*VideoAlias)(v),
		Channel:    &v.Channel,
		Editors:    v.Editors,
	})
}

func (v *Video) UnmarshalJSON(data []byte) error {
	type VideoAlias Video
	aux := &struct {
		*VideoAlias
		Channel *Channel `json:"channel"`
		Editors []Editor `json:"editors"`
	}{
		VideoAlias: (*VideoAlias)(v),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	v.Channel = *aux.Channel
	v.Editors = aux.Editors
	return nil
}

// Implement render.Binder for Video
func (v *Video) Bind(r *http.Request) error {
	// Decode the JSON body into a temporary struct that includes Channel and Editors
	type VideoWithChannelAndEditors struct {
		ID            string      `json:"id"`
		Status        Status      `json:"status"`
		Iterations    []Iteration `json:"iterations"`
		Resources     string      `json:"resources"`
		Title         string      `json:"title"`
		Description   string      `json:"description"`
		Keywords      []string    `json:"keywords"`
		Category      string      `json:"category"`
		PrivacyStatus bool        `json:"privacyStatus"`
		Channel       Channel     `json:"channel"`
		Editors       []Editor    `json:"editors"`
		CreatedAt     string      `json:"createdAt"`
		UpdatedAt     string      `json:"updatedAt"`
	}

	var temp VideoWithChannelAndEditors
	if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
		return err
	}

	// Copy the data from the temporary struct to the Video struct
	v.ID = temp.ID
	v.Status = temp.Status
	v.Iterations = temp.Iterations
	v.Resources = temp.Resources
	v.Title = temp.Title
	v.Description = temp.Description
	v.Keywords = temp.Keywords
	v.Category = temp.Category
	v.PrivacyStatus = temp.PrivacyStatus
	v.Channel = temp.Channel
	v.Editors = temp.Editors
	v.CreatedAt = temp.CreatedAt
	v.UpdatedAt = temp.UpdatedAt

	return nil
}

type Status string

const (
	Pending   Status = "pending"
	Published Status = "published"
	Draft     Status = "draft"
)

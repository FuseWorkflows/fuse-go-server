package models

import "encoding/json"

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

type Status string

const (
	Pending   Status = "pending"
	Published Status = "published"
	Draft     Status = "draft"
)

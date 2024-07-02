package models

import "encoding/json"

type Channel struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	API_KEY string  `json:"api_key"`
	Owner   User    `json:"owner"`
	Videos  []Video `json:"videos"`
}

func (c *Channel) MarshalJSON() ([]byte, error) {
	type ChannelAlias Channel
	return json.Marshal(&struct {
		*ChannelAlias
		Owner *User `json:"owner"`
	}{
		ChannelAlias: (*ChannelAlias)(c),
		Owner:        &c.Owner,
	})
}

func (c *Channel) UnmarshalJSON(data []byte) error {
	type ChannelAlias Channel
	aux := &struct {
		*ChannelAlias
		Owner *User `json:"owner"`
	}{
		ChannelAlias: (*ChannelAlias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	c.Owner = *aux.Owner
	return nil
}

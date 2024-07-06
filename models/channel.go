package models

import (
	"encoding/json"
	"net/http"
)

type Channel struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	API_KEY   string  `json:"api_key"`
	Owner     User    `json:"owner"`
	Videos    []Video `json:"videos"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

// func (c *Channel) MarshalJSON() ([]byte, error) {
// 	type ChannelAlias Channel
// 	return json.Marshal(&struct {
// 		*ChannelAlias
// 		Owner *User `json:"owner"`
// 	}{
// 		ChannelAlias: (*ChannelAlias)(c),
// 		Owner:        &c.Owner,
// 	})
// }

// func (c *Channel) UnmarshalJSON(data []byte) error {
// 	type ChannelAlias Channel
// 	aux := &struct {
// 		*ChannelAlias
// 		Owner *User `json:"owner"`
// 	}{
// 		ChannelAlias: (*ChannelAlias)(c),
// 	}

// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}

// 	c.Owner = *aux.Owner
// 	return nil
// }

// Implement render.Binder for Channel
func (c *Channel) Bind(r *http.Request) error {
	// Decode the JSON body into a temporary struct that includes Owner
	type ChannelWithOwner struct {
		ID        string  `json:"id"`
		Name      string  `json:"name"`
		API_KEY   string  `json:"api_key"`
		Owner     User    `json:"owner"`
		Videos    []Video `json:"videos"`
		CreatedAt string  `json:"createdAt"`
		UpdatedAt string  `json:"updatedAt"`
	}

	var temp ChannelWithOwner
	if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
		return err
	}

	// Copy the data from the temporary struct to the Channel struct
	c.ID = temp.ID
	c.Name = temp.Name
	c.API_KEY = temp.API_KEY
	c.Owner = temp.Owner
	c.Videos = temp.Videos
	c.CreatedAt = temp.CreatedAt
	c.UpdatedAt = temp.UpdatedAt

	return nil
}

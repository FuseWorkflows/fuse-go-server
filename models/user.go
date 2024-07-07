package models

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
	Tier      Tier      `json:"tier"`
	Trial     bool      `json:"trial"`
	Channels  []Channel `json:"channels"`
	// Editors []
}

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	aux := &struct {
		*Alias
		Password string    `json:"password,omitempty"`
		Channels []Channel `json:"channels,omitempty"`
	}{
		Alias: (*Alias)(u),
	}
	return json.Marshal(aux)
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

func (u *User) Bind(r *http.Request) error {
	// Decode the JSON body into a temporary struct that includes Channels
	type UserWithChannels struct {
		ID        string    `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		CreatedAt string    `json:"createdAt"`
		UpdatedAt string    `json:"updatedAt"`
		Tier      Tier      `json:"tier"`
		Trial     bool      `json:"trial"`
		Channels  []Channel `json:"channels"`
	}

	var temp UserWithChannels
	if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
		return err
	}

	// Copy the data from the temporary struct to the User struct
	u.ID = temp.ID
	u.Username = temp.Username
	u.Email = temp.Email
	u.Password = temp.Password
	u.CreatedAt = temp.CreatedAt
	u.UpdatedAt = temp.UpdatedAt
	u.Tier = temp.Tier
	u.Trial = temp.Trial
	u.Channels = temp.Channels

	return nil
}

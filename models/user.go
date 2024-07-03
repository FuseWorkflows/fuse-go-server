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
}

func (u *User) MarshalJSON() ([]byte, error) {
	type UserAlias User
	return json.Marshal(&struct {
		*UserAlias
		Channels []Channel `json:"channels"`
	}{
		UserAlias: (*UserAlias)(u),
		Channels:  u.Channels,
	})
}

func (u *User) UnmarshalJSON(data []byte) error {
	type UserAlias User
	aux := &struct {
		*UserAlias
		Channels []Channel `json:"channels"`
	}{
		UserAlias: (*UserAlias)(u),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	u.Channels = aux.Channels
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

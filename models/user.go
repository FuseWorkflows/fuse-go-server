package models

import "encoding/json"

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

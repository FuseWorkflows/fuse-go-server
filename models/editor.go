package models

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

type Tier string

const (
	Premium Tier = "premium"
	Basic   Tier = "basic"
	Free    Tier = "free"
)

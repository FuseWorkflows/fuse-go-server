package models

type AISuggestions struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Chapters    []string `json:"chapters"`
	Thumbnail   string   `json:"thumbnail"`
}

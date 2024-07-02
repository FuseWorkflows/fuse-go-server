package models

type Iteration struct {
	ID        string          `json:"id"`
	URL       string          `json:"url"`
	Length    string          `json:"length"`
	Status    IterationStatus `json:"status"`
	Notes     string          `json:"notes"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

type IterationStatus string

const (
	Processing IterationStatus = "processing"
	Completed  IterationStatus = "completed"
	Failed     IterationStatus = "failed"
)

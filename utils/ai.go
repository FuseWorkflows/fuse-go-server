package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/your-username/youtube-backend/config"
	"github.com/your-username/youtube-backend/models"
)

// GetAISuggestions retrieves AI suggestions for video metadata
func GetAISuggestions(video *models.Video, cfg *config.Config) (*models.AISuggestions, error) {
	// Define AI request payload
	aiRequest := map[string]interface{}{
		"videoTitle":       video.Title,
		"videoDescription": video.Description,
		"videoKeywords":    video.Keywords,
		"videoCategory":    video.Category,
	}

	// Encode request payload to JSON
	requestBody, err := json.Marshal(aiRequest)
	if err != nil {
		return nil, fmt.Errorf("error encoding AI request: %w", err)
	}

	// Create HTTP client
	client := &http.Client{}

	// Send POST request to AI service
	req, err := http.NewRequest(http.MethodPost, cfg.AIService, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating AI request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending AI request: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		var apiError struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&apiError)
		return nil, fmt.Errorf("AI service returned error: %s", apiError.Error)
	}

	// Decode AI response
	var aiSuggestions models.AISuggestions
	if err := json.NewDecoder(resp.Body).Decode(&aiSuggestions); err != nil {
		return nil, fmt.Errorf("error decoding AI response: %w", err)
	}

	return &aiSuggestions, nil
}

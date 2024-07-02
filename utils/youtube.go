package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/youtube/v3"

	"github.com/FuseWorkflows/fuse-go-server/models"
)

// UploadVideoToYouTube uploads a video to YouTube
func UploadVideoToYouTube(videoURL, apiKey string, video *models.Video) error {
	// Create a new YouTube service
	ctx := context.Background()
	service, err := youtube.New(ctx) // Replace with your actual client secret path
	// os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
	// or use a service account key (using the client secret file is highly discouraged)
	// os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),

	if err != nil {
		return fmt.Errorf("error creating YouTube service: %w", err)
	}
	service.Client.Transport = &googleapi.DefaultTransport{
		// use your actual API key instead of this
		// Key:      os.Getenv("YOUTUBE_API_KEY"),
		Key: apiKey,
	}

	// Create a new video upload request
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       video.Title,
			Description: video.Description,
			Tags:        video.Keywords,
			CategoryId:  video.Category,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: "public", // Use video.PrivacyStatus to set correct value
		},
	}

	// Get the last iteration
	lastIteration := video.Iterations[len(video.Iterations)-1]

	// Upload the video to YouTube
	// This code assumes you are uploading the video from a URL
	// If you are uploading the video from a local file, you'll need to adjust this code accordingly
	call := service.Videos.Insert("snippet,status", upload)
	call.Media(videoURL, "video/mp4") // Assumes video is in MP4 format

	response, err := call.Do()
	if err != nil {
		return fmt.Errorf("error uploading video to YouTube: %w", err)
	}
	defer response.Body.Close()

	// Handle potential errors
	if response.StatusCode != http.StatusOK {
		// Handle error codes from YouTube API
		var apiError googleapi.Error
		if err := json.NewDecoder(response.Body).Decode(&apiError); err != nil {
			return fmt.Errorf("error decoding YouTube API error: %w", err)
		}
		return fmt.Errorf("YouTube API error: %s", apiError.Message)
	}

	// Extract the video ID from the response
	videoID := response.Data["id"].(string)

	// Update the video in the database with the YouTube video ID
	video.ID = videoID                    // Replace with your own logic to update video ID
	_, err = updateVideo(video.ID, video) // Replace with your own updateVideo function
	if err != nil {
		return fmt.Errorf("error updating video in database: %w", err)
	}

	// Optionally: fetch the video from YouTube to update metadata from the API
	// This is optional but can be useful for updating the video metadata in your database
	// You can fetch the video using the video ID you just got:
	// fetchedVideo, err := service.Videos.List("snippet,statistics,contentDetails").Id(videoID).Do()
	// Handle errors and use the fetched data from fetchedVideo

	return nil
}

// Helper function for updating video in the database (replace with your own logic)
func updateVideo(videoID string, video *models.Video) (int64, error) {
	// This is a placeholder, replace with your own database interaction
	// You can use your existing database library
	// For example, using a database/sql driver:
	ctx := context.Background()
	result, err := db.ExecContext(ctx, "UPDATE videos SET id = $1 WHERE id = $2", video.ID, videoID)
	if err != nil {
		return 0, fmt.Errorf("error updating video: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting rows affected: %w", err)
	}
	return rowsAffected, nil
}

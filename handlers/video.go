package handlers

import (
	"errors"
	"net/http"

	"github.com/FuseWorkflows/fuse-go-server/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/your-username/youtube-backend/config"
	"github.com/your-username/youtube-backend/database"
	"github.com/your-username/youtube-backend/models"
	"github.com/your-username/youtube-backend/utils"
)

// GetVideoHandler retrieves a list of videos
func GetVideoHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserIDFromContext(r)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "User not authenticated"})
			return
		}

		videos, err := db.GetVideosByUser(userID)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to fetch videos"})
			return
		}

		render.JSON(w, r, videos)
	}
}

// CreateVideoHandler creates a new video
func CreateVideoHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserIDFromContext(r)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "User not authenticated"})
			return
		}

		var video models.Video
		if err := render.Bind(r, &video); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid video data"})
			return
		}

		// Ensure the user owns the channel
		channel, err := db.GetChannelByID(video.Channel.ID)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Channel not found"})
			return
		}
		if channel.Owner.ID != userID {
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, map[string]string{"error": "You are not authorized to create videos in this channel"})
			return
		}

		createdVideo, err := db.CreateVideo(&video)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to create video"})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, createdVideo)
	}
}

// GetVideoByIDHandler retrieves a specific video by ID
func GetVideoByIDHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID := chi.URLParam(r, "videoID")
		if videoID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Video ID is required"})
			return
		}

		video, err := db.GetVideoByID(videoID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Video not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to fetch video"})
			return
		}

		render.JSON(w, r, video)
	}
}

// UpdateVideoHandler updates a video by ID
func UpdateVideoHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID := chi.URLParam(r, "videoID")
		if videoID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Video ID is required"})
			return
		}

		var video models.Video
		if err := render.Bind(r, &video); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid video data"})
			return
		}

		updatedVideo, err := db.UpdateVideo(videoID, &video)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Video not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to update video"})
			return
		}

		render.JSON(w, r, updatedVideo)
	}
}

// DeleteVideoHandler deletes a video by ID
func DeleteVideoHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID := chi.URLParam(r, "videoID")
		if videoID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Video ID is required"})
			return
		}

		err := db.DeleteVideo(videoID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Video not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to delete video"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"message": "Video deleted successfully"})
	}
}

// UploadVideoHandler uploads a video to YouTube
func UploadVideoHandler(db *database.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID := chi.URLParam(r, "videoID")
		if videoID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Video ID is required"})
			return
		}

		video, err := db.GetVideoByID(videoID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Video not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to fetch video"})
			return
		}

		// Get the last iteration
		lastIteration := video.Iterations[len(video.Iterations)-1]

		// Upload the video to YouTube
		err = utils.UploadVideoToYouTube(lastIteration.URL, video.Channel.API_KEY, &video)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to upload video to YouTube"})
			return
		}

		// Update video status to "published"
		video.Status = models.Published
		updatedVideo, err := db.UpdateVideo(videoID, &video)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to update video status"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, updatedVideo)
	}
}

// GetAISuggestionsHandler retrieves AI suggestions for video metadata
func GetAISuggestionsHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var video models.Video
		if err := render.Bind(r, &video); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid video data"})
			return
		}

		aiSuggestions, err := utils.GetAISuggestions(&video, cfg)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to get AI suggestions"})
			return
		}

		render.JSON(w, r, aiSuggestions)
	}
}

// CreateVideoHandler creates a new video
func CreateVideoHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserIDFromContext(r)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "User not authenticated"})
			return
		}

		var video models.Video
		if err := render.Bind(r, &video); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid video data"})
			return
		}

		// Ensure the user owns the channel
		channel, err := db.GetChannelByID(video.Channel.ID)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Channel not found"})
			return
		}
		if channel.Owner.ID != userID {
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, map[string]string{"error": "You are not authorized to create videos in this channel"})
			return
		}

		// Get AI suggestions
		aiSuggestions, err := utils.GetAISuggestions(&video, cfg)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to get AI suggestions"})
			return
		}

		// Update video metadata with AI suggestions
		video.Title = aiSuggestions.Title
		video.Description = aiSuggestions.Description
		video.Keywords = aiSuggestions.Keywords
		video.Category = aiSuggestions.Category

		createdVideo, err := db.CreateVideo(&video)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to create video"})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, createdVideo)
	}
}

// UpdateVideoHandler updates a video by ID
func UpdateVideoHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID := chi.URLParam(r, "videoID")
		if videoID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Video ID is required"})
			return
		}

		var video models.Video
		if err := render.Bind(r, &video); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid video data"})
			return
		}

		updatedVideo, err := db.UpdateVideo(videoID, &video)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Video not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to update video"})
			return
		}

		render.JSON(w, r, updatedVideo)
	}
}

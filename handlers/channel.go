package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/FuseWorkflows/fuse-go-server/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/FuseWorkflows/fuse-go-server/database"
	"github.com/FuseWorkflows/fuse-go-server/models"
)

// GetChannelHandler retrieves a list of channels
func GetChannelHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserIDFromContext(r)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "User not authenticated"})
			return
		}

		channels, err := db.GetChannelsByUser(userID)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to fetch channels"})
			return
		}

		render.JSON(w, r, channels)
	}
}

// CreateChannelHandler creates a new channel
func CreateChannelHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserIDFromContext(r)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "User not authenticated"})
			return
		}

		var channel models.Channel
		// if err := render.Bind(r, &channel); err != nil {
		// 	render.Status(r, http.StatusBadRequest)
		// 	render.JSON(w, r, map[string]string{"error": "Invalid channel data"})
		// 	return
		// }

		if err := json.NewDecoder(r.Body).Decode(&channel); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid channel data"})
			return
		}

		channel.Owner.ID = userID

		createdChannel, err := db.CreateChannel(&channel)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to create channel"})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, createdChannel)
	}
}

// GetChannelByIDHandler retrieves a specific channel by ID
func GetChannelByIDHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channelID := chi.URLParam(r, "channelID")
		if channelID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Channel ID is required"})
			return
		}

		channel, err := db.GetChannelByID(channelID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Channel not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to fetch channel"})
			return
		}

		render.JSON(w, r, channel)
	}
}

// UpdateChannelHandler updates a channel by ID
func UpdateChannelHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channelID := chi.URLParam(r, "channelID")
		if channelID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Channel ID is required"})
			return
		}

		var channel models.Channel
		if err := json.NewDecoder(r.Body).Decode(&channel); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid channel data"})
			return
		}

		updatedChannel, err := db.UpdateChannel(channelID, &channel)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Channel not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to update channel"})
			return
		}

		render.JSON(w, r, updatedChannel)
	}
}

// DeleteChannelHandler deletes a channel by ID
func DeleteChannelHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channelID := chi.URLParam(r, "channelID")
		if channelID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Channel ID is required"})
			return
		}

		err := db.DeleteChannel(channelID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Channel not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to delete channel"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"message": "Channel deleted successfully"})
	}
}

package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/your-username/youtube-backend/database"
	"github.com/your-username/youtube-backend/models"
)

// GetEditorHandler retrieves a list of editors
func GetEditorHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		editors, err := db.GetEditors()
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to fetch editors"})
			return
		}

		render.JSON(w, r, editors)
	}
}

// CreateEditorHandler creates a new editor
func CreateEditorHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var editor models.Editor
		if err := render.Bind(r, &editor); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid editor data"})
			return
		}

		createdEditor, err := db.CreateEditor(&editor)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to create editor"})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, createdEditor)
	}
}

// GetEditorByIDHandler retrieves a specific editor by ID
func GetEditorByIDHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		editorID := chi.URLParam(r, "editorID")
		if editorID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Editor ID is required"})
			return
		}

		editor, err := db.GetEditorByID(editorID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Editor not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to fetch editor"})
			return
		}

		render.JSON(w, r, editor)
	}
}

// UpdateEditorHandler updates an editor by ID
func UpdateEditorHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		editorID := chi.URLParam(r, "editorID")
		if editorID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Editor ID is required"})
			return
		}

		var editor models.Editor
		if err := render.Bind(r, &editor); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid editor data"})
			return
		}

		updatedEditor, err := db.UpdateEditor(editorID, &editor)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Editor not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to update editor"})
			return
		}

		render.JSON(w, r, updatedEditor)
	}
}

// DeleteEditorHandler deletes an editor by ID
func DeleteEditorHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		editorID := chi.URLParam(r, "editorID")
		if editorID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Editor ID is required"})
			return
		}

		err := db.DeleteEditor(editorID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Editor not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to delete editor"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"message": "Editor deleted successfully"})
	}
}

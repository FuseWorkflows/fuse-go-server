package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/FuseWorkflows/fuse-go-server/database"
	"github.com/FuseWorkflows/fuse-go-server/models"
)

// GetIterationHandler retrieves a list of iterations
func GetIterationHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iterations, err := db.GetIterations()
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to fetch iterations"})
			return
		}

		render.JSON(w, r, iterations)
	}
}

// CreateIterationHandler creates a new iteration
func CreateIterationHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var iteration models.Iteration
		if err := render.Bind(r, &iteration); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid iteration data"})
			return
		}

		createdIteration, err := db.CreateIteration(&iteration)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to create iteration"})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, createdIteration)
	}
}

// GetIterationByIDHandler retrieves a specific iteration by ID
func GetIterationByIDHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iterationID := chi.URLParam(r, "iterationID")
		if iterationID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Iteration ID is required"})
			return
		}

		iteration, err := db.GetIterationByID(iterationID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Iteration not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to fetch iteration"})
			return
		}

		render.JSON(w, r, iteration)
	}
}

// UpdateIterationHandler updates an iteration by ID
func UpdateIterationHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iterationID := chi.URLParam(r, "iterationID")
		if iterationID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Iteration ID is required"})
			return
		}

		var iteration models.Iteration
		if err := render.Bind(r, &iteration); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid iteration data"})
			return
		}

		updatedIteration, err := db.UpdateIteration(iterationID, &iteration)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Iteration not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to update iteration"})
			return
		}

		render.JSON(w, r, updatedIteration)
	}
}

// DeleteIterationHandler deletes an iteration by ID
func DeleteIterationHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iterationID := chi.URLParam(r, "iterationID")
		if iterationID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Iteration ID is required"})
			return
		}

		err := db.DeleteIteration(iterationID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Iteration not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to delete iteration"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"message": "Iteration deleted successfully"})
	}
}

// AddNoteToIterationHandler adds a note to an iteration
func AddNoteToIterationHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iterationID := chi.URLParam(r, "iterationID")
		if iterationID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Iteration ID is required"})
			return
		}

		var note models.Note
		if err := render.Bind(r, &note); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid note data"})
			return
		}

		err := db.AddNoteToIteration(iterationID, &note)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "Iteration not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to add note to iteration"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"message": "Note added successfully"})
	}
}

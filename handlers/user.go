package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/FuseWorkflows/fuse-go-server/database"
	"github.com/FuseWorkflows/fuse-go-server/models"
)

// GetUserHandler retrieves a list of users
func GetUserHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := db.GetUsers()
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to fetch users"})
			return
		}

		render.JSON(w, r, users)
	}
}

// CreateUserHandler creates a new user
func CreateUserHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := render.Bind(r, &user); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid user data"})
			return
		}

		createdUser, err := db.CreateUser(&user)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to create user"})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, createdUser)
	}
}

// UpdateUserHandler updates a user by ID
func UpdateUserHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "User ID is required"})
			return
		}

		var user models.User
		if err := render.Bind(r, &user); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid user data"})
			return
		}

		updatedUser, err := db.UpdateUser(userID, &user)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "User not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to update user"})
			return
		}

		render.JSON(w, r, updatedUser)
	}
}

// DeleteUserHandler deletes a user by ID
func DeleteUserHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "User ID is required"})
			return
		}

		err := db.DeleteUser(userID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "User not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to delete user"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"message": "User deleted successfully"})
	}
}

package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"

	"github.com/FuseWorkflows/fuse-go-server/config"
	"github.com/FuseWorkflows/fuse-go-server/database"
	"github.com/FuseWorkflows/fuse-go-server/models"
)

// SignupHandler handles user signup
func SignupHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid user data"})
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)

		// Create the user
		createdUser, err := db.CreateUser(&user)
		if err != nil {
			fmt.Println(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to create user"})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, createdUser)
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginRequest) Bind(r *http.Request) error {
	if l.Email == "" || l.Password == "" {
		return errors.New("missing required fields")
	}
	return nil
}

// LoginHandler handles user login
func LoginHandler(db *database.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest LoginRequest

		if err := render.Bind(r, &loginRequest); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid login data"})
			return
		}

		// Find the user by email
		user, err := db.GetUserByEmail(loginRequest.Email)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "User not found"})
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to find user"})
			return
		}

		// Compare the password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Incorrect password"})
			return
		}

		// Create JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID":    user.ID,
			"issuedAt":  time.Now().Unix(),
			"expiresAt": time.Now().Add(time.Hour * 24).Unix(),
		})

		// Sign the token
		tokenString, err := token.SignedString([]byte(cfg.JWTKey))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to sign token"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]interface{}{
			"token": tokenString,
		})
	}
}

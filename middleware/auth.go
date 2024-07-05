package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"

	"github.com/FuseWorkflows/fuse-go-server/database"
	"github.com/FuseWorkflows/fuse-go-server/models"
)

// Auth middleware for JWT authentication
func Auth(db *database.DB, jwtKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Path == "/auth/signup" {
				// Skip authentication for signup
				next.ServeHTTP(w, r)
				return
			}

			if r.URL.Path == "/users/" {
				// Skip authentication for signup
				next.ServeHTTP(w, r)
				return
			}

			authorizationHeader := r.Header.Get("Authorization")
			if authorizationHeader == "" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "Authorization header is required"})
				return
			}

			tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtKey), nil
			})

			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "Invalid token"})
				return
			}

			if !token.Valid {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "Token is invalid"})
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "Invalid token claims"})
				return
			}

			userID, ok := claims["userID"].(string)
			if !ok {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "User ID is missing in token claims"})
				return
			}

			// Fetch user from database
			user, err := db.GetUserByID(userID)
			if err != nil {
				if errors.Is(err, database.ErrNotFound) {
					render.Status(r, http.StatusUnauthorized)
					render.JSON(w, r, map[string]string{"error": "User not found"})
					return
				}
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, map[string]string{"error": "Failed to fetch user"})
				return
			}

			// Attach user to context
			ctx := context.WithValue(r.Context(), "user", user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(r *http.Request) (string, error) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		return "", errors.New("user not found in context")
	}
	return user.ID, nil
}

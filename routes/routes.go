package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/FuseWorkflows/fuse-go-server/config"
	"github.com/FuseWorkflows/fuse-go-server/database"
	"github.com/FuseWorkflows/fuse-go-server/handlers"
)

func InitRoutes(r *chi.Mux, db *database.DB, cfg *config.Config) {

	// Authentication routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", handlers.SignupHandler(db))
		r.Post("/login", handlers.LoginHandler(db, cfg))
	})

	// User routes
	r.Route("/users", func(r chi.Router) {
		r.Get("/", handlers.GetUserHandler(db))
		r.Post("/", handlers.CreateUserHandler(db))
		r.Patch("/{userID}", handlers.UpdateUserHandler(db))
		r.Delete("/{userID}", handlers.DeleteUserHandler(db))
	})

	// Channel routes
	r.Route("/channels", func(r chi.Router) {
		r.Get("/", handlers.GetChannelHandler(db))
		r.Post("/", handlers.CreateChannelHandler(db))
		r.Get("/{channelID}", handlers.GetChannelByIDHandler(db))
		r.Patch("/{channelID}", handlers.UpdateChannelHandler(db))
		r.Delete("/{channelID}", handlers.DeleteChannelHandler(db))
	})

	// Video routes
	r.Route("/videos", func(r chi.Router) {
		r.Get("/", handlers.GetVideoHandler(db))
		r.Post("/", handlers.CreateVideoHandler(db))
		r.Get("/{videoID}", handlers.GetVideoByIDHandler(db))
		r.Patch("/{videoID}", handlers.UpdateVideoHandler(db))
		r.Delete("/{videoID}", handlers.DeleteVideoHandler(db))
		r.Post("/{videoID}/upload", handlers.UploadVideoHandler(db, cfg))
	})

	// Iteration routes
	r.Route("/iterations", func(r chi.Router) {
		r.Get("/", handlers.GetIterationHandler(db))
		r.Post("/", handlers.CreateIterationHandler(db))
		r.Get("/{iterationID}", handlers.GetIterationByIDHandler(db))
		r.Patch("/{iterationID}", handlers.UpdateIterationHandler(db))
		r.Delete("/{iterationID}", handlers.DeleteIterationHandler(db))
		r.Post("/{iterationID}/notes", handlers.AddNoteToIterationHandler(db))
	})

	// Editor routes
	r.Route("/editors", func(r chi.Router) {
		r.Get("/", handlers.GetEditorHandler(db))
		r.Post("/", handlers.CreateEditorHandler(db))
		r.Get("/{editorID}", handlers.GetEditorByIDHandler(db))
		r.Patch("/{editorID}", handlers.UpdateEditorHandler(db))
		r.Delete("/{editorID}", handlers.DeleteEditorHandler(db))
	})

	// AI routes
	r.Route("/ai", func(r chi.Router) {
		r.Post("/suggestions", handlers.GetAISuggestionsHandler(cfg))
	})
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/FuseWorkflows/fuse-go-server/config"
	"github.com/FuseWorkflows/fuse-go-server/database"
	"github.com/FuseWorkflows/fuse-go-server/middleware"
	"github.com/FuseWorkflows/fuse-go-server/routes"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Initialize config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Error initializing config:", err)
	}

	// Initialize router
	r := chi.NewRouter()

	// Basic middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS middleware
	corsCfg := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Replace with allowed origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Replace with allowed headers
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsCfg.Handler)

	// Authentication middleware
	r.Use(middleware.Auth(db, cfg.JWTKey))

	// Routes
	routes.InitRoutes(r, db, cfg)

	// Start server
	fmt.Printf("Server listening on port %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r))
}

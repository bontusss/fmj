package main

import (
	"context"
	"fmj/config"
	"fmj/internal/auth"
	"fmj/internal/email"
	"fmj/middleware"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	gowebly "github.com/gowebly/helpers"
)

// runServer runs a new HTTP server with the loaded environment variables.
func runServer(db *mongo.Database, cfg *config.Config) error {
	// Validate environment variables.
	port, err := strconv.Atoi(gowebly.Getenv("BACKEND_PORT", "7000"))
	if err != nil {
		return err
	}

	// Initialize services
	emailService := email.NewService(cfg)
	authRepo := auth.NewRepository(db, context.Context(context.Background()))
	authService := auth.NewService(authRepo, emailService)
	authHandler := auth.NewHandler(authService, cfg)

	// Create a new gin server.
	router := gin.Default()

	// Handle static files.
	router.Static("/static", "./static")

	// Setup sessions
	store := cookie.NewStore([]byte(cfg.SessionSecret))
	store.Options(sessions.Options{
		MaxAge: 60 * 60 * 24 * 10, // 10 days
		Path:   "/",
	})
	router.Use(sessions.Sessions("auth_session", store))
	// Apply CheckAuth to public routes
	router.Use(middleware.CheckAuth())

	// Register auth routes
	authHandler.RegisterRoutes(router)

	// Handle index page view.
	router.GET("/", indexViewHandler)

	// Handle API endpoints.
	router.GET("/api/hello-world", showContentAPIHandler)

	// protected ungrouped routes
	protected := router.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/dashboard", showDashboardHandler)
	}
	// Create a new server instance with options from environment variables.
	// For more information, see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	// Send log message.
	slog.Info("Starting server...", "port", port)

	return server.ListenAndServe()
}

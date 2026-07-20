package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/example/yt-downloader/internal/db"
	"github.com/example/yt-downloader/internal/handlers"
	"github.com/example/yt-downloader/internal/middleware"
	"github.com/example/yt-downloader/internal/telegram"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// Database connection setup
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "yt_downloader"
	}

	database, err := db.Connect(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Execute migrations
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	r := gin.Default()
	// Trust only local reverse proxies by default; avoids "trusted all proxies" warning.
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	// CORS Setup — origins can be overridden via CORS_ORIGINS env var (comma-separated)
	corsOrigins := []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175", "http://localhost:3000"}
	if envOrigins := os.Getenv("CORS_ORIGINS"); envOrigins != "" {
		corsOrigins = strings.Split(envOrigins, ",")
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize Rate Limiters
	// 1. Core API (Playlist Fetching & Video Streaming) - 5 requests per minute (limit=5/60, burst=5)
	apiLimiter := middleware.NewIPRateLimiter(rate.Limit(5.0/60.0), 5)
	// 2. Auth API (Login/Register attempts) - 10 requests per minute (limit=10/60, burst=10)
	authLimiter := middleware.NewIPRateLimiter(rate.Limit(10.0/60.0), 10)

	authHandler := handlers.NewAuthHandler(database)

	api := r.Group("/api")
	{
		// Apply IP-based rate limiting to heavy download/info endpoints
		coreAPI := api.Group("")
		coreAPI.Use(middleware.RateLimitMiddleware(apiLimiter))
		{
			coreAPI.POST("/playlist/info", handlers.GetPlaylist)
			coreAPI.GET("/stream", handlers.StreamDownload) // Direct streaming endpoint
		}

		// Auth & Profile Routes
		auth := api.Group("/auth")
		auth.Use(middleware.RateLimitMiddleware(authLimiter))
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)

			// Authenticated Routes
			authenticated := auth.Group("")
			authenticated.Use(middleware.AuthMiddleware())
			{
				authenticated.GET("/profile", authHandler.GetProfile)
				authenticated.PUT("/profile", authHandler.UpdateProfile)
				authenticated.POST("/download", authHandler.LogDownload)
			}
		}
	}

	// Telegram Bot Setup
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken != "" {
		bot, err := telegram.NewBot(botToken, database)
		if err != nil {
			log.Printf("Failed to initialize Telegram bot: %v", err)
		} else {
			go func() {
				log.Println("Telegram bot starting...")
				bot.Start()
			}()
		}
	} else {
		log.Println("TELEGRAM_BOT_TOKEN not set, bot will not start.")
	}

	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

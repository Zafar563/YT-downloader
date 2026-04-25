package main

import (
	"log"
	"os"
	"time"

	"github.com/example/yt-downloader/internal/handlers"
	"github.com/example/yt-downloader/internal/telegram"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS Setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175", "http://localhost:3000"}, // Allow Vite frontend ports
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		api.POST("/playlist/info", handlers.GetPlaylist)
		api.GET("/stream", handlers.StreamDownload) // Direct streaming endpoint
	}

	// Telegram Bot Setup
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken != "" {
		bot, err := telegram.NewBot(botToken)
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

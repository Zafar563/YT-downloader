package main

import (
	"log"
	"time"

	"github.com/example/yt-downloader/internal/handlers"
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
		api.POST("/download", handlers.StartDownload)
		api.GET("/stream", handlers.StreamDownload) // New streaming endpoint
	}

	r.GET("/ws", handlers.WSHandler)
	r.Static("/downloads", "./downloads")
	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

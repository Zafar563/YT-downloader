package handlers

import (
	"net/http"
	"path/filepath"
	"sync"
    "fmt"

	"github.com/example/yt-downloader/internal/downloader"
	"github.com/example/yt-downloader/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all CORS for dev
	},
}

// Global variable to keep track of active connections
// In a real app, this should be better managed
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.DownloadProgress)
var clientsMu sync.Mutex

func setupWebsocketHub() {
    go func() {
        for msg := range broadcast {
            clientsMu.Lock()
            for client := range clients {
                err := client.WriteJSON(msg)
                if err != nil {
                    client.Close()
                    delete(clients, client)
                }
            }
            clientsMu.Unlock()
        }
    }()
}

func init() {
    setupWebsocketHub()
}

// WSHandler handles WebSocket connections
func WSHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WS Upgrade Error:", err)
		return
	}
	clientsMu.Lock()
	clients[ws] = true
	clientsMu.Unlock()
}

// GetPlaylist handles request to fetch playlist metadata
func GetPlaylist(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	playlist, err := downloader.GetPlaylistInfo(req.URL)
	if err != nil {
        fmt.Println("Error fetching playlist:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, playlist)
}

// StartDownload handles request to start downloading videos
func StartDownload(c *gin.Context) {
	var req models.DownloadRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Simple semaphore for concurrency control (limit 3 concurrent downloads)
	sem := make(chan struct{}, 3)
	// Output directory
	outputDir, _ := filepath.Abs("downloads")

	// Launch a goroutine to handle the entire batch download so the HTTP request returns immediately
	go func(urls []string) {
		var wg sync.WaitGroup
		for _, url := range urls {
            // Wait for a slot
			sem <- struct{}{}
			wg.Add(1)
            
            // Extract Video ID from URL (naive approach, better pass ID from frontend)
            // assuming url is https://www.youtube.com/watch?v=VIDEO_ID
            // For now, let's use the URL as ID for progress tracking if ID not available, 
            // but the frontend should ideally pass IDs. Let's assume the frontend passes full URLs.
            // A better way is to pass Video objects. 
            // Let's generate a temporary ID or use the URL as key.
            videoID := url 

			go func(u string, vID string) {
				defer wg.Done()
				defer func() { <-sem }() // Release slot

				progressChan := make(chan models.DownloadProgress)
                
                // Forward progress from downloader to websocket
                go func() {
                    for p := range progressChan {
                        broadcast <- p
                    }
                }()

				downloader.DownloadVideo(vID, u, req.Format, outputDir, progressChan)
			}(url, videoID)
		}
		wg.Wait()
        // Notify all finished? Maybe not needed as individual finished events are sent.
	}(req.URLs)

	c.JSON(http.StatusOK, gin.H{"message": "Download started", "count": len(req.URLs)})
}

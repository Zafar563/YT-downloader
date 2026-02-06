package handlers

import (
	"net/http"
	"path/filepath"
	"sync"
	"fmt"
    "regexp"


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
            
            // Extract Video ID from URL
            var videoID string
            regex := regexp.MustCompile(`(?:v=|/)([0-9A-Za-z_-]{11}).*`)
            matches := regex.FindStringSubmatch(url)
            if len(matches) > 1 {
                videoID = matches[1]
            } else {
                videoID = url // Fallback to URL if extraction fails
            } 

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

// StreamDownload handles direct streaming request
func StreamDownload(c *gin.Context) {
    url := c.Query("url")
    format := c.Query("format")
    title := c.Query("title")

    if url == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
        return
    }

    // Set headers for download
    filename := fmt.Sprintf("%s.mp4", title)
    contentType := "video/mp4"
    if format == "mp3" {
        filename = fmt.Sprintf("%s.mp3", title)
        contentType = "audio/mpeg"
    }

    // Sanitize filename (basic)
    filename = filepath.Base(filename) 

    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
    c.Header("Content-Type", contentType)

    err := downloader.StreamVideo(url, format, c.Writer)
    if err != nil {
        // Can't write JSON error if headers already sent, but we can log
        fmt.Println("Streaming error:", err)
    }
}

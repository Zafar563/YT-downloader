package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/example/yt-downloader/internal/downloader"
	"github.com/gin-gonic/gin"
)


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


// StreamDownload handles direct streaming request
func StreamDownload(c *gin.Context) {
	url := c.Query("url")
	format := c.Query("format")
	quality := c.Query("quality")
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

	err := downloader.StreamVideo(url, format, quality, c.Writer)
	if err != nil {
		// Can't write JSON error if headers already sent, but we can log
		fmt.Println("Streaming error:", err)
	}
}

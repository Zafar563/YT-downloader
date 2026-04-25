package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

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

	if format == "mp3" {
		// Download to a temporary file on the server, convert, and then send it
		tempFileName := fmt.Sprintf("yt_%d.mp3", time.Now().UnixNano())
		tempPath := filepath.Join(os.TempDir(), tempFileName)

		// Remove the temp file after we are done
		defer os.Remove(tempPath)

		err := downloader.DownloadToPath(url, format, quality, tempPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download and convert: " + err.Error()})
			return
		}

		filename := fmt.Sprintf("%s.mp3", title)
		filename = filepath.Base(filename)

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		c.Header("Content-Type", "audio/mpeg")
		
		c.File(tempPath)
		return
	}

	// For video, we can stream directly
	filename := fmt.Sprintf("%s.mp4", title)
	filename = filepath.Base(filename)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "video/mp4")

	err := downloader.StreamVideo(url, format, quality, c.Writer)
	if err != nil {
		fmt.Println("Streaming error:", err)
	}
}

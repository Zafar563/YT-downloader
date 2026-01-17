package downloader

import (
	"bufio"
	"encoding/json"
	"fmt"
    "os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/example/yt-downloader/internal/models"
)

// GetPlaylistInfo fetches metadata for a playlist or video
func GetPlaylistInfo(url string) (*models.Playlist, error) {
	// --flat-playlist is much faster for large playlists as it doesn't extract full info for every video immediately
    // --dump-single-json outputs the result as a single JSON object
	// Use local yt-dlp.exe if available, or assume it's in PATH
    exePath := "./yt-dlp.exe"
    if _, err := os.Stat(exePath); os.IsNotExist(err) {
        exePath = "yt-dlp"
    }
    
	cmd := exec.Command(exePath, "--dump-single-json", "--flat-playlist", "--no-warnings", url)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute yt-dlp: %w", err)
	}

	var result models.Playlist
    // Sometimes yt-dlp returns a single video as a playlist structure, sometimes just the video.
    // For simplicity, we decode into a generic map first or try to decode into Playlist.
    // However, since we defined Playlist struct matching expected output, let's try direct unmarshal.
    
    // Check if it's a playlist or single video.
    var raw map[string]interface{}
    if err := json.Unmarshal(output, &raw); err != nil {
        return nil, err
    }

    if _, ok := raw["entries"]; ok {
        // It's a playlist
        if err := json.Unmarshal(output, &result); err != nil {
            return nil, err
        }
    } else {
        // It's a single video, wrap it in a playlist
        var video models.Video
        if err := json.Unmarshal(output, &video); err != nil {
            return nil, err
        }
        result.Title = "Single Video"
        result.Entries = []models.Video{video}
    }

	return &result, nil
}

// DownloadVideo downloads a video and sends progress updates
func DownloadVideo(videoID string, url string, outputDir string, progressChan chan<- models.DownloadProgress) {
    defer close(progressChan)

    // Output template to save in downloads folder
    outputPath := fmt.Sprintf("%s/%%(title)s.%%(ext)s", outputDir)

    // Command to download
    // --newline forces progress to be printed on new lines for easier parsing
    // --progress-template prints custom progress format
    exePath := "./yt-dlp.exe"
    if _, err := os.Stat(exePath); os.IsNotExist(err) {
        exePath = "yt-dlp"
    }

    cmd := exec.Command(exePath, "-x", "--audio-format", "mp3", "-o", outputPath, "--newline", "--no-warnings", url)

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        progressChan <- models.DownloadProgress{VideoID: videoID, Status: "error", Message: err.Error()}
        return
    }

    // Capture stderr
    stderr, _ := cmd.StderrPipe()
    
    if err := cmd.Start(); err != nil {
        fmt.Println("Download start error:", err)
        progressChan <- models.DownloadProgress{VideoID: videoID, Status: "error", Message: err.Error()}
        return
    }

    // Parse stdout for progress
    scanner := bufio.NewScanner(stdout)
    
    // Read stderr in a separate goroutine
    var stderrOut strings.Builder
    go func() {
        scannerErr := bufio.NewScanner(stderr)
        for scannerErr.Scan() {
            stderrOut.WriteString(scannerErr.Text() + "\n")
        }
    }()
    
    // Regex to match standard yt-dlp progress output [download]  45.0% of 10.00MiB at  2.50MiB/s ETA 00:05
    progressRegex := regexp.MustCompile(`\[download\]\s+(\d+\.?\d*)%`)

    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "[download]") {
            matches := progressRegex.FindStringSubmatch(line)
            if len(matches) > 1 {
                percent, _ := strconv.ParseFloat(matches[1], 64)
                progressChan <- models.DownloadProgress{
                    VideoID: videoID,
                    Status:  "downloading",
                    Percent: percent,
                }
            }
            if strings.Contains(line, "100%") || strings.Contains(line, "100.0%") {
                 progressChan <- models.DownloadProgress{VideoID: videoID, Status: "finished", Percent: 100}
            }
        }
    }

    if err := cmd.Wait(); err != nil {
         fmt.Println("Download wait error:", err)
         fmt.Println("Stderr Output:", stderrOut.String()) // Print captured stderr
         progressChan <- models.DownloadProgress{VideoID: videoID, Status: "error", Message: "Download failed: " + stderrOut.String()}
    } else {
         progressChan <- models.DownloadProgress{VideoID: videoID, Status: "finished", Percent: 100}
    }
}

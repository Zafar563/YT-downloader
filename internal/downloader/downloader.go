package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/example/yt-downloader/internal/models"
)

func getYTdlpPath() string {
    // Priority: environment variable, then path, then local file (mostly for dev)
    if path := os.Getenv("YT_DLP_PATH"); path != "" {
        return path
    }
    if path, err := exec.LookPath("yt-dlp"); err == nil {
        return path
    }
    // Fallback to local file if it exists
    if _, err := os.Stat("./yt-dlp"); err == nil {
        return "./yt-dlp"
    }
    return "yt-dlp" // Assume it's in PATH if all else fails
}

// GetPlaylistInfo fetches metadata for a playlist or video
func GetPlaylistInfo(url string) (*models.Playlist, error) {
    exePath := getYTdlpPath()
    
	cmd := exec.Command(exePath, 
        "--dump-single-json", 
        "--flat-playlist", 
        "--no-warnings", 
        "--extractor-args", "youtube:player_client=android,web",
        url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute yt-dlp (%s): %w, output: %s", exePath, err, string(output))
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


// StreamVideo streams the video directly to the writer
func StreamVideo(url string, format string, quality string, writer io.Writer) error {
    exePath := getYTdlpPath()

    var cmd *exec.Cmd
    if format == "mp3" {
        // -o - directs output to stdout
        cmd = exec.Command(exePath, "-x", "--audio-format", "mp3", "-o", "-", "--no-warnings", "--extractor-args", "youtube:player_client=android,web", url)
    } else {
        formatArgs := "bestvideo+bestaudio/best"
        if quality != "" && quality != "best" {
            formatArgs = fmt.Sprintf("bestvideo[height<=?%s]+bestaudio/best", strings.TrimSuffix(quality, "p"))
        }
        cmd = exec.Command(exePath, "-f", formatArgs, "-o", "-", "--no-warnings", "--extractor-args", "youtube:player_client=android,web", url)
    }

    cmd.Stdout = writer
    cmd.Stderr = os.Stderr // Pipe stderr to server logs for debugging

    return cmd.Run()
}

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

func getCookiesArgs() []string {
    cookiesPath := "/app/cookies.txt"
    if _, err := os.Stat(cookiesPath); err == nil {
        return []string{"--cookies", cookiesPath}
    }
    // Fallback for local development
    if _, err := os.Stat("./cookies.txt"); err == nil {
        return []string{"--cookies", "./cookies.txt"}
    }
    return nil
}

// GetPlaylistInfo fetches metadata for a playlist or video
func GetPlaylistInfo(url string) (*models.Playlist, error) {
    exePath := getYTdlpPath()
    
	args := []string{
		"--dump-single-json",
		"--flat-playlist",
		"--no-warnings",
		"--extractor-args", "youtube:player_client=android,web",
	}
	args = append(args, getCookiesArgs()...)
	args = append(args, url)

	cmd := exec.Command(exePath, args...)
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

	var args []string
	if format == "mp3" {
		// yt-dlp ignores -x and --audio-format when streaming to stdout (-o -), 
		// so we just request the best audio stream
		args = []string{"-f", "bestaudio[ext=m4a]/bestaudio/best", "-o", "-", "--no-warnings"}
	} else {
		formatArgs := "bestvideo+bestaudio/best"
		if quality != "" && quality != "best" {
			height := strings.TrimSuffix(quality, "p")
			formatArgs = fmt.Sprintf("bestvideo[height<=?%s]+bestaudio/best[height<=?%s]/best", height, height)
		}
		args = []string{"-f", formatArgs, "-o", "-", "--no-warnings"}
	}
	args = append(args, "--extractor-args", "youtube:player_client=android,web")
	args = append(args, getCookiesArgs()...)
	args = append(args, url)

	cmd := exec.Command(exePath, args...)

    cmd.Stdout = writer
    cmd.Stderr = os.Stderr // Pipe stderr to server logs for debugging

    return cmd.Run()
}

// DownloadToPath downloads the video/audio to a specific file path
func DownloadToPath(url string, format string, quality string, outputPath string) error {
    exePath := getYTdlpPath()

	var args []string
	if format == "mp3" {
		// Download best audio first to save bandwidth, then convert to mp3
		args = []string{"-f", "bestaudio/best", "-x", "--audio-format", "mp3", "-o", outputPath, "--no-warnings"}
	} else {
		formatArgs := "bestvideo+bestaudio/best"
		if quality != "" && quality != "best" {
			height := strings.TrimSuffix(quality, "p")
			formatArgs = fmt.Sprintf("bestvideo[height<=?%s]+bestaudio/best[height<=?%s]/best", height, height)
		}
		args = []string{"-f", formatArgs, "-o", outputPath, "--no-warnings"}
	}
	args = append(args, "--extractor-args", "youtube:player_client=android,web")
	args = append(args, getCookiesArgs()...)
	args = append(args, url)

	cmd := exec.Command(exePath, args...)

    cmd.Stderr = os.Stderr
    return cmd.Run()
}


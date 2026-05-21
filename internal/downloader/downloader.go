package downloader

import (
	"bytes"
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

// getFormatArgs returns the appropriate format string based on format and quality
func getFormatArgs(format string, quality string) string {
	if format == "mp3" {
		return "bestaudio[ext=m4a]/bestaudio/best"
	}
	formatArgs := "bestvideo+bestaudio/bestvideo/best"
	if quality != "" && quality != "best" {
		height := strings.TrimSuffix(quality, "p")
		formatArgs = fmt.Sprintf("bestvideo[height<=?%s]+bestaudio/bestvideo[height<=?%s]/best[height<=?%s]/best", height, height, height)
	}
	return formatArgs
}

// getCommonArgs returns arguments shared by all yt-dlp invocations.
// Forces the web player client to avoid the degraded TV client that only
// serves muxed streams with a limited format list.
func getCommonArgs() []string {
	return []string{
		"--no-warnings",
		"--extractor-args", "youtube:player_client=web,default",
		"--no-check-certificates",
	}
}

// GetPlaylistInfo fetches metadata for a playlist or video
func GetPlaylistInfo(url string) (*models.Playlist, error) {
	exePath := getYTdlpPath()

	args := []string{
		"--dump-single-json",
		"--flat-playlist",
		"--no-warnings",
	}
	args = append(args, getCookiesArgs()...)
	args = append(args, url)

	cmd := exec.Command(exePath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute yt-dlp: %w\nCommand: %s %s\nOutput: %s", err, exePath, strings.Join(args, " "), string(output))
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
	formatArgs := getFormatArgs(format, quality)

	if format == "mp3" {
		// yt-dlp ignores -x and --audio-format when streaming to stdout (-o -),
		// so we just request the best audio stream
		args = []string{"-f", formatArgs, "-o", "-"}
	} else {
		args = []string{"-f", formatArgs, "-o", "-"}
	}
	args = append(args, getCommonArgs()...)
	args = append(args, getCookiesArgs()...)
	args = append(args, url)

	cmd := exec.Command(exePath, args...)

	cmd.Stdout = writer
	var stderr bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr) // Keep logs and return message upstream

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("stream failed: %w\nOutput: %s", err, stderr.String())
	}
	return nil
}

// DownloadToPath downloads the video/audio to a specific file path
func DownloadToPath(url string, format string, quality string, outputPath string) error {
	exePath := getYTdlpPath()

	var args []string
	formatArgs := getFormatArgs(format, quality)

	if format == "mp3" {
		// Download best audio first to save bandwidth, then convert to mp3
		args = []string{"-f", "bestaudio/best", "-x", "--audio-format", "mp3", "-o", outputPath}
	} else {
		args = []string{"-f", formatArgs, "-o", outputPath}
	}
	args = append(args, getCommonArgs()...)
	args = append(args, getCookiesArgs()...)
	args = append(args, url)

	cmd := exec.Command(exePath, args...)

	var stderr bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("download failed: %w\nOutput: %s", err, stderr.String())
	}
	return nil
}

// ValidateStream checks whether yt-dlp can resolve requested formats before headers are sent.
func ValidateStream(url string, format string, quality string) error {
	exePath := getYTdlpPath()

	formatArgs := getFormatArgs(format, quality)
	args := []string{"-f", formatArgs, "--skip-download"}
	args = append(args, getCommonArgs()...)
	args = append(args, getCookiesArgs()...)
	args = append(args, url)

	cmd := exec.Command(exePath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("preflight failed: %w\nCommand: %s %s\nOutput: %s", err, exePath, strings.Join(args, " "), string(output))
	}
	return nil
}

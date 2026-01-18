package models

// Video represents a single video's metadata
type Video struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Duration  float64    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	WebpageURL string `json:"webpage_url"`
    URL       string `json:"url"` // Added to capture 'url' from flat-playlist
    Formats   []Format `json:"formats,omitempty"`
}

// Format represents a video format (e.g., 1080p, 720p)
type Format struct {
    FormatID string `json:"format_id"`
    Note     string `json:"format_note"`
    Ext      string `json:"ext"`
}

// Playlist represents a detailed playlist
type Playlist struct {
	Title string  `json:"title"`
	Entries []Video `json:"entries"`
}

// DownloadRequest represents the request payload for downloading videos
type DownloadRequest struct {
	URLs   []string `json:"urls"`
	Format string   `json:"format"` // "mp3" or "video"
}

// DownloadProgress represents the progress update sent via WebSocket
type DownloadProgress struct {
    VideoID string  `json:"video_id"`
    Status  string  `json:"status"` // "downloading", "finished", "error"
    Percent float64 `json:"percent"`
    Message string  `json:"message,omitempty"`
}

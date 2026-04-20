package models

// Video represents a single video's metadata
type Video struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Duration   float64  `json:"duration"`
	Thumbnail  string   `json:"thumbnail"`
	WebpageURL string   `json:"webpage_url"`
	URL        string   `json:"url"` // Added to capture 'url' from flat-playlist
	Formats    []Format `json:"formats,omitempty"`
}

// Format represents a video format (e.g., 1080p, 720p)
type Format struct {
	FormatID string `json:"format_id"`
	Note     string `json:"format_note"`
	Ext      string `json:"ext"`
}

// Playlist represents a detailed playlist
type Playlist struct {
	Title   string  `json:"title"`
	Entries []Video `json:"entries"`
}


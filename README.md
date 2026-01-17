# YouTube Downloader

A robust full-stack application for downloading YouTube videos and playlists with real-time progress tracking.

## Features

- **Playlist Support**: Fetch and download entire playlists or single videos.
- **Batch Processing**: Concurrent downloads with a managed worker pool.
- **Real-time Progress**: WebSocket integration for live download progress updates on the frontend.
- **Automatic Cleanup**: Background task automatically deletes downloaded files after 1 hour to save space.
- **Modern UI**: React-based frontend for a smooth user experience.

## Tech Stack

### Backend
- **Go**: Core programming language.
- **Gin**: HTTP web framework for routing and API handling.
- **Gorilla WebSocket**: For real-time client communication.

### Frontend
- **React**: UI library.
- **Vite**: Build tool and development server.
- **Axios**: HTTP client.

### External Tools
- **yt-dlp**: Command-line audio/video downloader.
- **FFmpeg**: Complete, cross-platform solution to record, convert and stream audio and video.

## Prerequisites

Ensure you have the following installed on your system:
- [Go](https://go.dev/dl/) (1.18+)
- [Node.js](https://nodejs.org/) (16+)
- [yt-dlp](https://github.com/yt-dlp/yt-dlp/releases)
- [FFmpeg](https://ffmpeg.org/download.html)

> **Note**: `yt-dlp.exe` and `ffmpeg` executables should be placed in the root directory of the project or be available in your system's PATH.

## Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd yt-downloader
   ```

2. **Backend Setup:**
   Navigate to the root directory and install Go dependencies:
   ```bash
   go mod tidy
   ```
   Ensure `yt-dlp.exe` and `ffmpeg` are properly set up.

3. **Frontend Setup:**
   Navigate to the `web` directory and install Node dependencies:
   ```bash
   cd web
   npm install
   ```

## Usage

### Starting the Backend
From the project root directory, run:
```bash
go run cmd/main.go
```
The server will start on `http://localhost:8080`.

### Starting the Frontend
In a new terminal, navigate to the `web` directory and start the development server:
```bash
cd web
npm run dev
```
Open your browser and verify the local address (usually `http://localhost:5173`).

## API Documentation

### `POST /api/playlist/info`
Fetches metadata for a given YouTube URL (video or playlist).
- **Body**: `{"url": "https://youtube.com/..."}`
- **Response**: JSON object containing playlist title and entries.

### `POST /api/download`
Initiates the download process for a list of videos.
- **Body**: `{"urls": ["url1", "url2", ...]}`
- **Response**: `{"message": "Download started", "count": N}`

### `GET /ws`
WebSocket endpoint for listening to download progress events.
- **Events**:
    - `status`: "downloading" | "finished" | "error"
    - `percent`: Current download percentage.
    - `video_id`: ID/URL of the video.

## Project Structure

```
yt-downloader/
├── cmd/
│   └── main.go           # Entry point for the application
├── internal/
│   ├── downloader/       # Wrapper around yt-dlp
│   ├── handlers/         # HTTP and WebSocket handlers
│   └── models/           # Data structures
├── web/                  # React Frontend
│   ├── src/
│   └── public/
├── downloads/            # Temp folder for downloaded files
├── yt-dlp.exe            # External tool (optional if in PATH)
└── README.md             # Project documentation
```

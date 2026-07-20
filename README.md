# Premium YT & IG Media Downloader

[![Go Version](https://img.shields.io/badge/Go-1.25.0-blue.svg?style=flat-square&logo=go)](https://golang.org)
[![React Version](https://img.shields.io/badge/React-19.2.0-blue.svg?style=flat-square&logo=react)](https://react.dev)
[![Docker Support](https://img.shields.io/badge/Docker-Supported-blue.svg?style=flat-square&logo=docker)](https://www.docker.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)

A self-hosted, high-performance, and feature-rich media downloader for YouTube and Instagram. This project consists of a Go backend serving as a robust streaming proxy and controller, a sleek and interactive React single-page web application, and a Telegram Bot interface allowing direct downloads in messenger chats.

---

## 🌟 Key Features

*   **⚡ High-Speed Direct Streaming:** Streams video files from YouTube directly to the browser client without keeping temporary files on the host disk (except for MP3 conversions).
*   **🤖 Integrated Telegram Bot:** Send YouTube/Instagram links in Telegram and get direct MP3/MP4 files delivered as chat attachment uploads (automatically formats and keeps file sizes under Telegram's 50MB limit).
*   **🔑 User Accounts & Authentication:** Robust user registration and login system powered by JWT authentication and secure bcrypt password hashing.
*   **⚙️ Default User Preferences:** Set your preferred download format (Video/MP3) and maximum video quality (up to 4K / 2160p) inside your profile to apply settings automatically on every session.
*   **📂 Synced Download History:** Logs and saves your download history in a PostgreSQL database, letting you track and access past downloads instantly across any device.
*   **🔒 IP-Based Rate Limiting:** Built-in security middleware restricting abusive requests to heavy operations (playlist fetching, auth registration/login, streaming) to guarantee server uptime and prevent IP bans.
*   **🍪 Extractor Integrity (`cookies.txt`):** Includes mounting capabilities for YouTube cookies (`cookies.txt` in Netscape format) to bypass sign-in checks and rate-limit obstacles imposed by YouTube extractor changes.
*   **🐳 Fully Containerized:** Built and deployed in seconds using Docker Compose, bundling the React frontend, Nginx reverse proxy, Go backend server, and a PostgreSQL database.

---

## 🛠️ Technology Stack

### Backend
*   **Programming Language:** Go (v1.25.0)
*   **Web Framework:** [Gin Gonic](https://github.com/gin-gonic/gin)
*   **Rate Limiter:** `golang.org/x/time/rate` (Token Bucket Algorithm)
*   **Telegram Bot SDK:** [go-telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api)
*   **CLI Extractor:** [yt-dlp](https://github.com/yt-dlp/yt-dlp) + [Deno](https://deno.com/) (JS runtime to solve extractor challenges) + [ffmpeg](https://ffmpeg.org/) (media format transcoder)
*   **Database Migrations:** Embedded SQL migrator in Go code (`internal/db/migrations`)

### Frontend
*   **Framework:** React (v19.2.0)
*   **Build Tool:** Vite (v7.2.4)
*   **Client Communication:** Axios (v1.13.2)
*   **Styling:** Custom responsive CSS layout with modern typography, glassmorphism panel overlays, and fluid dark-mode theme variables.

### Database & Web Server
*   **Database:** PostgreSQL 16 (Alpine-based Docker container)
*   **Reverse Proxy / File server:** Nginx 1.25 (Alpine-based Docker container)

---

## 📂 Directory Structure

```text
├── cmd
│   └── main.go                  # Go Backend Entry Point
├── internal
│   ├── db                       # DB Connection & Embedded SQL Migrations
│   ├── downloader               # Executable wrappers for yt-dlp commands & format maps
│   ├── handlers                 # HTTP handlers for auth, playlist info, and stream endpoints
│   ├── middleware               # CORS, JWT Auth, and IP Rate Limiting middleware
│   ├── models                   # Go structs mapping database models & API requests
│   └── telegram                 # Telegram Bot client, routers, and request limiters
├── web
│   ├── src
│   │   ├── App.jsx              # Main React Application
│   │   ├── index.css            # Clean, Responsive CSS Styling & Profile Drawer Styles
│   │   └── main.jsx             # React DOM Mounting
│   ├── nginx.conf               # Nginx routing config (API proxying, buffering off, timeouts)
│   ├── Dockerfile.frontend      # Multi-stage production build for React UI in Nginx
│   └── package.json             # Node dependencies and scripts
├── Dockerfile.backend           # Go Builder & Runtime with yt-dlp, deno, ffmpeg, python3
├── docker-compose.yml           # Multi-container orchestration config
├── DEPLOY.md                    # Production deployment instructions
└── .env.example                 # Template for configuration settings
```

---

## 🚀 Getting Started

### 1. Prerequisites
Make sure you have the following installed on your host system:
*   [Docker](https://docs.docker.com/get-docker/) (v20.10+)
*   [Docker Compose](https://docs.docker.com/compose/install/) (v2.0+)

### 2. Environment Variables Configuration
Copy the template `.env.example` file to `.env` in the root folder:
```bash
cp .env.example .env
```
Open the `.env` file and set the variables:
*   `PORT`: The port on which the web application will be exposed (defaults to `80`).
*   `TELEGRAM_BOT_TOKEN`: Enter your Telegram Bot Token if you want the bot interface to run.
*   `JWT_SECRET`: Input a long random string used to secure user JWT sessions.
*   `DB_USER`, `DB_PASSWORD`, `DB_NAME`: Database credentials (default: `postgres`).

### 3. Setup YouTube Cookies (Critical)
To prevent YouTube extractor errors (such as *"Sign in to confirm you are not a bot"* or rate limits):
1.  Install a browser extension like **Get cookies.txt LOCALLY** (for Chrome/Brave/Firefox).
2.  Log in to your YouTube account in the browser.
3.  Click the extension and export the cookies in **Netscape format**.
4.  Rename the downloaded file to `cookies.txt` and place it in the root folder of this project.
5.  In production, this file is automatically mounted into the backend container at `/app/cookies.txt`.

### 4. Build and Run (Docker Compose)
Start the entire stack using Docker Compose:

```bash
# 1. Build the Docker images
docker compose build

# 2. Start the services in detached mode (background)
docker compose up -d

# 3. Check container status
docker compose ps
```

The application will be accessible at `http://localhost` (or your configured port/domain). 

To view logs for troubleshooting:
```bash
docker compose logs -f backend
```

---

## 🛠️ Local Development (Without Docker)

To run the backend and frontend separately for development:

### Running the Backend
Ensure you have `go 1.25+`, `python3`, `ffmpeg`, and `yt-dlp` installed on your local path.

1.  Run a local PostgreSQL database instance and configure your environment variables.
2.  Execute the main entry point:
    ```bash
    go run cmd/main.go
    ```
    The server will run on `http://localhost:8080` and automatically create/migrate tables on startup.

### Running the Frontend
1.  Navigate to the `web` folder:
    ```bash
    cd web
    ```
2.  Install development dependencies:
    ```bash
    npm install
    ```
3.  Start the Vite dev server:
    ```bash
    npm run dev
    ```
    The UI will run on `http://localhost:5173`. Make sure to configure CORS origins or `.env` if custom URLs are used.

---

## 🔒 Security & Performance Details

### ⚡ Proxy Buffering Off
For streaming video directly from YouTube to clients, Nginx is configured with `proxy_buffering off;`. This ensures that downstream users receive the stream chunks in real-time, decreasing server memory load and eliminating delay buffers on the gateway.

### ⏱️ Timeouts
Large media streams might take several minutes to download and stream. Nginx read/write timeouts (`proxy_read_timeout`, `proxy_send_timeout`) are set to `600s` (10 minutes) under `web/nginx.conf` to avoid premature connection closures.

### 🛡️ Rate Limiting Rules
*   **Web Client Core API:** Max 5 requests per minute (with burst room) on heavy routes (e.g. `/api/playlist/info`).
*   **User Registration / Login:** Max 10 attempts per minute per IP to mitigate brute-force attacks.
*   **Telegram Bot Requests:** Max 3 commands/callbacks per minute per Telegram ID.

---

## 🔄 Maintenance & Extractor Updates

YouTube changes its API and website codebase frequently, which can occasionally break the download mechanism. If downloads begin to fail, rebuild the backend Docker image without cache to pull the latest version of `yt-dlp`:

```bash
docker compose build --no-cache backend
docker compose up -d backend
```

---

## 📄 License

This project is open-source and licensed under the [MIT License](LICENSE).

# Production Deployment Guide

This guide describes how to deploy the YouTube Downloader application to a production server using Docker Compose.

---

## 1. Prerequisites

Ensure your target server has the following installed:
- [Docker](https://docs.docker.com/get-docker/) (v20.10+)
- [Docker Compose](https://docs.docker.com/compose/install/) (v2.0+)

---

## 2. Environment Configuration

1. Copy the `.env.example` template to `.env` in the root folder:
   ```bash
   cp .env.example .env
   ```
2. Open `.env` and configure the settings:
   - **`PORT`**: Set the public-facing port (defaults to `80`).
   - **`TELEGRAM_BOT_TOKEN`**: Provide your Telegram Bot API token if you want to use the bot interface.

---

## 3. YouTube Cookies (`cookies.txt`)

To prevent rate-limiting, IP bans, or "Sign in to confirm you are not a bot" errors from YouTube, you **must** supply a valid `cookies.txt` file from a logged-in YouTube account.

1. Install a browser extension like **Get cookies.txt LOCALLY** (for Chrome/Brave/Firefox).
2. Log into YouTube in your browser.
3. Open the extension and export your cookies as a Netscape-format text file.
4. Rename the exported file to `cookies.txt` and place it in the root folder of this project.
5. In production, this file is automatically mounted into the backend container at `/app/cookies.txt`.

---

## 4. Run the Application

Once your `.env` and `cookies.txt` are ready, run:

1. **Build the Docker containers:**
   ```bash
   docker compose build
   ```
2. **Start the services in the background (detached mode):**
   ```bash
   docker compose up -d
   ```
3. **Verify the status of the containers:**
   ```bash
   docker compose ps
   ```

---

## 5. Monitoring & Logs

To check if the services are running correctly:

- **Check logs:**
  ```bash
  docker compose logs -f
  ```
- **Check backend logs only:**
  ```bash
  docker compose logs -f backend
  ```
- **Check frontend (Nginx proxy) logs only:**
  ```bash
  docker compose logs -f frontend
  ```

---

## 6. Maintenance & Updates

### Updating `yt-dlp`
YouTube changes its platform code frequently. If downloads start failing, rebuild the backend image to fetch the latest version of `yt-dlp`:
```bash
docker compose build --no-cache backend
docker compose up -d backend
```

### Temporary File Cleanup
The application deletes temporary MP3 conversions and Telegram-sent files automatically. However, you can monitor the local `downloads` directory to make sure it doesn't build up unneeded cache.

package telegram

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/example/yt-downloader/internal/db"
	"github.com/example/yt-downloader/internal/downloader"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	db      *sql.DB
	limiter *BotLimiter
}

func NewBot(token string, database *sql.DB) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:     api,
		db:      database,
		limiter: NewBotLimiter(),
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		}
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	if msg.From != nil {
		// Rate limit: 3 requests per minute
		if !b.limiter.Allow(msg.From.ID, 3, 1*time.Minute) {
			reply := tgbotapi.NewMessage(msg.Chat.ID, "So'rovlar soni ko'payib ketdi. Iltimos, birozdan keyin qaytadan urining (Limit: daqiqasiga 3 ta so'rov).")
			b.api.Send(reply)
			return
		}

		err := db.UpsertUser(b.db, msg.From.ID, msg.From.FirstName, msg.From.LastName, msg.From.UserName, msg.From.LanguageCode, msg.From.IsBot)
		if err != nil {
			log.Printf("Failed to upsert user in database: %v", err)
		}
	}

	if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			reply := tgbotapi.NewMessage(msg.Chat.ID, "Assalomu alaykum! YouTube yoki Instagram havolasini yuboring va men uni yuklab beraman.")
			b.api.Send(reply)
		}
		return
	}

	url := strings.TrimSpace(msg.Text)
	isYouTube := strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be")
	isInstagram := strings.Contains(url, "instagram.com") || strings.Contains(url, "ddinstagram.com")
	if !isYouTube && !isInstagram {
		return
	}

	// Notify user that we are fetching info
	waitMsg := tgbotapi.NewMessage(msg.Chat.ID, "Video ma'lumotlari yuklanmoqda...")
	sentWait, _ := b.api.Send(waitMsg)

	playlist, err := downloader.GetPlaylistInfo(url)
	if err != nil {
		edit := tgbotapi.NewEditMessageText(msg.Chat.ID, sentWait.MessageID, fmt.Sprintf("Xatolik: %v", err))
		b.api.Send(edit)
		return
	}

	if len(playlist.Entries) == 0 {
		edit := tgbotapi.NewEditMessageText(msg.Chat.ID, sentWait.MessageID, "Hech qanday video topilmadi.")
		b.api.Send(edit)
		return
	}

	video := playlist.Entries[0]
	text := fmt.Sprintf("Nomi: %s\nDavomiyligi: %.0f soniya\n\nQaysi formatda yuklamoqchisiz?", video.Title, video.Duration)

	// Callback data: format|quality|source_id
	// Telegram callback data has a limit of 64 bytes. We prepend 'yt:' or 'ig:' to identify the platform.
	prefix := "yt:"
	if isInstagram {
		prefix = "ig:"
	}
	mp4Data := fmt.Sprintf("mp4|best|%s%s", prefix, video.ID)
	mp3Data := fmt.Sprintf("mp3|best|%s%s", prefix, video.ID)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Video (MP4)", mp4Data),
			tgbotapi.NewInlineKeyboardButtonData("Audio (MP3)", mp3Data),
		),
	)

	edit := tgbotapi.NewEditMessageText(msg.Chat.ID, sentWait.MessageID, text)
	edit.ReplyMarkup = &keyboard
	b.api.Send(edit)
}

func (b *Bot) handleCallback(query *tgbotapi.CallbackQuery) {
	if query.From != nil {
		// Rate limit: 3 requests per minute
		if !b.limiter.Allow(query.From.ID, 3, 1*time.Minute) {
			b.api.Send(tgbotapi.NewCallback(query.ID, "So'rov cheklovi! Iltimos kuting."))
			reply := tgbotapi.NewMessage(query.Message.Chat.ID, "So'rovlar soni ko'payib ketdi. Iltimos, birozdan keyin qaytadan urining (Limit: daqiqasiga 3 ta so'rov).")
			b.api.Send(reply)
			return
		}
	}

	data := strings.Split(query.Data, "|")
	if len(data) < 3 {
		return
	}

	format := data[0]
	quality := data[1]
	videoID := data[2]

	var url string
	if strings.HasPrefix(videoID, "ig:") {
		shortcode := strings.TrimPrefix(videoID, "ig:")
		url = fmt.Sprintf("https://www.instagram.com/p/%s/", shortcode)
	} else if strings.HasPrefix(videoID, "yt:") {
		ytID := strings.TrimPrefix(videoID, "yt:")
		url = fmt.Sprintf("https://www.youtube.com/watch?v=%s", ytID)
	} else {
		// Backwards compatibility for old buttons
		url = fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	}

	// Safe ID for output filename
	safeID := strings.ReplaceAll(videoID, ":", "_")

	platform := "youtube"
	if strings.HasPrefix(videoID, "ig:") {
		platform = "instagram"
	}

	if query.From != nil {
		err := db.UpsertUser(b.db, query.From.ID, query.From.FirstName, query.From.LastName, query.From.UserName, query.From.LanguageCode, query.From.IsBot)
		if err != nil {
			log.Printf("Failed to upsert user on callback: %v", err)
		}
		err = db.LogDownload(b.db, query.From.ID, platform, url, format)
		if err != nil {
			log.Printf("Failed to log download on callback: %v", err)
		}
	}

	// Answer callback to remove loading state
	b.api.Send(tgbotapi.NewCallback(query.ID, "Yuklab olish boshlandi..."))

	// Create temporary file
	tempDir := "/app/downloads"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		tempDir = "./downloads" // Fallback for local dev
		os.MkdirAll(tempDir, 0777)
	}

	ext := "mp4"
	if format == "mp3" {
		ext = "mp3"
	}
	outputPath := filepath.Join(tempDir, fmt.Sprintf("%s_%s.%s", safeID, format, ext))

	// Notify user
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, "Fayl tayyorlanmoqda, iltimos kuting...")
	sentMsg, _ := b.api.Send(msg)

	err := downloader.DownloadToPath(url, format, quality, outputPath)
	if err != nil {
		edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, sentMsg.MessageID, fmt.Sprintf("Yuklab olishda xatolik: %v", err))
		b.api.Send(edit)
		return
	}

	// Check file size
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, sentMsg.MessageID, "Faylni tekshirishda xatolik yuz berdi.")
		b.api.Send(edit)
		return
	}

	if fileInfo.Size() > 50*1024*1024 {
		edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, sentMsg.MessageID, "Fayl hajmi 50MB dan katta. Uni bot orqali yuborib bo'lmaydi. Iltimos, veb-saytdan foydalaning.")
		b.api.Send(edit)
		// We keep the file in downloads directory so it can be accessed via web if needed
		return
	}

	// Send file
	b.api.Send(tgbotapi.NewDeleteMessage(query.Message.Chat.ID, sentMsg.MessageID))

	if format == "mp3" {
		audio := tgbotapi.NewAudio(query.Message.Chat.ID, tgbotapi.FilePath(outputPath))
		_, err = b.api.Send(audio)
	} else {
		video := tgbotapi.NewVideo(query.Message.Chat.ID, tgbotapi.FilePath(outputPath))
		_, err = b.api.Send(video)
	}

	if err != nil {
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, fmt.Sprintf("Faylni yuborishda xatolik: %v", err))
		b.api.Send(msg)
	}

	// Cleanup
	os.Remove(outputPath)
}

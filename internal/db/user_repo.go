package db

import (
	"database/sql"
	"time"
)

// User represents a user record in the database
type User struct {
	ID                int64     `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Username          string    `json:"username"`
	LanguageCode      string    `json:"language_code"`
	IsBot             bool      `json:"is_bot"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	LastInteractionAt time.Time `json:"last_interaction_at"`
	InteractionCount  int       `json:"interaction_count"`
}

// UpsertUser registers a user or updates their profile attributes and bumps their interaction count
func UpsertUser(db *sql.DB, id int64, firstName, lastName, username, languageCode string, isBot bool) error {
	query := `
		INSERT INTO users (id, first_name, last_name, username, language_code, is_bot, created_at, updated_at, last_interaction_at, interaction_count)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), NOW(), 1)
		ON CONFLICT (id) DO UPDATE 
		SET first_name = EXCLUDED.first_name,
		    last_name = EXCLUDED.last_name,
		    username = EXCLUDED.username,
		    language_code = EXCLUDED.language_code,
		    updated_at = NOW(),
		    last_interaction_at = NOW(),
		    interaction_count = users.interaction_count + 1
	`
	_, err := db.Exec(query, id, firstName, lastName, username, languageCode, isBot)
	return err
}

// LogDownload logs a successfully requested download audit record
func LogDownload(db *sql.DB, userID int64, platform, url, format string) error {
	query := `
		INSERT INTO user_downloads (user_id, platform, url, format, downloaded_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := db.Exec(query, userID, platform, url, format)
	return err
}

// WebUser represents a row in the web_users table
type WebUser struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	AvatarColor    string    `json:"avatar_color"`
	DefaultFormat  string    `json:"default_format"`
	DefaultQuality string    `json:"default_quality"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// DownloadLog represents a logged download event
type DownloadLog struct {
	ID           int       `json:"id"`
	WebUserID    *int      `json:"web_user_id,omitempty"`
	Platform     string    `json:"platform"`
	URL          string    `json:"url"`
	Format       string    `json:"format"`
	DownloadedAt time.Time `json:"downloaded_at"`
}

// CreateWebUser inserts a new registered user in the database
func CreateWebUser(db *sql.DB, username, email, passwordHash string) (int, error) {
	var id int
	query := `
		INSERT INTO web_users (username, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`
	err := db.QueryRow(query, username, email, passwordHash).Scan(&id)
	return id, err
}

// GetWebUserByUsernameOrEmail fetches a web user for login authentication
func GetWebUserByUsernameOrEmail(db *sql.DB, identifier string) (*WebUser, error) {
	query := `
		SELECT id, username, email, password_hash, avatar_color, default_format, default_quality, created_at, updated_at
		FROM web_users
		WHERE username = $1 OR email = $1
	`
	var user WebUser
	err := db.QueryRow(query, identifier).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.AvatarColor, &user.DefaultFormat, &user.DefaultQuality,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetWebUserByID fetches a web user by primary key
func GetWebUserByID(db *sql.DB, id int) (*WebUser, error) {
	query := `
		SELECT id, username, email, password_hash, avatar_color, default_format, default_quality, created_at, updated_at
		FROM web_users
		WHERE id = $1
	`
	var user WebUser
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.AvatarColor, &user.DefaultFormat, &user.DefaultQuality,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateWebUserPreferences updates the profile name/email and settings
func UpdateWebUserPreferences(db *sql.DB, id int, username, email, avatarColor, format, quality string) error {
	query := `
		UPDATE web_users
		SET username = $1, email = $2, avatar_color = $3, default_format = $4, default_quality = $5, updated_at = NOW()
		WHERE id = $6
	`
	_, err := db.Exec(query, username, email, avatarColor, format, quality, id)
	return err
}

// LogWebDownload logs a media download associated with a web user
func LogWebDownload(db *sql.DB, webUserID int, platform, url, format string) error {
	query := `
		INSERT INTO user_downloads (web_user_id, platform, url, format, downloaded_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := db.Exec(query, webUserID, platform, url, format)
	return err
}

// GetWebUserDownloadHistory retrieves recent download audits for a web user
func GetWebUserDownloadHistory(db *sql.DB, webUserID int) ([]DownloadLog, error) {
	query := `
		SELECT id, web_user_id, platform, url, format, downloaded_at
		FROM user_downloads
		WHERE web_user_id = $1
		ORDER BY downloaded_at DESC
		LIMIT 50
	`
	rows, err := db.Query(query, webUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []DownloadLog
	for rows.Next() {
		var l DownloadLog
		err := rows.Scan(&l.ID, &l.WebUserID, &l.Platform, &l.URL, &l.Format, &l.DownloadedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

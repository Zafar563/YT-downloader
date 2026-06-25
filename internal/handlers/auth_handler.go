package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/example/yt-downloader/internal/db"
	"github.com/example/yt-downloader/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db *sql.DB
}

func NewAuthHandler(database *sql.DB) *AuthHandler {
	return &AuthHandler{db: database}
}

type authClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func (h *AuthHandler) generateToken(userID int) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "yoursupersecretjwtkeyhere"
	}
	secretKey := []byte(jwtSecret)

	claims := authClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days expiration
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to secure password"})
		return
	}
	passwordHash := string(hashedBytes)

	// Create user
	userID, err := db.CreateWebUser(h.db, req.Username, req.Email, passwordHash)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "violates unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or Email is already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user record"})
		return
	}

	// Generate JWT
	token, err := h.generateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User: models.WebUserResponse{
			ID:             userID,
			Username:       req.Username,
			Email:          req.Email,
			AvatarColor:    "#6366f1",
			DefaultFormat:  "video",
			DefaultQuality: "best",
		},
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve user profile
	user, err := db.GetWebUserByUsernameOrEmail(h.db, req.UsernameOrEmail)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/email or password"})
		return
	}

	// Match password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/email or password"})
		return
	}

	// Generate JWT
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User: models.WebUserResponse{
			ID:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			AvatarColor:    user.AvatarColor,
			DefaultFormat:  user.DefaultFormat,
			DefaultQuality: user.DefaultQuality,
		},
	})
}

// GetProfile returns the authenticated user details and their download history
func (h *AuthHandler) GetProfile(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := val.(int)

	user, err := db.GetWebUserByID(h.db, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	history, err := db.GetWebUserDownloadHistory(h.db, userID)
	if err != nil {
		history = []db.DownloadLog{} // fallback to empty list
	}

	c.JSON(http.StatusOK, gin.H{
		"user": models.WebUserResponse{
			ID:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			AvatarColor:    user.AvatarColor,
			DefaultFormat:  user.DefaultFormat,
			DefaultQuality: user.DefaultQuality,
		},
		"history": history,
	})
}

type UpdateProfileReq struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	AvatarColor    string `json:"avatar_color"`
	DefaultFormat  string `json:"default_format"`
	DefaultQuality string `json:"default_quality"`
}

// UpdateProfile edits user profile and preferences
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := val.(int)

	var req UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.UpdateWebUserPreferences(h.db, userID, req.Username, req.Email, req.AvatarColor, req.DefaultFormat, req.DefaultQuality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile settings"})
		return
	}

	c.JSON(http.StatusOK, models.WebUserResponse{
		ID:             userID,
		Username:       req.Username,
		Email:          req.Email,
		AvatarColor:    req.AvatarColor,
		DefaultFormat:  req.DefaultFormat,
		DefaultQuality: req.DefaultQuality,
	})
}

type LogDownloadReq struct {
	Platform string `json:"platform" binding:"required"`
	URL      string `json:"url" binding:"required"`
	Format   string `json:"format" binding:"required"`
}

// LogDownload records a successfully completed download associated with the logged-in user
func (h *AuthHandler) LogDownload(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := val.(int)

	var req LogDownloadReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.LogWebDownload(h.db, userID, req.Platform, req.URL, req.Format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to audit download"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

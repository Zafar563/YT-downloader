package models

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// WebUserResponse represents the public user profile profile info
type WebUserResponse struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	AvatarColor    string `json:"avatar_color"`
	DefaultFormat  string `json:"default_format"`
	DefaultQuality string `json:"default_quality"`
}

// AuthResponse represents the response containing profile info and token
type AuthResponse struct {
	Token string          `json:"token"`
	User  WebUserResponse `json:"user"`
}

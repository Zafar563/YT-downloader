package telegram

import (
	"sync"
	"time"
)

// BotLimiter regulates request frequencies for Telegram bot users
type BotLimiter struct {
	mu       sync.Mutex
	requests map[int64][]time.Time
}

// NewBotLimiter creates a new Telegram request rate limiter
func NewBotLimiter() *BotLimiter {
	return &BotLimiter{
		requests: make(map[int64][]time.Time),
	}
}

// Allow checks if a user is within the limit (e.g. limit requests per duration window)
func (l *BotLimiter) Allow(userID int64, limit int, duration time.Duration) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-duration)

	// Keep only timestamps within the sliding window duration
	var active []time.Time
	for _, t := range l.requests[userID] {
		if t.After(cutoff) {
			active = append(active, t)
		}
	}

	if len(active) >= limit {
		l.requests[userID] = active
		return false
	}

	active = append(active, now)
	l.requests[userID] = active
	return true
}

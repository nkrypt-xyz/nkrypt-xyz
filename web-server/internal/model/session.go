package model

import "time"

// Session represents a row in the sessions audit table.
type Session struct {
	ID          string
	UserID      string
	APIKeyHash  string
	HasExpired  bool
	ExpiredAt   *time.Time
	ExpireReason *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RedisSessionData is stored as JSON in Redis.
type RedisSessionData struct {
	SessionID string `json:"sessionId"`
	UserID    string `json:"userId"`
	CreatedAt int64  `json:"createdAt"`
}

// SessionListItem is used for /user/list-all-sessions response.
type SessionListItem struct {
	IsCurrentSession bool    `json:"isCurrentSession"`
	HasExpired       bool    `json:"hasExpired"`
	ExpireReason     *string `json:"expireReason"`
	CreatedAt        int64   `json:"createdAt"`
	ExpiredAt        *int64  `json:"expiredAt"`
}


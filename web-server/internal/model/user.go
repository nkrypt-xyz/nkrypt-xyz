package model

import "time"

// User represents the users table.
type User struct {
	ID                   string
	DisplayName          string
	UserName             string
	PasswordHash         string
	PasswordSalt         string
	IsBanned             bool
	PermManageAllUser    bool
	PermCreateUser       bool
	PermCreateBucket     bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
	GlobalPermissionsMap map[string]bool // derived, not stored directly
}

// UserListItem is a lightweight view for listing users.
type UserListItem struct {
	ID          string
	UserName    string
	DisplayName string
	IsBanned    bool
}

// AuthData is stored in context for authenticated requests.
type AuthData struct {
	ApiKey    string
	UserID    string
	SessionID string
	User      *User
}


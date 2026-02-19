package model

// UserResponse is the JSON representation of a user in API responses.
type UserResponse struct {
	ID                string          `json:"_id"`
	UserName          string          `json:"userName"`
	DisplayName       string          `json:"displayName"`
	IsBanned          bool            `json:"isBanned"`
	GlobalPermissions map[string]bool `json:"globalPermissions,omitempty"`
}

// SessionResponse is the minimal session representation used in login/assert.
type SessionResponse struct {
	ID string `json:"_id"`
}


package models

import "time"

type OAuthProvider struct {
	UserID         int64
	ProviderName   string
	ProviderUserID string
	AccessToken    string
	RefreshToken   string
	ExpiresAt      time.Time
}

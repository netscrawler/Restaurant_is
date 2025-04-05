package models

import "time"

type OAuthProvider struct {
	ClientID    string
	Provider    string // google, yandex
	ProviderID  string
	AccessToken string
	ExpiresAt   time.Time
}

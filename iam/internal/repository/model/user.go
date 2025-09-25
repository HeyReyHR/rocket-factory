package model

import "time"

type User struct {
	Uuid      string `json:"uuid"`
	UserInfo  `json:"user_info"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type UserInfo struct {
	Login               string               `json:"login"`
	Email               string               `json:"email"`
	PasswordHash        string               `json:"password_hash"`
	NotificationMethods []NotificationMethod `json:"notification_methods"`
}

type NotificationMethod struct {
	ProviderName string `json:"provider_name"`
	Target       string `json:"target"`
}

type Session struct {
	Uuid      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at "`
}

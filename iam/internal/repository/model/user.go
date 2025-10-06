package model

import (
	"time"
)

type User struct {
	Uuid           string             `json:"uuid"`
	AdditionalInfo `json:"user_info"` // going to keep in mind that
	UpdatedAt      time.Time          `json:"updated_at"`
	CreatedAt      time.Time          `json:"created_at"`
}

type AdditionalInfo struct {
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

type UserDbDto struct {
	Uuid      string    `db:"uuid"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Login        string `db:"login"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`

	NotificationMethods []string `db:"notification_methods"`
}

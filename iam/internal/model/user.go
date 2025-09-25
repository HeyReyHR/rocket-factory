package model

import "time"

type User struct {
	Uuid string
	UserInfo
	UpdatedAt time.Time
	CreatedAt time.Time
}

type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

type NotificationMethod struct {
	ProviderName string
	Target       string
}

type TokenPair struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type Claims struct {
	Uuid  int64  `json:"uuid"`
	Login string `json:"login"`
}

type Session struct {
	Uuid      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

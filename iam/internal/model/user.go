package model

import "time"

type User struct {
	Uuid string
	AdditionalInfo
	UpdatedAt time.Time
	CreatedAt time.Time
}

type AdditionalInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

type NotificationMethod struct {
	ProviderName string
	Target       string
}

type Session struct {
	Uuid      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

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
	PasswordHash        string
	NotificationMethods []NotificationMethod
}

type NotificationMethod struct {
	ProviderName string
	Target       string
}

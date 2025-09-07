package model

import (
	"errors"
)

var (
	ErrEventNotFound   = errors.New("event not found")
	ErrEventScanFailed = errors.New("event scan failed")
)

package auth

import (
	"fmt"
	"time"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/cache"
)

const (
	sessionsCacheKeyPrefix = "auth:sessions:"
	userCacheKeyPrefix     = "auth:user-sessions:"
)

type repository struct {
	cache      cache.RedisClient
	sessionTtl time.Duration
}

func NewRepository(cache cache.RedisClient, sessionTtl time.Duration) *repository {
	return &repository{
		cache:      cache,
		sessionTtl: sessionTtl,
	}
}

func (r *repository) getSessionCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", sessionsCacheKeyPrefix, uuid)
}

func (r *repository) getUserCacheKey(userUuid string) string {
	return fmt.Sprintf("%s%s", userCacheKeyPrefix, userUuid)
}

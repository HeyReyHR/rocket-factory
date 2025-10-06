package auth

import (
	"fmt"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/cache"
)

const (
	sessionsCacheKeyPrefix = "auth:sessions:"
	userCacheKeyPrefix     = "auth:user-sessions:"
)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

func (r *repository) getSessionCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", sessionsCacheKeyPrefix, uuid)
}

func (r *repository) getUserCacheKey(userUuid string) string {
	return fmt.Sprintf("%s%s", userCacheKeyPrefix, userUuid)
}

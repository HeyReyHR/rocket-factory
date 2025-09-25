package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
	"github.com/google/uuid"
)

func (r *repository) Create(ctx context.Context, user model.User) (string, error) {
	sessionUuid := uuid.NewString()
	now := time.Now().UTC()

	session := model.Session{
		Uuid:      sessionUuid,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(r.sessionTtl),
	}
	payload, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	payloadUser, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	if err = r.cache.SetWithTTL(ctx, r.getSessionCacheKey(sessionUuid), payload, r.sessionTtl); err != nil {
		return "", err
	}

	if err = r.cache.SetWithTTL(ctx, r.getUserCacheKey(sessionUuid), payloadUser, r.sessionTtl); err != nil {
		_ = r.cache.Del(ctx, r.getSessionCacheKey(sessionUuid))
		return "", err
	}

	if err = r.cache.SAdd(ctx, r.getUserCacheKey(user.Uuid), sessionUuid); err != nil {
		_ = r.cache.Del(ctx, r.getSessionCacheKey(sessionUuid))
		_ = r.cache.Del(ctx, r.getUserCacheKey(user.Uuid))
		return "", err
	}

	return sessionUuid, nil
}

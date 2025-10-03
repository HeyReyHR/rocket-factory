package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, user model.User, session model.Session, sessionTtl time.Duration) error {
	payload, err := json.Marshal(session)
	if err != nil {
		return err
	}
	payloadUser, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err = r.cache.SetWithTTL(ctx, r.getSessionCacheKey(session.Uuid), payload, sessionTtl); err != nil {
		return err
	}

	if err = r.cache.SetWithTTL(ctx, r.getUserCacheKey(session.Uuid), payloadUser, sessionTtl); err != nil {
		_ = r.cache.Del(ctx, r.getSessionCacheKey(session.Uuid))
		return err
	}

	if err = r.cache.SAdd(ctx, r.getUserCacheKey(user.Uuid), session.Uuid); err != nil {
		_ = r.cache.Del(ctx, r.getSessionCacheKey(session.Uuid))
		_ = r.cache.Del(ctx, r.getUserCacheKey(user.Uuid))
		return err
	}

	return nil
}

package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
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
		if err = r.cache.Del(ctx, r.getSessionCacheKey(session.Uuid)); err != nil {
			logger.Error(ctx, "Failed to delete cache", zap.Error(err))
		}
		return err
	}

	if err = r.cache.SAdd(ctx, r.getUserCacheKey(user.Uuid), session.Uuid); err != nil {
		if err = r.cache.Del(ctx, r.getSessionCacheKey(session.Uuid)); err != nil {
			logger.Error(ctx, "Failed to delete cache", zap.Error(err))
		}
		if err = r.cache.Del(ctx, r.getUserCacheKey(user.Uuid)); err != nil {
			logger.Error(ctx, "Failed to delete cache", zap.Error(err))
		}
		return err
	}

	return nil
}

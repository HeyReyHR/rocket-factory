package auth

import (
	"context"
	"encoding/json"

	serviceModel "github.com/HeyReyHR/rocket-factory/iam/internal/model"
	"github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, sessionUuid string) (model.User, model.Session, error) {
	sessionRaw, err := r.cache.Get(ctx, r.getSessionCacheKey(sessionUuid))
	if err != nil {
		if sessionRaw == nil {
			return model.User{}, model.Session{}, serviceModel.ErrSessionNotFound
		}
		return model.User{}, model.Session{}, err
	}

	var session model.Session
	if err = json.Unmarshal(sessionRaw, &session); err != nil {
		return model.User{}, model.Session{}, err
	}

	userRaw, err := r.cache.Get(ctx, r.getUserCacheKey(sessionUuid))
	if err != nil {
		return model.User{}, model.Session{}, err
	}

	var user model.User
	if err = json.Unmarshal(userRaw, &user); err != nil {
		return model.User{}, model.Session{}, err
	}
	return user, session, nil
}

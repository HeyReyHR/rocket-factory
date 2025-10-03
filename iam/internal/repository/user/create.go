package user

import (
	"context"
	"time"

	"github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *repository) Create(ctx context.Context, user model.AdditionalInfo, passwordHash string) (string, error) {
	userUuid := uuid.NewString()

	tx, err := r.dbConn.Begin(ctx)

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	if err != nil {
		logger.Error(ctx, "begin tx failed", zap.Error(err))
		return "", err
	}

	if _, err = tx.Exec(ctx, "INSERT INTO users (uuid, created_at, updated_at) VALUES ($1, $2, $3)", userUuid, time.Now(), time.Now()); err != nil {
		logger.Error(ctx, "insert user failed", zap.Error(err))
		return "", err
	}

	if _, err = tx.Exec(ctx, "INSERT INTO user_infos (user_uuid, login, email, password_hash) VALUES ($1, $2, $3, $4)", userUuid, user.Login, user.Email, passwordHash); err != nil {
		logger.Error(ctx, "insert user failed", zap.Error(err))
		return "", err
	}

	_, err = tx.Prepare(ctx, "ins_notification_methods",
		"INSERT INTO notification_methods (user_infos_uuid, provider_name, target) VALUES ($1, $2, $3)")
	if err != nil {
		return "", err
	}
	for _, notifMethod := range user.NotificationMethods {
		if _, err = tx.Exec(ctx, "ins_notification_methods", userUuid, notifMethod.ProviderName, notifMethod.Target); err != nil {
			return "", err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Error(ctx, "commit tx failed", zap.Error(err))
		_ = tx.Rollback(ctx)
		return "", err
	}

	committed = true
	logger.Info(ctx, "user created", zap.String("registered users", userUuid))
	return userUuid, nil
}

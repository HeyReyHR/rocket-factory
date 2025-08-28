package assembly

import (
	"context"
	"encoding/json"

	"github.com/HeyReyHR/rocket-factory/assembly/internal/repository/model"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (r *repository) Create(ctx context.Context, uuid string, eventType model.EventType, orderUuid string, userUuid string, buildTimeSec int64) error {
	payload := map[string]interface{}{
		"order_uuid":     orderUuid,
		"user_uuid":      userUuid,
		"build_time_sec": buildTimeSec,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		logger.Error(ctx, "marshal payload failed", zap.Error(err))
		return err
	}

	_, err = r.dbConn.Exec(ctx, "INSERT INTO outbox (uuid, type, payload, status) VALUES ($1, $2, $3::jsonb, $4)", uuid, eventType, b, model.PendingStatus)
	if err != nil {
		logger.Error(ctx, "create outbox record failed", zap.Error(err))
		return err
	}
	return nil
}

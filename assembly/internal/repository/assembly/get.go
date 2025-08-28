package assembly

import (
	"context"
	"errors"

	"github.com/HeyReyHR/rocket-factory/assembly/internal/model"
	repoModel "github.com/HeyReyHR/rocket-factory/assembly/internal/repository/model"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Get(ctx context.Context) (model.OrderAssembledEvent, error) {
	row := r.dbConn.QueryRow(ctx, `SELECT uuid, type, payload, status 
	FROM outbox WHERE status = $1 LIMIT 1`, repoModel.PendingStatus)

	var event repoModel.OrderAssembled
	err := row.Scan(&event.EventUuid, &event.EventType, &event.Payload, &event.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.OrderAssembledEvent{}, model.ErrEventNotFound
		}
		return model.OrderAssembledEvent{}, model.ErrEventScanFailed
	}

	orderAssembledEvent := model.OrderAssembledEvent{
		EventUuid:    event.EventUuid,
		OrderUuid:    event.Payload.OrderUuid,
		UserUuid:     event.Payload.UserUuid,
		BuildTimeSec: event.Payload.BuildTimeSec,
	}
	return orderAssembledEvent, nil
}

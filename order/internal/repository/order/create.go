package order

import (
	"context"
)

func (r *repository) Create(ctx context.Context, uuid, userUuid string, partUuids []string, totalPrice float64) (string, float64, error) {
	_, err := r.dbConn.Exec(ctx, "INSERT INTO orders (uuid, user_uuid, part_uuids, total_price) VALUES ($1, $2, $3, $4)", uuid, userUuid, partUuids, totalPrice)
	if err != nil {
		return "", 0, err
	}
	return uuid, totalPrice, nil
}

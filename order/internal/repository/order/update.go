package order

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/repository/model"
)

func (r *repository) Update(ctx context.Context, uuid string, order model.Order) error {
	_, err := r.dbConn.Exec(ctx,
		"UPDATE orders SET user_uuid = $1, part_uuids = $2, total_price = $3, transaction_uuid = $4, status = $5, payment_method = $6, updated_at = now() WHERE uuid = $7",
		order.UserUuid, order.PartUuids, order.TotalPrice, order.TransactionUuid, order.Status, order.PaymentMethod, uuid)
	if err != nil {
		return err
	}

	return nil
}

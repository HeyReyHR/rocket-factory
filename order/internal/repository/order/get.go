package order

import (
	"context"
	"errors"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/HeyReyHR/rocket-factory/order/internal/repository/model"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Order, error) {
	row := r.dbConn.QueryRow(ctx, `SELECT user_uuid, part_uuids, total_price, transaction_uuid, payment_method, status 
	FROM orders WHERE uuid = $1`, uuid)

	var userUuid string
	var partUuids []string
	var totalPrice float64
	var transactionUuid *string
	var paymentMethod *repoModel.PaymentMethod
	var status repoModel.Status
	err := row.Scan(&userUuid, &partUuids, &totalPrice, &transactionUuid, &paymentMethod, &status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, model.ErrOrderScanFailed
	}
	order := repoModel.Order{
		Uuid:            uuid,
		UserUuid:        userUuid,
		PartUuids:       partUuids,
		TotalPrice:      totalPrice,
		TransactionUuid: transactionUuid,
		PaymentMethod:   paymentMethod,
		Status:          status,
	}

	return converter.RepoOrderToServiceOrder(order), nil
}

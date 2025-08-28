package assembly

import (
	"context"

	repoModel "github.com/HeyReyHR/rocket-factory/assembly/internal/repository/model"
)

func (r *repository) Update(ctx context.Context, uuid string) error {
	_, err := r.dbConn.Exec(ctx,
		"UPDATE outbox SET status = $1 WHERE uuid = $2", repoModel.Done, uuid)
	if err != nil {
		return err
	}
	return nil
}

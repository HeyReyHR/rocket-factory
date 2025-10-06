package user

import (
	repository2 "github.com/HeyReyHR/rocket-factory/iam/internal/repository"
	"github.com/jackc/pgx/v5"
)

var _ repository2.UserRepository = (*repository)(nil)

type repository struct {
	dbConn *pgx.Conn
}

func NewRepository(dbConn *pgx.Conn) *repository {
	return &repository{
		dbConn: dbConn,
	}
}

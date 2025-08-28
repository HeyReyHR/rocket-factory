package assembly

import (
	repository2 "github.com/HeyReyHR/rocket-factory/assembly/internal/repository"
	"github.com/jackc/pgx/v5"
)

var _ repository2.AssemblyRepository = (*repository)(nil)

type repository struct {
	dbConn *pgx.Conn
}

func NewRepository(dbConn *pgx.Conn) *repository {
	return &repository{
		dbConn: dbConn,
	}
}

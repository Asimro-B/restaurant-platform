package persistencedb

import (
	"context"
	db "restaurant-platform/database/sqlc/gen"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
)

type PersistenceDB struct {
	*db.Queries
	Pool   *pgxpool.Pool
	GormDB *gorm.DB
}

func New(pool *pgxpool.Pool, GormDB *gorm.DB) *PersistenceDB {
	return &PersistenceDB{
		Queries: db.New(pool),
		Pool:    pool,
		GormDB:  GormDB,
	}
}

func (p *PersistenceDB) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return p.Pool.BeginTx(ctx, pgx.TxOptions{})
}

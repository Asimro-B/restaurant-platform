package persistencedb

import (
	"context"
	db "restaurant-platform/database/sqlc/gen"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PersistenceDB struct {
	*db.Queries
	Pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PersistenceDB {
	return &PersistenceDB{
		Queries: db.New(pool),
		Pool:    pool,
	}
}

func (p *PersistenceDB) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return p.Pool.BeginTx(ctx, pgx.TxOptions{})
}

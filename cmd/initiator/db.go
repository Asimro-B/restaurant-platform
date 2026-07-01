package initiator

import (
	"context"
	"fmt"
	"os"
	persistencedb "restaurant-platform/database/sqlc/persistenceDB"
	"restaurant-platform/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitiatePersistenceDB() (*persistencedb.PersistenceDB, error) {
	ctx := context.Background()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	logger.Log.Info("Connected to PostgreSQL")
	return persistencedb.New(pool), nil
}

package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

// Store defines all function to execute db queries and transaction
type Store interface {
	Querier
	CreateTransferTx(ctx context.Context, arg CreateTransactionParams) (*Transfer, error)
	SetSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUserTx(ctx context.Context, arg CreateUserParams) (User, error)
	Health(ctx context.Context) StoreHealthRecord
}

// SqlStore provides all function to execute sql queries and transaction

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// ExecTX executes a function within a database transaction
func (t *SQLStore) ExecTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := t.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err : %v, rb err : %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

// Health checks the health of the database connection by pinging the database.
func (t *SQLStore) Health(ctx context.Context) StoreHealthRecord {
	var record StoreHealthRecord

	err := t.connPool.Ping(ctx)
	if err != nil {
		record.Status = "down"
		record.Error = fmt.Sprintf("db down: %v", err)
		return record
	}

	poolStat := t.connPool.Stat()

	return StoreHealthRecord{
		Status:              "up",
		Message:             "It's healthy",
		TotalConnections:    poolStat.TotalConns(),
		IdleConnections:     poolStat.IdleConns(),
		AcquiredConnections: poolStat.AcquiredConns(),
	}
}

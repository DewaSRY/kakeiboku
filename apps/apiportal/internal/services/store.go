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

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg *CreateTransferParams) (*TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg *CreateUserTxParams) (*CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg *VerifyEmailTxParams) (*VerifyEmailTxResult, error)
}

// Store provides all functions to execute db queries and transactions
type SqlStore struct {
	*Queries
	db *pgxpool.Pool
}

// NewStore creates a new Store
func NewStore(db *pgxpool.Pool) Store {
	return &SqlStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (s *SqlStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.Queries.WithTx(tx)
	err = fn(qtx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

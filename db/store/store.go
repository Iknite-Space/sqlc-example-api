package store

import (
	"context"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Begin(ctx context.Context) (repo.Querier, pgx.Tx, error)
	Do() repo.Querier
}

type Impl struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Impl {
	return &Impl{db: db}
}

func (u *Impl) Begin(ctx context.Context) (repo.Querier, pgx.Tx, error) {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	return repo.New(tx), tx, nil
}

func (u *Impl) Do() repo.Querier {
	return repo.New(u.db)
}

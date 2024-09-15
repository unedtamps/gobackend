package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Store struct {
	Querier
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		Querier: New(db),
	}
}

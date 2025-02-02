package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{db: db}
}

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"wallet-app/internal/app/app_errors"
	"wallet-app/internal/app/domain"
)

// CreateWallet создает новый кошелек с начальным балансом
func (r *WalletRepository) CreateWallet(ctx context.Context) (domain.Wallet, error) {
	// Генерируем новый UUID для кошелька
	walletID := uuid.New()

	// Вставляем новый кошелек в базу данных с начальным балансом 0
	_, err := r.db.Exec(ctx, "INSERT INTO wallets(wallet_id, balance) VALUES($1, $2)", walletID, decimal.Zero.String())
	if err != nil {
		return domain.Wallet{}, err
	}

	// Возвращаем созданный кошелек с балансом 0
	newWallet := domain.Wallet{
		ID:      walletID,
		Balance: decimal.Zero,
	}

	return newWallet, nil
}

func (r *WalletRepository) GetWallet(ctx context.Context, walletID uuid.UUID) (*decimal.Decimal, error) {
	var wallet domain.Wallet
	var balanceStr string

	err := r.db.QueryRow(ctx, "SELECT wallet_id, balance FROM wallets WHERE wallet_id=$1", walletID).
		Scan(&wallet.ID, &balanceStr)
	if err != nil {
		return nil, err
	}

	wallet.Balance, err = decimal.NewFromString(balanceStr)
	if err != nil {
		return nil, err
	}

	return &wallet.Balance, nil
}

func (r *WalletRepository) UpdateBalance(ctx context.Context, walletID uuid.UUID, amount decimal.Decimal) error {
	// Начало транзакции
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Установка уровня изоляции транзакции
	_, err = tx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	defer tx.Rollback(ctx)

	// Получение текущего баланса кошелька с блокировкой строки для обновления
	var balanceStr string
	err = tx.QueryRow(ctx, "SELECT balance FROM wallets WHERE wallet_id=$1 FOR UPDATE", walletID).Scan(&balanceStr)
	if err != nil {
		return err
	}

	balance, err := decimal.NewFromString(balanceStr)
	if err != nil {
		return err
	}

	newBalance := balance.Add(amount)
	if newBalance.LessThan(decimal.Zero) {
		return app_errors.ErrInsufficientFunds
	}

	_, err = tx.Exec(ctx, "UPDATE wallets SET balance = $1 WHERE wallet_id = $2", newBalance.String(), walletID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

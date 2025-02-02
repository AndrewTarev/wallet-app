package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"wallet-app/internal/app/domain"
	"wallet-app/internal/app/repository"
)

type WalletService struct {
	repo *repository.WalletRepository
}

func NewWalletService(repo *repository.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

// CreateWallet создает новый кошелек с нулевым балансом
func (s *WalletService) CreateWallet(ctx context.Context) (domain.Wallet, error) {
	return s.repo.CreateWallet(ctx)
}

// ProcessOperation обрабатывает операцию пополнения или снятия средств
func (s *WalletService) ProcessOperation(ctx context.Context, op domain.WalletOperation) error {
	signedAmount, err := op.GetSignedAmount()
	if err != nil {
		return fmt.Errorf("failed to get signed amount: %w", err)
	}

	// Обновляем баланс
	return s.repo.UpdateBalance(ctx, op.WalletID, signedAmount)
}

// GetBalance возвращает баланс кошелька
func (s *WalletService) GetBalance(ctx context.Context, walletID uuid.UUID) (decimal.Decimal, error) {
	// Получаем баланс кошелька через репозиторий
	balance, err := s.repo.GetWallet(ctx, walletID)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return *balance, nil
}

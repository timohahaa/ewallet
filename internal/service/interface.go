package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/timohahaa/ewallet/internal/entity"
)

type WalletService interface {
	CreateWallet(ctx context.Context) (entity.Wallet, error)
	Transfer(ctx context.Context, from, to uuid.UUID, amount float32) error
	TransactionHistory(ctx context.Context, walletId uuid.UUID) ([]entity.Transaction, error)
	WalletStatus(ctx context.Context, walletId uuid.UUID) (entity.Wallet, error)
}

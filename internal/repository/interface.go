package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/timohahaa/ewallet/internal/entity"
)

type WalletRepo interface {
	CreateWallet(ctx context.Context) (entity.Wallet, error)
	Transfer(ctx context.Context, from, to uuid.UUID, amount int64) error
	GetTransactionHistory(ctx context.Context, walletId uuid.UUID) ([]entity.Transaction, error)
	GetWalletStatus(ctx context.Context, walletId uuid.UUID) (entity.Wallet, error)
}

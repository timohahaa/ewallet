package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/timohahaa/ewallet/internal/entity"
	"github.com/timohahaa/ewallet/internal/repository"
	"github.com/timohahaa/ewallet/internal/repository/repoerrors"
)

type walletServiceImpl struct {
	walletRepo repository.WalletRepo
	log        *logrus.Logger
}

func NewWalletService(wr repository.WalletRepo, log *logrus.Logger) *walletServiceImpl {
	return &walletServiceImpl{
		walletRepo: wr,
		log:        log,
	}
}

func (ws *walletServiceImpl) CreateWallet(ctx context.Context) (entity.Wallet, error) {
	wallet, err := ws.walletRepo.CreateWallet(ctx)
	if err != nil {
		ws.log.Error("walletServiceImpl.CreateWallet - walletRepo.CreateWallet", "err", err)
		return entity.Wallet{}, err
	}
	return wallet, nil
}

func (ws *walletServiceImpl) Transfer(ctx context.Context, from, to uuid.UUID, amount float32) error {
	err := ws.walletRepo.Transfer(ctx, from, to, amount)
	if errors.Is(err, repoerrors.ErrWalletNotFound) {
		return ErrWalletNotFound
	}
	if errors.Is(err, repoerrors.ErrTargetWalletNotFound) {
		return ErrTargetWalletNotFound
	}
	if errors.Is(err, repoerrors.ErrNotEnoughBalance) {
		return ErrNotEnoughBalance
	}
	return err
}

func (ws *walletServiceImpl) TransactionHistory(ctx context.Context, walletId uuid.UUID) ([]entity.Transaction, error) {
	txs, err := ws.walletRepo.GetTransactionHistory(ctx, walletId)
	if errors.Is(err, repoerrors.ErrWalletNotFound) {
		return nil, ErrWalletNotFound
	}
	return txs, err
}

func (ws *walletServiceImpl) WalletStatus(ctx context.Context, walletId uuid.UUID) (entity.Wallet, error) {
	wallet, err := ws.walletRepo.GetWalletStatus(ctx, walletId)
	if errors.Is(err, repoerrors.ErrWalletNotFound) {
		return entity.Wallet{}, ErrWalletNotFound
	}
	return wallet, err
}

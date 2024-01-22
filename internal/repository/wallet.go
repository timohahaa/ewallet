package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/timohahaa/ewallet/internal/entity"
	"github.com/timohahaa/ewallet/internal/repository/repoerrors"
	"github.com/timohahaa/postgres"
)

const (
	InitialWalletBalance = 100.0
)

type walletRepoImpl struct {
	db *postgres.Postgres
}

func NewWalletRepo(db *postgres.Postgres) *walletRepoImpl {
	return &walletRepoImpl{
		db: db,
	}
}

// создание нового кошелька
func (wr *walletRepoImpl) CreateWallet(ctx context.Context) (entity.Wallet, error) {
	newWalletID, err := uuid.NewRandom()
	if err != nil {
		slog.Error("walletRepoImpl.CreateWallet - uuid.NewRandom", "err", err)
		return entity.Wallet{}, err
	}

	sql, args, err := wr.db.Builder.
		Insert("wallets").
		Columns("id", "balance").
		Values(newWalletID, InitialWalletBalance).
		ToSql()
	if err != nil {
		slog.Error("walletRepoImpl.CreateWallet - db.Builder", "err", err)
		return entity.Wallet{}, err
	}

	_, err = wr.db.ConnPool.Exec(ctx, sql, args)

	if err != nil {
		slog.Error("walletRepoImpl.CreateWallet - db.ConnPool.Exec", "err", err)
		return entity.Wallet{}, err
	}
	return entity.Wallet{Id: newWalletID, Balance: InitialWalletBalance}, nil
}

// вспомогательные функции для совершения транзакции - Dont Repeat Youtself ;)
func (wr *walletRepoImpl) getWallet(ctx context.Context, walletId uuid.UUID) (entity.Wallet, error) {
	sql, args, err := wr.db.Builder.
		Select("wallets").
		Columns("id", "balance").
		Where("id = ?", walletId).
		ToSql()
	if err != nil {
		slog.Error("walletRepoImpl.getWallet - db.Builder", "err", err)
		return entity.Wallet{}, err
	}

	var wallet entity.Wallet
	err = wr.db.ConnPool.QueryRow(ctx, sql, args).Scan(&wallet)
	// кошелек не найден
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Wallet{}, repoerrors.ErrWalletNotFound
	}
	if err != nil {
		slog.Error("walletRepoImpl.getWallet - db.ConnPool.QueryRow", "err", err)
		return entity.Wallet{}, err
	}

	return wallet, nil
}

func (wr *walletRepoImpl) updateWallet(ctx context.Context, walletId uuid.UUID, newBalance float32) error {
	sql, args, err := wr.db.Builder.
		Update("wallets").
		Set("balance", newBalance).
		Where("id = ?", walletId).
		ToSql()

	_, err = wr.db.ConnPool.Exec(ctx, sql, args)
	if err != nil {
		slog.Error("walletRepoImpl.updateWallet - db.ConnPool.Exec", "err", err)
		return err
	}
	return nil
}

// совершение транзакции
func (wr *walletRepoImpl) Transfer(ctx context.Context, from, to uuid.UUID, amount float32) error {
	// проверяем исходящий кошелек
	fromWallet, err := wr.getWallet(ctx, from)

	// исходящий кошелек не найден
	if errors.Is(err, pgx.ErrNoRows) {
		return repoerrors.ErrWalletNotFound
	}
	if err != nil {
		slog.Error("walletRepoImpl.Transfer - getWallet", "err", err)
		return err
	}
	// баланса не достаточно для перевода
	if fromWallet.Balance-amount < 0 {
		return repoerrors.ErrNotEnoughBalance
	}

	// проверяем целевой кошелек
	toWallet, err := wr.getWallet(ctx, to)
	// целевой кошелек не найден
	if errors.Is(err, pgx.ErrNoRows) {
		return repoerrors.ErrTargetWalletNotFound
	}
	if err != nil {
		slog.Error("walletRepoImpl.Transfer - getWallet", "err", err)
		return err
	}

	//обновляем балансы и сохраняем транзакцию
	txTime := time.Now().UTC()
	fromWallet.Balance -= amount
	toWallet.Balance += amount
	err = wr.updateWallet(ctx, fromWallet.Id, fromWallet.Balance)
	if err != nil {
		slog.Error("walletRepoImpl.Transfer - updateWallet", "err", err)
		return err
	}
	err = wr.updateWallet(ctx, toWallet.Id, toWallet.Balance)
	if err != nil {
		slog.Error("walletRepoImpl.Transfer - updateWallet", "err", err)
		return err
	}

	sql, args, err := wr.db.Builder.
		Insert("transactions").
		Columns("made_at", "transfered_from", "transfered_to", "amount").
		Values(txTime, fromWallet.Id, toWallet.Id, amount).
		ToSql()
	if err != nil {
		slog.Error("walletRepoImpl.Transfer - db.Builder", "err", err)
		return err
	}

	_, err = wr.db.ConnPool.Exec(ctx, sql, args)
	if err != nil {
		slog.Error("walletRepoImpl.Transfer - db.ConnPool.Exec", "err", err)
		return err
	}

	return nil
}

func (wr *walletRepoImpl) GetTransactionHistory(ctx context.Context, walletId uuid.UUID) ([]entity.Transaction, error) {
	return nil, nil
}

func (wr *walletRepoImpl) GetWalletStatus(ctx context.Context, walletId uuid.UUID) (entity.Wallet, error) {
	wallet, err := wr.getWallet(ctx, walletId)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Wallet{}, repoerrors.ErrWalletNotFound
	}
	if err != nil {
		slog.Error("walletRepoImpl.GetWalletStatus - getWallet", "err", err)
		return entity.Wallet{}, err
	}

	return wallet, nil
}

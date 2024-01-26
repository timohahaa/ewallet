package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"github.com/timohahaa/ewallet/internal/entity"
	"github.com/timohahaa/ewallet/internal/repository/repoerrors"
	"github.com/timohahaa/postgres"
)

const (
	InitialWalletBalance = 100.0
)

type walletRepoImpl struct {
	db  *postgres.Postgres
	log *logrus.Logger
}

func NewWalletRepo(db *postgres.Postgres, log *logrus.Logger) *walletRepoImpl {
	return &walletRepoImpl{
		db:  db,
		log: log,
	}
}

// создание нового кошелька
func (wr *walletRepoImpl) CreateWallet(ctx context.Context) (entity.Wallet, error) {
	newWalletID, err := uuid.NewRandom()
	if err != nil {
		wr.log.Error("walletRepoImpl.CreateWallet - uuid.NewRandom", "err", err)
		return entity.Wallet{}, err
	}

	sql, args, err := wr.db.Builder.
		Insert("wallets").
		Columns("id", "balance").
		Values(newWalletID, InitialWalletBalance).
		ToSql()

	if err != nil {
		wr.log.Error("walletRepoImpl.CreateWallet - db.Builder", "err", err)
		return entity.Wallet{}, err
	}

	_, err = wr.db.ConnPool.Exec(ctx, sql, args...)

	if err != nil {
		wr.log.Error("walletRepoImpl.CreateWallet - db.ConnPool.Exec", "err", err)
		return entity.Wallet{}, err
	}
	return entity.Wallet{Id: newWalletID, Balance: InitialWalletBalance}, nil
}

// вспомогательные функции для совершения транзакции - Dont Repeat Youtself ;)
func (wr *walletRepoImpl) getWallet(ctx context.Context, walletId uuid.UUID) (entity.Wallet, error) {
	sql, args, err := wr.db.Builder.
		Select("id", "balance").
		From("wallets").
		Where("id = ?", walletId).
		ToSql()
	if err != nil {
		wr.log.Error("walletRepoImpl.getWallet - db.Builder", "err", err)
		return entity.Wallet{}, err
	}

	var wallet entity.Wallet
	err = wr.db.ConnPool.QueryRow(ctx, sql, args...).Scan(&wallet.Id, &wallet.Balance)
	// кошелек не найден
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Wallet{}, pgx.ErrNoRows
	}
	if err != nil {
		wr.log.Error("walletRepoImpl.getWallet - db.ConnPool.QueryRow", "err", err)
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

	_, err = wr.db.ConnPool.Exec(ctx, sql, args...)
	if err != nil {
		wr.log.Error("walletRepoImpl.updateWallet - db.ConnPool.Exec", "err", err)
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
		wr.log.Error("walletRepoImpl.Transfer - getWallet", "err", err)
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
		wr.log.Error("walletRepoImpl.Transfer - getWallet", "err", err)
		return err
	}

	//обновляем балансы и сохраняем транзакцию
	txTime := time.Now().UTC()
	fromWallet.Balance -= amount
	toWallet.Balance += amount
	err = wr.updateWallet(ctx, fromWallet.Id, fromWallet.Balance)
	if err != nil {
		wr.log.Error("walletRepoImpl.Transfer - updateWallet", "err", err)
		return err
	}
	err = wr.updateWallet(ctx, toWallet.Id, toWallet.Balance)
	if err != nil {
		wr.log.Error("walletRepoImpl.Transfer - updateWallet", "err", err)
		return err
	}

	sql, args, err := wr.db.Builder.
		Insert("transactions").
		Columns("made_at", "transfered_from", "transfered_to", "amount").
		Values(txTime, fromWallet.Id, toWallet.Id, amount).
		ToSql()
	if err != nil {
		wr.log.Error("walletRepoImpl.Transfer - db.Builder", "err", err)
		return err
	}

	_, err = wr.db.ConnPool.Exec(ctx, sql, args...)
	if err != nil {
		wr.log.Error("walletRepoImpl.Transfer - db.ConnPool.Exec", "err", err)
		return err
	}

	return nil
}

func (wr *walletRepoImpl) GetTransactionHistory(ctx context.Context, walletId uuid.UUID) ([]entity.Transaction, error) {
	// проверка на существование кошелька
	wallet, err := wr.getWallet(ctx, walletId)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repoerrors.ErrWalletNotFound
	}
	if err != nil {
		wr.log.Error("walletRepoImpl.Transfer - getWallet", "err", err)
		return nil, err
	}

	sql, args, err := wr.db.Builder.
		Select("made_at", "transfered_from", "transfered_to", "amount").
		From("transactions").
		Where("transfered_from = ? OR transfered_to = ?", wallet.Id, wallet.Id).
		ToSql()
	if err != nil {
		wr.log.Error("walletRepoImpl.GetTransactionHistory - db.Builder", "err", err)
		return nil, err
	}

	rows, err := wr.db.ConnPool.Query(ctx, sql, args...)
	if err != nil {
		wr.log.Error("walletRepoImpl.GetTransactionHistory - db.ConnPool.Query", "err", err)
		return nil, err
	}

	var transactions []entity.Transaction
	for rows.Next() {
		var tx entity.Transaction
		// игнорируем ошибку, но:
		// можно бы было сделать ошибку ErrScan или типа того, и записывать ее в переменную
		// в скоупе вне цикла, а затем возвращать неполный список транзакций и ошибку
		_ = rows.Scan(&tx.Time, &tx.From, &tx.To, &tx.Amount)
		transactions = append(transactions, tx)
	}

	// если не нашли транзакции - вернем nil-слайс, он в json пойдет как null
	if len(transactions) == 0 {
		return []entity.Transaction(nil), nil
	}
	return transactions, nil
}

func (wr *walletRepoImpl) GetWalletStatus(ctx context.Context, walletId uuid.UUID) (entity.Wallet, error) {
	wallet, err := wr.getWallet(ctx, walletId)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Wallet{}, repoerrors.ErrWalletNotFound
	}
	if err != nil {
		wr.log.Error("walletRepoImpl.GetWalletStatus - getWallet", "err", err)
		return entity.Wallet{}, err
	}

	return wallet, nil
}

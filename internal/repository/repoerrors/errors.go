package repoerrors

import "errors"

var (
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrTargetWalletNotFound = errors.New("target wallet not found")
	ErrNotEnoughBalance     = errors.New("not enough balance")
)

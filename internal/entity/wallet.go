package entity

import "github.com/google/uuid"

type Wallet struct {
	Id      uuid.UUID `json:"id"`
	Balance float32   `json:"balance"`
}

func NewWallet(id uuid.UUID, balance float32) *Wallet {
	return &Wallet{
		Id:      id,
		Balance: balance,
	}
}

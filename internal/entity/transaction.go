package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Time   time.Time `json:"time"`
	From   uuid.UUID `json:"from"`
	To     uuid.UUID `json:"to"`
	Amount float32   `json:"amount"`
}

func NewTransaction(time time.Time, from, to uuid.UUID, amount float32) *Transaction {
	return &Transaction{
		Time:   time,
		From:   from,
		To:     to,
		Amount: amount,
	}
}

package entity

import (
	"sync"
)

type Amount int64

type balance struct {
	amount Amount

	mu sync.Mutex
}

func (b *balance) Amount() Amount {
	return b.amount
}

func (b *balance) withdraw(amount Amount) error {
	if b.amount-amount < 0 {
		return ErrInsufficientBalance
	}
	b.amount -= amount
	return nil
}

func (b *balance) deposit(amount Amount) error {
	b.amount += amount
	return nil
}

package entity

import (
	"errors"
	"sync"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrZeroAmount          = errors.New("zero amount isn't permit")
)

type Amount int64

type Balance struct {
	Amount Amount

	mu sync.Mutex
}

func (b *Balance) withdraw(amount Amount) error {
	if amount == 0 {
		return ErrZeroAmount
	}
	if b.Amount-amount < 0 {
		return ErrInsufficientBalance
	}
	b.Amount -= amount
	return nil
}

func (b *Balance) deposit(amount Amount) {
	b.Amount += amount
}

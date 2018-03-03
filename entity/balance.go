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

type balance struct {
	amount Amount

	mu sync.Mutex
}

func (b *balance) Amount() Amount {
	return b.amount
}

func (b *balance) withdraw(amount Amount) error {
	if amount == 0 {
		return ErrZeroAmount
	}
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

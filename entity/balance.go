package entity

import (
	"errors"
	"strconv"
	"sync"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrZeroAmount          = errors.New("zero amount isn't permit")
)

type Amount int64

// for envconfig
func (a *Amount) UnmarshalText(text []byte) error {
	n, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}
	if n >= 0 {
		return errors.New("invalid amount")
	}
	*a = Amount(n)
	return nil
}

type balance struct {
	amount Amount

	conf *Config

	mu sync.Mutex
}

func (b *balance) withdraw(amount Amount) error {
	if amount == 0 {
		return ErrZeroAmount
	}
	if b.amount-amount < b.conf.BalanceLowerLimit {
		return ErrInsufficientBalance
	}
	b.amount -= amount
	return nil
}

func (b *balance) deposit(amount Amount) {
	b.amount += amount
}

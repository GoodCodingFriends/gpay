package entity

import (
	"errors"
	"strconv"
	"sync"

	"github.com/GoodCodingFriends/gpay/config"
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

	conf *config.Config

	mu sync.Mutex
}

func (b *balance) checkAmount(amount Amount) error {
	if amount == 0 {
		return ErrZeroAmount
	}
	if b.amount-amount < Amount(b.conf.Entity.BalanceLowerLimit) {
		return ErrInsufficientBalance
	}
	return nil
}

func (b *balance) withdraw(amount Amount) error {
	err := b.checkAmount(amount)
	if err != nil {
		return err
	}
	b.amount -= amount
	return nil
}

func (b *balance) deposit(amount Amount) {
	b.amount += amount
}

func GetBalanceAmount(user *User) Amount {
	return user.balance.amount
}

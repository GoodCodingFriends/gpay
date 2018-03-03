package entity

import (
	"errors"
)

var (
	ErrSameUser            = errors.New("cannot pay/claim between same users")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type UserID string

type User struct {
	ID          UserID
	FirstName   string
	LastName    string
	DisplayName string
	balance     *balance
}

func (u *User) Pay(to *User, amount Amount, message string) (*Transaction, error) {
	if u.ID == to.ID {
		return nil, ErrSameUser
	}

	u.balance.mu.Lock()
	defer u.balance.mu.Unlock()

	var err error
	rollback := func() func() {
		b1, b2 := *u.balance, *to.balance
		return func() {
			if err != nil {
				u.balance = &b1
				to.balance = &b2
			}
		}
	}()
	defer rollback()

	err = u.balance.withdraw(amount)
	if err != nil {
		return nil, err
	}
	err = to.balance.deposit(amount)
	if err != nil {
		return nil, err
	}

	return newTransaction(TxTypePay, u.ID, to.ID, amount, message), nil
}

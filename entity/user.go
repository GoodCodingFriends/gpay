package entity

import (
	"errors"

	"github.com/GoodCodingFriends/gpay/entity/internal/id"
)

var (
	ErrSameUser         = errors.New("cannot pay/claim between same users")
	ErrWrongDestination = errors.New("destination user is wrong")
	ErrCompletedInvoice = errors.New("the invoice is already completed")
)

type UserID string

type User struct {
	ID          UserID
	FirstName   string
	LastName    string
	DisplayName string
	balance     balance
}

func NewUser(conf *Config, firstName, lastName, displayName string) *User {
	return &User{
		ID:          UserID(id.New()),
		FirstName:   firstName,
		LastName:    lastName,
		DisplayName: displayName,
		balance: balance{
			conf: conf,
		},
	}
}

func (u *User) Pay(to *User, amount Amount, message string) (*Transaction, error) {
	if u.ID == to.ID {
		return nil, ErrSameUser
	}

	if err := send(u, to, amount); err != nil {
		return nil, err
	}

	return newTransaction(TxTypePay, u.ID, to.ID, amount, message), nil
}

func (u *User) Claim(to *User, amount Amount, message string) (*Invoice, error) {
	if u.ID == to.ID {
		return nil, ErrSameUser
	}
	if err := to.balance.checkAmount(amount); err != nil {
		return nil, err
	}

	return newInvoice(u.ID, to.ID, amount, message), nil
}

func (u *User) AcceptInvoice(invoice *Invoice, from *User) (*Transaction, error) {
	if invoice.to != u.ID {
		return nil, ErrWrongDestination
	}

	if invoice.IsCompleted {
		return nil, ErrCompletedInvoice
	}

	if err := send(u, from, invoice.amount); err != nil {
		return nil, err
	}

	return newTransaction(TxTypeClaim, from.ID, u.ID, invoice.amount, invoice.message), nil
}

package entity

import (
	"errors"
)

var (
	ErrSameUser         = errors.New("cannot pay/claim between same users")
	ErrWrongDestination = errors.New("destination user is wrong")
)

type UserID string

type User struct {
	ID          UserID
	FirstName   string
	LastName    string
	DisplayName string
	balance     *Balance
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
	if amount == Amount(0) {
		return nil, ErrZeroAmount
	}

	return newInvoice(u.ID, to.ID, amount, message), nil
}

func (u *User) AcceptInvoice(invoice *Invoice, from *User) (*Transaction, error) {
	if invoice.to != u.ID {
		return nil, ErrWrongDestination
	}

	if err := send(u, from, invoice.amount); err != nil {
		return nil, err
	}

	return newTransaction(TxTypeClaim, from.ID, u.ID, invoice.amount, invoice.message), nil
}

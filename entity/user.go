package entity

import (
	"errors"

	"github.com/GoodCodingFriends/gpay/config"
)

var (
	ErrSameUser         = errors.New("cannot pay/claim between same users")
	ErrWrongDestination = errors.New("destination user is wrong")
	ErrCompletedInvoice = errors.New("the invoice is already completed")
)

type UserIDGenerator interface {
	Generate() UserID
}

type UserID string

type User struct {
	ID          UserID
	FirstName   string
	LastName    string
	DisplayName string
	balance     balance
}

// TODO: remove the dependency related to id generation

func NewUser(conf *config.Config, id UserID, firstName, lastName, displayName string) *User {
	return &User{
		ID:          id,
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
	if invoice.ToID != u.ID {
		return nil, ErrWrongDestination
	}

	if invoice.Status != StatusPending {
		return nil, ErrCompletedInvoice
	}

	invoice.Status = StatusAccepted

	if err := send(u, from, invoice.amount); err != nil {
		return nil, err
	}

	return newTransaction(TxTypeClaim, from.ID, u.ID, invoice.amount, invoice.message), nil
}

func (u *User) RejectInvoice(invoice *Invoice, from *User) error {
	if invoice.ToID != u.ID {
		return ErrWrongDestination
	}

	if invoice.Status != StatusPending {
		return ErrCompletedInvoice
	}

	invoice.Status = StatusRejected

	return nil
}

func (u *User) BalanceAmount() Amount {
	return u.balance.amount
}

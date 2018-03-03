package entity

import "errors"

var (
	ErrWrongDestination = errors.New("destination user is wrong")
)

type InvoiceID string

type Invoice struct {
	ID       InvoiceID
	from, to UserID
	amount   Amount
	message  string
}

func newInvoice(from, to UserID, amount Amount, message string) *Invoice {
	return &Invoice{
		ID:      InvoiceID(newID()),
		from:    from,
		to:      to,
		amount:  amount,
		message: message,
	}
}

func (i *Invoice) accept(user *User) error {
	if user.ID != i.to {
		return ErrWrongDestination
	}
	return nil
}

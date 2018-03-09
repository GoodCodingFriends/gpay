package entity

import "github.com/GoodCodingFriends/gpay/entity/internal/id"

type InvoiceID string

type Status int

const (
	StatusPending  = 0
	StatusAccepted = 1
	StatusRejected = 2
)

type Invoice struct {
	ID           InvoiceID
	Status       Status
	FromID, ToID UserID
	Amount       Amount
	Message      string
}

func newInvoice(from, to UserID, amount Amount, message string) *Invoice {
	return &Invoice{
		ID:      InvoiceID(id.New()),
		FromID:  from,
		ToID:    to,
		Amount:  amount,
		Message: message,
	}
}

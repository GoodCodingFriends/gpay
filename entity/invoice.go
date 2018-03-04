package entity

import "github.com/GoodCodingFriends/gpay/entity/internal/id"

type InvoiceID string

type Invoice struct {
	ID       InvoiceID
	from, to UserID
	amount   Amount
	message  string
}

func newInvoice(from, to UserID, amount Amount, message string) *Invoice {
	return &Invoice{
		ID:      InvoiceID(id.New()),
		from:    from,
		to:      to,
		amount:  amount,
		message: message,
	}
}

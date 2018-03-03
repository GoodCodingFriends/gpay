package entity

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

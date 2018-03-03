package entity

type TxID string

type TxType int

const (
	TxTypePay TxType = iota
	TxTypeClaim
)

type Transaction struct {
	ID       TxID
	Type     TxType
	From, To UserID
	Amount   Amount
	Message  string
}

func newTransaction(typ TxType, from, to UserID, amount Amount, message string) *Transaction {
	return &Transaction{
		ID:      TxID(newID()),
		Type:    typ,
		From:    from,
		To:      to,
		Amount:  amount,
		Message: message,
	}
}

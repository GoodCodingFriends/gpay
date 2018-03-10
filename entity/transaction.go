package entity

import "github.com/GoodCodingFriends/gpay/entity/internal/id"

type TxID string

type TxType int

var txTypeToString = map[TxType]string{
	TxTypePay:   "pay",
	TxTypeClaim: "claim",
}

func (t TxType) String() string {
	return txTypeToString[t]
}

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
		ID:      TxID(id.New()),
		Type:    typ,
		From:    from,
		To:      to,
		Amount:  amount,
		Message: message,
	}
}

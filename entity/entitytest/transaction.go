package entitytest

import (
	"testing"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/entity/internal/id"
)

func NewTransaction(t *testing.T) *entity.Transaction {
	return &entity.Transaction{
		ID:      entity.TxID(id.New()),
		Type:    entity.TxTypePay,
		From:    entity.UserID("asuka"),
		To:      entity.UserID("kaori"),
		Amount:  entity.Amount(100),
		Message: "",
	}
}

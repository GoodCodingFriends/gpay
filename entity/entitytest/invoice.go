package entitytest

import (
	"testing"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/entity/internal/id"
)

func NewInvoice(t *testing.T) *entity.Invoice {
	return &entity.Invoice{
		ID:      entity.InvoiceID(id.New()),
		Status:  entity.StatusPending,
		FromID:  entity.UserID("nakagawa"),
		ToID:    entity.UserID("yoshikawa"),
		Amount:  entity.Amount(100),
		Message: "",
	}
}

package repository

import (
	"context"

	"github.com/GoodCodingFriends/gpay/entity"
)

type Repository struct {
	User        userRepository
	Invoice     invoiceRepository
	Transaction txRepository

	transactor
}

type userRepository interface {
	FindByID(context.Context, entity.UserID) (*entity.User, error)
	FindAll(context.Context) ([]*entity.User, error)
	Store(context.Context, *entity.User) error
	StoreAll(context.Context, []*entity.User) error
}

type invoiceRepository interface {
	FindByID(context.Context, entity.InvoiceID) (*entity.Invoice, error)
	FindAll(context.Context) ([]*entity.Invoice, error)
	Store(context.Context, *entity.Invoice) error
	StoreAll(context.Context, []*entity.Invoice) error
}

type txRepository interface {
	FindByID(context.Context, entity.TxID) (*entity.Transaction, error)
	FindAll(context.Context) ([]*entity.Transaction, error)
	Store(context.Context, *entity.Transaction) error
}

package repository

import (
	"context"

	"github.com/GoodCodingFriends/gpay/entity"
)

type UserRepository interface {
	FindByID(context.Context, entity.UserID) (*entity.User, error)
	FindAll(context.Context) ([]*entity.User, error)
	Store(context.Context, *entity.User) error
	StoreAll(context.Context, []*entity.User) error
}

type InvoiceRepository interface {
	FindByID(context.Context, entity.InvoiceID) (*entity.Invoice, error)
	FindAll(context.Context) ([]*entity.Invoice, error)
	Store(context.Context, *entity.Invoice) error
	StoreAll(context.Context, []*entity.Invoice) error
}

type TxRepository interface {
	FindByID(context.Context, entity.TxID) (*entity.Transaction, error)
	FindAll(context.Context) ([]*entity.Transaction, error)
	FindAllByUserID(context.Context, entity.UserID) ([]*entity.Transaction, error)
	Store(context.Context, *entity.Transaction) error
}

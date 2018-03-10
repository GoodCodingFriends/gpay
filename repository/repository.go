package repository

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("specified entity not found")
)

type Repository struct {
	beginner TxBeginner

	User        UserRepository
	Invoice     InvoiceRepository
	Transaction TxRepository
}

func (r *Repository) BeginTx(ctx context.Context) (*Tx, context.Context, error) {
	committer, wrappedCtx, err := r.beginner.BeginTx(ctx)
	if err != nil {
		return nil, nil, err
	}
	return &Tx{
		committer:   committer,
		User:        r.User,
		Invoice:     r.Invoice,
		Transaction: r.Transaction,
	}, wrappedCtx, nil
}

func (r *Repository) Close() error {
	return r.beginner.Close()
}

func New(txBeginner TxBeginner, userRepo UserRepository, invoiceRepo InvoiceRepository, txRepo TxRepository) *Repository {
	return &Repository{
		beginner:    txBeginner,
		User:        userRepo,
		Invoice:     invoiceRepo,
		Transaction: txRepo,
	}
}

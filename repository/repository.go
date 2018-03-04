package repository

import (
	"context"
)

type Repository struct {
	beginner TxBeginner

	User        UserRepository
	Invoice     InvoiceRepository
	Transaction TxRepository
}

func (r *Repository) BeginTx(ctx context.Context) (*Tx, error) {
	committer, err := r.beginner.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	return &Tx{
		committer:   committer,
		User:        r.User,
		Invoice:     r.Invoice,
		Transaction: r.Transaction,
	}, nil
}

func New(txBeginner TxBeginner, userRepo UserRepository, invoiceRepo InvoiceRepository, txRepo TxRepository) *Repository {
	return &Repository{
		beginner:    txBeginner,
		User:        userRepo,
		Invoice:     invoiceRepo,
		Transaction: txRepo,
	}
}

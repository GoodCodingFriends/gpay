package repository

import "context"

type transactor interface {
	BeginTx(context.Context) *Tx
}

type Tx struct {
	User        userRepository
	Invoice     invoiceRepository
	Transaction txRepository
}

func (t *Tx) Commit() error {
	panic("not implemented yet")
	return nil
}

func (t *Tx) Rollback() error {
	panic("not implemented yet")
	return nil
}

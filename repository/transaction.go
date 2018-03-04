package repository

import "context"

type TxBeginner interface {
	BeginTx(context.Context) (TxCommitter, error)
}

type TxCommitter interface {
	Commit() error
	Rollback() error
}

type Tx struct {
	committer TxCommitter

	User        UserRepository
	Invoice     InvoiceRepository
	Transaction TxRepository
}

func (t *Tx) Commit() error {
	panic("not implemented yet")
	return nil
}

func (t *Tx) Rollback() error {
	panic("not implemented yet")
	return nil
}

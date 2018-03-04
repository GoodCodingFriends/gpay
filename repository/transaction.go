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
	return t.committer.Commit()
}

func (t *Tx) Rollback() error {
	return t.committer.Rollback()
}

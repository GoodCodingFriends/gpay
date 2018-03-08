package usecase

import (
	"context"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/repository"
)

type ClaimParam struct {
	From, To *entity.User
	Amount   entity.Amount
	Message  string
}

func Claim(repo *repository.Repository, p *ClaimParam) (*entity.Invoice, error) {
	from, to := p.From, p.To
	invoice, err := from.Claim(to, p.Amount, p.Message)
	if err != nil {
		return nil, err
	}

	err = repo.Invoice.Store(context.TODO(), invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

type AcceptInvoiceParam struct {
	InvoiceID entity.InvoiceID
}

func AcceptInvoice(repo *repository.Repository, p *AcceptInvoiceParam) (*entity.Transaction, error) {
	invoice, err := repo.Invoice.FindByID(context.Background(), p.InvoiceID)
	if err != nil {
		return nil, err
	}

	from, to, err := FindBothUsers(repo, invoice.FromID, invoice.ToID)
	if err != nil {
		return nil, err
	}

	tx, err := to.AcceptInvoice(invoice, from)
	if err != nil {
		return nil, err
	}

	dbtx, ctx, err := repo.BeginTx(context.Background())
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := recover(); err != nil {
			dbtx.Rollback()
			panic(err)
		} else if err != nil {
			dbtx.Rollback()
		}
	}()

	err = dbtx.User.StoreAll(ctx, []*entity.User{from, to})
	if err != nil {
		return nil, err
	}

	err = dbtx.Invoice.Store(ctx, invoice)
	if err != nil {
		return nil, err
	}

	err = dbtx.Transaction.Store(ctx, tx)
	if err != nil {
		return nil, err
	}

	return tx, dbtx.Commit()
}

type RejectInvoiceParam struct {
	InvoiceID entity.InvoiceID
}

func RejectInvoice(repo *repository.Repository, p *RejectInvoiceParam) error {
	invoice, err := repo.Invoice.FindByID(context.Background(), p.InvoiceID)
	if err != nil {
		return err
	}

	from, to, err := FindBothUsers(repo, invoice.FromID, invoice.ToID)
	if err != nil {
		return err
	}

	err = to.RejectInvoice(invoice, from)
	if err != nil {
		return err
	}

	dbtx, ctx, err := repo.BeginTx(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			dbtx.Rollback()
			panic(err)
		} else if err != nil {
			dbtx.Rollback()
		}
	}()

	err = dbtx.User.StoreAll(ctx, []*entity.User{from, to})
	if err != nil {
		return err
	}

	err = dbtx.Invoice.Store(ctx, invoice)
	if err != nil {
		return err
	}

	return dbtx.Commit()
}

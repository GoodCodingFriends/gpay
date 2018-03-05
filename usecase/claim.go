package usecase

import (
	"context"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/repository"
)

type ClaimParam struct {
	FromID, ToID entity.UserID
	Amount       entity.Amount
	Message      string
}

func Claim(repo *repository.Repository, p *ClaimParam) (*entity.Invoice, error) {
	from, to, err := findBothUsers(repo, p.FromID, p.ToID)
	if err != nil {
		return nil, err
	}

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
	InvoiceID    entity.InvoiceID
	FromID, ToID entity.UserID
}

func AcceptInvoice(repo *repository.Repository, p *AcceptInvoiceParam) (*entity.Transaction, error) {
	invoice, err := repo.Invoice.FindByID(context.Background(), p.InvoiceID)
	if err != nil {
		return nil, err
	}

	from, to, err := findBothUsers(repo, p.FromID, p.ToID)
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

	err = dbtx.Transaction.Store(ctx, tx)
	if err != nil {
		return nil, err
	}

	return tx, err
}

func findBothUsers(repo *repository.Repository, fromID, toID entity.UserID) (*entity.User, *entity.User, error) {
	from, err := repo.User.FindByID(context.Background(), fromID)
	if err != nil {
		return nil, nil, err
	}
	to, err := repo.User.FindByID(context.Background(), toID)
	if err != nil {
		return nil, nil, err
	}
	return from, to, nil
}
